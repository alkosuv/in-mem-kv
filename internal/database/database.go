package database

import (
	"context"
	"github.com/alkosuv/in-mem-kv/internal/database/compute"
	"log/slog"
)

type storage interface {
	Set(ctx context.Context, key, value string) (string, error)
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) (string, error)
}

type Database struct {
	storage storage
}

func NewDatabase(storage storage) *Database {
	return &Database{
		storage: storage,
	}
}

func (db *Database) HandlerQuery(ctx context.Context, req []byte) ([]byte, error) {
	operationID, ok := ctx.Value("operation_id").(string)
	if !ok {
		slog.Info("ProcessingQuery", "error", "empty operation_id")
	}

	query, err := compute.Compute(ctx, string(req))
	if err != nil {
		slog.Debug("processing query", "operation_id", operationID, "error", err)
		return nil, err
	}

	var resp string
	switch query.Command {
	case compute.SetCommand:
		resp, err = db.storage.Set(ctx, query.Key, query.Value)
	case compute.GetCommand:
		resp, err = db.storage.Get(ctx, query.Key)
	case compute.DelCommand:
		resp, err = db.storage.Del(ctx, query.Key)
	default:
		return nil, compute.ErrInvalidCommand
	}

	return []byte(resp), err
}
