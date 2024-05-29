package compute

import (
	"context"
	"log/slog"
	"strings"
)

// Compute функция обрабатывает полученный запрос
func Compute(ctx context.Context, req string) (Query, error) {
	tokens, err := parsingQuery(req)
	if err != nil {
		return Query{}, err
	}

	slog.Debug("processing query", "operation_id", ctx.Value("operation_id"), "tokens", tokens)

	if err := analyzeQuery(tokens); err != nil {
		return Query{}, err
	}

	query := Query{
		Command: StringToCommand[strings.ToUpper(tokens[commandTokenIndex])],
		Key:     tokens[keyTokenIndex],
	}
	// Проверка нужна только для команды SET
	if len(tokens) == maxLenToken {
		query.Value = tokens[valueTokenIndex]
	}

	return query, nil
}
