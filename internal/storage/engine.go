package storage

import "context"

const minNumberMap int = 100_000

type Engine struct {
	m map[string]string
}

func New() *Engine {
	return &Engine{m: make(map[string]string, minNumberMap)}
}

func (e *Engine) Set(_ context.Context, key, value string) (string, error) {
	e.m[key] = value
	return "ok", nil
}
func (e *Engine) Get(_ context.Context, key string) (string, error) {
	if value, ok := e.m[key]; ok {
		return value, nil
	}
	return "", nil
}

func (e *Engine) Del(_ context.Context, key string) (string, error) {
	delete(e.m, key)
	return "ok", nil
}
