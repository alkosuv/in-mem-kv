package compute

import (
	"context"
	"reflect"
	"testing"
)

func TestProcessingCompute(t *testing.T) {
	// Определяем контекст
	ctx := context.Background()

	// Создаем тестовые случаи
	testCases := []struct {
		name      string
		input     string
		expected  Query
		expectErr bool
	}{
		{
			name:  "SET command",
			input: "SET key value",
			expected: Query{
				Command: SetCommand,
				Key:     "key",
				Value:   "value",
			},
			expectErr: false,
		},
		{
			name:  "GET command",
			input: "GET key",
			expected: Query{
				Command: GetCommand,
				Key:     "key",
			},
			expectErr: false,
		},
		{
			name:  "DEL command",
			input: "DEL key",
			expected: Query{
				Command: DelCommand,
				Key:     "key",
			},
			expectErr: false,
		},
		{
			name:      "Invalid command",
			input:     "INVALID_COMMAND key value",
			expected:  Query{},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Compute(ctx, tc.input)
			if (err != nil) != tc.expectErr {
				t.Errorf("processingQuery(%q) вернул ошибку: %v, ожидается ошибка: %v", tc.input, err, tc.expectErr)
				return
			}
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("processingQuery(%q) = %v; ожидается %v", tc.input, result, tc.expected)
			}
		})
	}
}
