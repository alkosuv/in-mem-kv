package main

import (
	"context"
	"flag"
	"github.com/alkosuv/in-mem-kv/internal/config"
	"github.com/alkosuv/in-mem-kv/internal/database"
	"github.com/alkosuv/in-mem-kv/internal/database/storage"
	"github.com/alkosuv/in-mem-kv/internal/logger"
	"github.com/alkosuv/in-mem-kv/internal/network"
	"log/slog"
	"os"
)

func main() {
	var ctx context.Context

	configPath := flag.String("configPath", config.DefaultConfigPath, "path to config file")
	flag.Parse()

	cfg, err := config.New(*configPath)
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	logger.Init(cfg.Logging)

	s := storage.New()
	db := database.NewDatabase(s)

	slog.Info("database starting...")
	n := network.NewTCPServer(cfg.Network)
	if err := n.Listen(ctx, db.HandlerQuery); err != nil {
		slog.Error("Failed to listen and serve", "err", err)
		return
	}
}
