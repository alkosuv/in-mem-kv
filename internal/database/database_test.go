package database

import (
	"context"
	"github.com/alkosuv/in-mem-kv/internal/database/compute"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MockStorage struct {
	SetFunc func(ctx context.Context, key, value string) (string, error)
	GetFunc func(ctx context.Context, key string) (string, error)
	DelFunc func(ctx context.Context, key string) (string, error)
}

func (m *MockStorage) Set(ctx context.Context, key, value string) (string, error) {
	return m.SetFunc(ctx, key, value)
}

func (m *MockStorage) Get(ctx context.Context, key string) (string, error) {
	return m.GetFunc(ctx, key)
}

func (m *MockStorage) Del(ctx context.Context, key string) (string, error) {
	return m.DelFunc(ctx, key)
}

func TestDatabase_HandlerQuery(t *testing.T) {
	storage := new(MockStorage)

	testCases := []struct {
		name    string
		request string
		ctx     context.Context
		setFunc func(context.Context, string, string) (string, error)
		getFunc func(context.Context, string) (string, error)
		delFunc func(context.Context, string) (string, error)
		want    string
		wantErr error
	}{
		{
			name:    "SET command",
			request: "SET key value",
			ctx:     context.WithValue(context.Background(), "operation_id", "test"),
			setFunc: func(ctx context.Context, key, value string) (string, error) {
				return "OK", nil
			},
			want: "OK",
		},
		{
			name:    "GET command",
			request: "GET key",
			ctx:     context.WithValue(context.Background(), "operation_id", "test"),
			getFunc: func(ctx context.Context, key string) (string, error) {
				return "value", nil
			},
			want: "value",
		},
		{
			name:    "DEL command",
			request: "DEL key",
			ctx:     context.WithValue(context.Background(), "operation_id", "test"),
			delFunc: func(ctx context.Context, key string) (string, error) {
				return "OK", nil
			},
			want: "OK",
		},
		{
			name:    "Invalid command",
			request: "INVALID command",
			ctx:     context.WithValue(context.Background(), "operation_id", "test"),
			wantErr: compute.ErrInvalidCommand,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storage.SetFunc = tc.setFunc
			storage.GetFunc = tc.getFunc
			storage.DelFunc = tc.delFunc

			db := NewDatabase(storage)
			resp, err := db.HandlerQuery(tc.ctx, []byte(tc.request))

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, string(resp))
			}
		})
	}
}
