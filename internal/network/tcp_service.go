package network

import (
	"context"
	"errors"
	"github.com/alkosuv/in-mem-kv/internal/tools"
	"github.com/alkosuv/in-mem-kv/internal/uuid"
	"io"
	"log/slog"
	"net"
	"time"
)

// DefaultAddress - дефолтный адрес базы данных
const defaultAddress string = "127.0.0.1:13666"

// defaultMessageSize - размер буфера для сообщения (запроса)
const defaultMessageSize int = 1024

// defaultMaxConnIdleTime - максимальное время бездействия соединения. Если соединение не используется в течение этого времени, оно будет закрыто:
const defaultMaxConnIdleTime time.Duration = time.Duration(time.Second * 1800)

// defaultMaxConn - максимальное количество соединений
const defaultMaxConn int = 100

// defaultMaxConnPoolRetryAttempts - максимальное количество попыток занять место в пуле соединений
const defaultMaxConnPoolRetryAttempts int = 5

// defaultConnPoolRetryTimeout - временя ожидания между попытками занять место в пуле соединений
const defaultConnPoolRetryTimeout time.Duration = time.Second

var (
	ErrConnectionLimitReached = errors.New("error connection limit reached")
)

var defaultTCPService = &TCPServer{
	address:               defaultAddress,
	idleTimeout:           defaultMaxConnIdleTime,
	messageSize:           defaultMessageSize,
	maxConn:               defaultMaxConn,
	connPoolRetryAttempts: defaultMaxConnPoolRetryAttempts,
	connPoolRetryTimeout:  defaultConnPoolRetryTimeout,
}

type TCPHandler func(context.Context, []byte) ([]byte, error)

type TCPServerConfig struct {
	IdleMax               int    `yaml:"idle-max"`
	MessageSize           int    `yaml:"message-size"`
	MaxConn               int    `yaml:"max-conn"`
	ConnPoolRetryAttempts int    `yaml:"conn-pool-retry-attempts"`
	ConnPoolRetryTimeout  int    `yaml:"conn-pool-retry-timeout"`
	Address               string `yaml:"address"`
}

type TCPServer struct {
	messageSize           int
	maxConn               int
	connPoolRetryAttempts int
	idleTimeout           time.Duration
	connPoolRetryTimeout  time.Duration
	address               string
}

func NewTCPServer(cfg TCPServerConfig) *TCPServer {
	tcp := defaultTCPService

	if cfg.Address != "" {
		tcp.address = cfg.Address
	}

	if cfg.MaxConn > 0 {
		tcp.maxConn = cfg.MaxConn
	}

	if cfg.MessageSize > 0 {
		tcp.messageSize = cfg.MessageSize
	}

	if cfg.IdleMax > 0 {
		tcp.idleTimeout = time.Duration(cfg.IdleMax) * time.Second
	}

	if cfg.ConnPoolRetryAttempts > 0 {
		tcp.connPoolRetryAttempts = cfg.ConnPoolRetryAttempts
	}

	if cfg.ConnPoolRetryTimeout > 0 {
		tcp.connPoolRetryTimeout = time.Duration(cfg.ConnPoolRetryTimeout) * time.Second
	}

	return tcp
}

func (tcp *TCPServer) Listen(ctx context.Context, handler TCPHandler) error {
	semaphore := tools.NewSemaphore(tcp.maxConn)

	listener, err := net.Listen("tcp", tcp.address)
	if err != nil {
		slog.Error("invalid init listener", "error", err)
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return err
			}

			slog.Error("invalid accept", "error", err)
			continue
		}

		go func() {
			defer conn.Close()

			isAcquire := false
			for range tcp.connPoolRetryAttempts {
				if isAcquire = semaphore.TryAcquire(); isAcquire {
					break
				}
				time.Sleep(tcp.connPoolRetryTimeout)
			}
			if !isAcquire {
				conn.Write([]byte(ErrConnectionLimitReached.Error()))
				return
			}
			defer semaphore.Release()

			if _, err := conn.Write([]byte("Ok")); err != nil {
				slog.Warn("failed to write", "error", err)
				return
			}

			tcp.handlerConn(context.Background(), conn, handler)
		}()
	}
}

func (tcp *TCPServer) handlerConn(ctx context.Context, conn net.Conn, handler TCPHandler) {
	req := make([]byte, tcp.messageSize)

	for {
		operationID, err := uuid.Generate()
		if err != nil {
			slog.Error("generate uuid", "error", err)
			break
		}
		ctx := context.WithValue(ctx, "operation_id", operationID)

		if err := conn.SetDeadline(time.Now().Add(tcp.idleTimeout)); err != nil {
			slog.Warn("filed to set deadline", "operation_id", operationID, "error", err)
			break
		}

		count, err := conn.Read(req)
		if err != nil {
			if errors.Is(err, io.EOF) {
				slog.Warn("failed to read", "operation_id", operationID, "error", err)
			}

			break
		}

		resp, err := handler(ctx, req[0:count])
		if err != nil {
			slog.Warn("failed to handle request", "operation_id", operationID, "error", err)
			resp = []byte(err.Error())
		}
		if _, err := conn.Write(resp); err != nil {
			slog.Warn("failed to write", "operation_id", operationID, "error", err)
			break
		}
	}
}
