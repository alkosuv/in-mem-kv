package logger

import (
	"io"
	"log/slog"
	"os"
)

const defaultLogLevel = slog.LevelInfo

type Config struct {
	Level  string `yaml:"level"`
	Output string `yaml:"output"`
}

func DefaultInit() {
	Init(Config{})
}

func Init(cfg Config) {
	var err error

	level := new(slog.Level)
	if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
		*level = defaultLogLevel
	}

	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	var w io.Writer = os.Stdout
	if cfg.Output != "" {
		w, err = os.OpenFile(cfg.Output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			slog.Error("Failed to open log file", "file", cfg.Output, "error", err)
		}
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(w, opts)))
}
