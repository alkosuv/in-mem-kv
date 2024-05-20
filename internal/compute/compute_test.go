package compute

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestComputeQuery(t *testing.T) {
	ctx := context.Background()
	storage := new(MockStorage)
	compute := New(storage)

	testCases := []struct {
		name    string
		request string
		setFunc func(context.Context, string, string) (string, error)
		getFunc func(context.Context, string) (string, error)
		delFunc func(context.Context, string) (string, error)
		want    string
		wantErr error
	}{
		{
			name:    "SET command",
			request: "SET key value",
			setFunc: func(ctx context.Context, key, value string) (string, error) {
				return "OK", nil
			},
			want: "OK",
		},
		{
			name:    "GET command",
			request: "GET key",
			getFunc: func(ctx context.Context, key string) (string, error) {
				return "value", nil
			},
			want: "value",
		},
		{
			name:    "DEL command",
			request: "DEL key",
			delFunc: func(ctx context.Context, key string) (string, error) {
				return "OK", nil
			},
			want: "OK",
		},
		{
			name:    "Invalid command",
			request: "INVALID command",
			wantErr: ErrInvalidCommand,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storage.SetFunc = tc.setFunc
			storage.GetFunc = tc.getFunc
			storage.DelFunc = tc.delFunc

			res, err := compute.Query(ctx, tc.request)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, res)
			}
		})
	}
}
