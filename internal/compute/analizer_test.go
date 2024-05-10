package compute

import (
	"testing"
)

func TestAnalyzeQuery(t *testing.T) {
	testCases := []struct {
		name     string
		tokens   tokens
		expected error
	}{
		{
			name:     "Valid GET command",
			tokens:   tokens([]string{"GET", "key"}),
			expected: nil,
		},
		{
			name:     "Valid SET command",
			tokens:   tokens([]string{"SET", "key", "value"}),
			expected: nil,
		},
		{
			name:     "Valid DEL command",
			tokens:   tokens([]string{"DEL", "key"}),
			expected: nil,
		},
		{
			name:     "Invalid command",
			tokens:   tokens([]string{"INVALID", "command"}),
			expected: ErrInvalidCommand,
		},
		{
			name:     "Extra token for GET command",
			tokens:   tokens([]string{"GET", "key", "extra"}),
			expected: ErrInvalidNumberArgument,
		},
		{
			name:     "Extra token for DEL command",
			tokens:   tokens([]string{"DEL", "key", "extra"}),
			expected: ErrInvalidNumberArgument,
		},
		{
			name:     "Key is several words",
			tokens:   []string{"SET", "key with spaces", "value"},
			expected: ErrInvalidSymbol,
		},
		{
			name:     "Russian letters in key",
			tokens:   []string{"SET", "ключ", "value"},
			expected: ErrInvalidSymbol,
		},
		{
			name:     "Invalid use of special characters",
			tokens:   []string{"SET", "ke+y", "value"},
			expected: ErrInvalidSymbol,
		},
		{
			name:     "Empty tokens",
			tokens:   tokens([]string{}),
			expected: ErrInvalidQuery,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := analyzeQuery(tc.tokens)

			if err != tc.expected {
				t.Errorf("analyzeQuery(%v) = %v; expected %v", tc.tokens, err, tc.expected)
			}
		})
	}
}
