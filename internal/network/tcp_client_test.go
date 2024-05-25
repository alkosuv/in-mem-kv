package network

import (
	"context"
	"testing"
	"time"
)

func TestTCPClient(t *testing.T) {
	serverCfg := TCPServerConfig{
		Address:               "127.0.0.1:13667",
		MessageSize:           1024,
		MaxConn:               10,
		ConnPoolRetryAttempts: 3,
		ConnPoolRetryTimeout:  1,
		IdleMax:               10,
	}
	server := NewTCPServer(serverCfg)

	var serverErr error
	go func() {
		serverErr = server.Listen(context.Background(), echoHandler)
	}()

	time.Sleep(1 * time.Second) // Give the server time to start

	tests := []struct {
		name           string
		address        string
		messageSize    int
		message        string
		expectResponse string
		expectError    bool
	}{
		{
			name:           "Connect and Send Echo",
			address:        "127.0.0.1:13667",
			messageSize:    1024,
			message:        "hello",
			expectResponse: "hello",
			expectError:    false,
		},
		{
			name:           "Connection Limit Reached",
			address:        "127.0.0.1:13668",
			messageSize:    1024,
			message:        "hello",
			expectResponse: ErrConnectionLimitReached.Error(),
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewTCPClient(tt.address, tt.messageSize)

			err := client.Connect()
			if tt.expectError && err == nil {
				t.Fatalf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if err != nil {
				return // Skip the rest of the test if Connect failed as expected
			}

			resp, err := client.Send([]byte(tt.message))
			if tt.expectError && err == nil {
				t.Fatalf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if string(resp) != tt.expectResponse {
				t.Fatalf("Expected response %q, got %q", tt.expectResponse, resp)
			}

			err = client.Close()
			if err != nil {
				t.Fatalf("Failed to close connection: %v", err)
			}
		})
	}

	if serverErr != nil {
		t.Fatalf("Server error: %v", serverErr)
	}
}
