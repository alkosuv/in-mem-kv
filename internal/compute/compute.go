package compute

import (
	"context"
	"log/slog"
)

type storage interface {
	Set(ctx context.Context, key, value string) (string, error)
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) (string, error)
}

type Compute struct {
	s storage
}

func New(s storage) *Compute {
	return &Compute{s: s}
}

// Query ...
func (c *Compute) Query(ctx context.Context, req string) (string, error) {
	operationID, ok := ctx.Value("operation_id").(string)
	if !ok {
		slog.Info("ProcessingQuery", "error", "empty operation_id")
	}

	query, err := processingQuery(ctx, req)
	if err != nil {
		slog.Debug("processing query", "operation_id", operationID, "error", err)
		return "", err
	}

	switch query.Command {
	case SetCommand:
		return c.s.Set(ctx, query.Key, query.Value)
	case GetCommand:
		return c.s.Get(ctx, query.Key)
	case DelCommand:
		return c.s.Del(ctx, query.Key)
	default:
		slog.Error("compute query invalid command")
		return "", ErrInvalidCommand
	}
}
