package network

import (
	"context"
	"io"
	"net"
	"testing"
	"time"
)

func echoHandler(ctx context.Context, req []byte) ([]byte, error) {
	return req, nil
}

func TestTCPServer_Listen(t *testing.T) {
	tests := []struct {
		name           string
		cfg            TCPServerConfig
		message        string
		expectResponse string
		expectError    error
	}{
		{
			name: "Echo message",
			cfg: TCPServerConfig{
				Address:               "127.0.0.1:13661",
				MessageSize:           1024,
				MaxConn:               10,
				ConnPoolRetryAttempts: 3,
				ConnPoolRetryTimeout:  1,
				IdleMax:               1000,
			},
			message:        "hello",
			expectResponse: "hello",
			expectError:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewTCPServer(tt.cfg)

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			// Start the server in a separate goroutine
			go func() {
				if err := server.Listen(ctx, echoHandler); err != nil {
					t.Fatalf("Server failed to start: %v", err)
				}
			}()

			// Connect to the server as a client
			conn, err := net.Dial("tcp", tt.cfg.Address)
			if err != nil {
				t.Fatalf("Failed to connect to server: %v", err)
			}
			defer conn.Close()

			// connection check
			resp := make([]byte, len(tt.message))
			if _, err = conn.Read(resp); err != nil {
				t.Fatalf("Failed to read from server: %v", err)
			}

			// Test sending and receiving data
			if _, err = conn.Write([]byte(tt.message)); err != nil {
				t.Fatalf("Failed to send data: %v", err)
			}

			count, err := conn.Read(resp)
			if err != nil {
				t.Fatalf("Failed to read response: %v", err)
			}

			if string(resp[:count]) != tt.expectResponse {
				t.Fatalf("Expected response %q, got %q", tt.expectResponse, resp)
			}
		})
	}
}

func TestTCPServer_Listen_ExceedConnectionLimit(t *testing.T) {
	cfg := TCPServerConfig{
		Address:               "127.0.0.1:13664",
		MessageSize:           1024,
		MaxConn:               2,
		ConnPoolRetryAttempts: 3,
		ConnPoolRetryTimeout:  1,
		IdleMax:               10,
	}
	expectResponse := ErrConnectionLimitReached.Error()

	server := NewTCPServer(cfg)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Start the server in a separate goroutine
	go func() {
		if err := server.Listen(ctx, echoHandler); err != nil {
			t.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Connect to the server as a client
	conn, err := net.Dial("tcp", cfg.Address)
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Open max connections
	var conns []net.Conn
	for i := 0; i < cfg.MaxConn; i++ {
		c, err := net.Dial("tcp", cfg.Address)
		if err != nil {
			t.Fatalf("Failed to connect to server: %v", err)
		}
		conns = append(conns, c)
	}
	defer func() {
		for _, c := range conns {
			c.Close()
		}
	}()

	// Attempt to exceed connection limit
	extraConn, err := net.Dial("tcp", cfg.Address)
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer extraConn.Close()

	resp := make([]byte, len(expectResponse))
	_, err = io.ReadFull(extraConn, resp)
	if err != nil {
		t.Fatalf("Failed to read response from extra connection: %v", err)
	}

	if string(resp) != expectResponse {
		t.Fatalf("Expected response %q, got %q", expectResponse, resp)
	}
}
