package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/alkosuv/in-mem-kv/internal/logger"
	"github.com/alkosuv/in-mem-kv/internal/network"
	"log/slog"
	"os"
)

func main() {

	var (
		address     = flag.String("address", "127.0.0.1:13666", "address to listen")
		messageSize = flag.Int("message-size", 1024, "message size")
	)
	flag.Parse()

	logger.DefaultInit()

	client := network.NewTCPClient(*address, *messageSize)
	if err := client.Connect(); err != nil {
		slog.Error("Failed to connect to server", "error", err)
		os.Exit(1)
	}
	defer client.Close()

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		message, err := reader.ReadBytes('\n')
		if err != nil {
			slog.Error("read request", "error", err, "query", message)
			continue
		}

		resp, err := client.Send(message)
		if err != nil {
			slog.Error("send request", "error", err, "query", message)
			continue
		}

		fmt.Printf("< %s\n", string(resp))
	}
}
