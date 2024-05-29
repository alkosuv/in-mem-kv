package tools

import (
	"testing"
)

func TestSemaphore(t *testing.T) {
	tests := []struct {
		name       string
		operations []string
		limit      int
		expected   []bool
	}{
		{
			name:       "Acquire and Release",
			operations: []string{"Acquire", "Acquire", "Release", "Acquire"},
			limit:      2,
			expected:   []bool{true, true, true, true},
		},
		{
			name:       "TryAcquire success",
			operations: []string{"TryAcquire", "TryAcquire"},
			limit:      2,
			expected:   []bool{true, true},
		},
		{
			name:       "TryAcquire limit reached",
			operations: []string{"TryAcquire", "TryAcquire", "TryAcquire"},
			limit:      2,
			expected:   []bool{true, true, false},
		},
		{
			name:       "TryAcquire after Release",
			operations: []string{"TryAcquire", "TryAcquire", "Release", "TryAcquire"},
			limit:      2,
			expected:   []bool{true, true, true, true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sem := NewSemaphore(tt.limit)
			var results []bool

			for _, op := range tt.operations {
				switch op {
				case "Acquire":
					sem.Acquire()
					results = append(results, true) // Always true for Acquire since it blocks until it can acquire
				case "Release":
					sem.Release()
					results = append(results, true) // Always true for Release since it always succeeds
				case "TryAcquire":
					result := sem.TryAcquire()
					results = append(results, result)
				}
			}

			for i, result := range results {
				if result != tt.expected[i] {
					t.Fatalf("Expected result for operation %s to be %v, got %v", tt.operations[i], tt.expected[i], result)
				}
			}
		})
	}
}
