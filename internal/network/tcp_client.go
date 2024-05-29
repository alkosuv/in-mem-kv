package network

import (
	"errors"
	"log/slog"
	"net"
	"os"
	"syscall"
	"time"
)

type TCPClient struct {
	address     string
	idleTimeout time.Duration
	messageSize int
	conn        net.Conn
}

var defaultTCPClient = &TCPClient{
	address:     defaultAddress,
	idleTimeout: defaultMaxConnIdleTime,
	messageSize: defaultMessageSize,
}

func NewTCPClient(address string, messageSize int) *TCPClient {
	client := defaultTCPClient

	if address != "" {
		client.address = address
	}

	if messageSize > 0 {
		client.messageSize = messageSize
	}

	return client
}

func (tcp *TCPClient) Connect() error {
	conn, err := net.Dial("tcp", tcp.address)
	if err != nil {
		slog.Error("connecting to server", "error", err)
		return err
	}

	resp := make([]byte, tcp.messageSize)
	count, err := conn.Read(resp)
	if err != nil {
		return err
	}
	if string(resp[:count]) == ErrConnectionLimitReached.Error() {
		return ErrConnectionLimitReached
	}

	if err := conn.SetDeadline(time.Now().Add(tcp.idleTimeout)); err != nil {
		return err
	}

	tcp.conn = conn
	return nil
}

func (tcp *TCPClient) Close() error {
	if tcp.conn != nil {
		return tcp.conn.Close()
	}
	return nil
}

func (tcp *TCPClient) Send(data []byte) ([]byte, error) {
	if err := tcp.conn.SetDeadline(time.Now().Add(tcp.idleTimeout)); err != nil {
		return nil, err
	}

	if _, err := tcp.conn.Write(data); err != nil {
		return nil, err
	}

	resp := make([]byte, tcp.messageSize)
	count, err := tcp.conn.Read(resp)
	if err != nil {
		if errors.Is(err, os.ErrDeadlineExceeded) {
			panic(err)
		}

		if errors.Is(err, syscall.EPIPE) {
			// TODO: можно вызвать повторное созадние подключения место паники
			panic(err)
		}
		return nil, err
	}

	return resp[:count], nil
}
