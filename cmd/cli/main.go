package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/alkosuv/in-mem-kv/internal/compute"
	"github.com/alkosuv/in-mem-kv/internal/storage"
	"github.com/alkosuv/in-mem-kv/internal/uuid"
	"log/slog"
	"os"
)

func main() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, opts)))

	s := storage.New()
	cmp := compute.New(s)

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		req, err := reader.ReadString('\n')
		if err != nil {
			slog.Error("read request", "error", err, "query", req)
			continue
		}

		if req == "q\n" {
			return
		}

		uuid, err := uuid.Generate()
		if err != nil {
			slog.Error("generate uuid", "error", err)
			continue
		}

		ctx := context.WithValue(context.Background(), "operation_id", uuid)
		value, err := cmp.Query(ctx, req)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("> %s\n", value)
		}

	}

}
