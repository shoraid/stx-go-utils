package structutil_test

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/shoraid/stx-go-utils/structutil"
	"github.com/stretchr/testify/assert"
)

type Sample struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type errorReader struct{}

func (e errorReader) Read(_ []byte) (int, error) {
	return 0, errors.New("read error")
}

func TestStructUtil_DecodeJSON(t *testing.T) {
	tests := []struct {
		name        string
		input       io.Reader
		target      any
		expected    any
		expectError bool
	}{
		{
			name:        "valid JSON",
			input:       strings.NewReader(`{"name":"Alice","age":30}`),
			target:      &Sample{},
			expected:    &Sample{Name: "Alice", Age: 30},
			expectError: false,
		},
		{
			name:        "invalid JSON syntax",
			input:       strings.NewReader(`{"name":"Alice",`),
			target:      &Sample{},
			expected:    &Sample{},
			expectError: true,
		},
		{
			name:        "type mismatch",
			input:       strings.NewReader(`{"name":true,"age":"old"}`),
			target:      &Sample{},
			expected:    &Sample{},
			expectError: true,
		},
		{
			name:        "empty input",
			input:       strings.NewReader(``),
			target:      &Sample{},
			expected:    &Sample{},
			expectError: true,
		},
		{
			name:        "read error",
			input:       errorReader{}, // simulate error from io.Reader
			target:      &Sample{},
			expected:    &Sample{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := structutil.DecodeJSON(tt.input, tt.target)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, tt.target)
			}
		})
	}
}

func BenchmarkStructUtil_DecodeJSON(b *testing.B) {
	var jsonInput = `{"name":"Alice","age":30}`

	for i := 0; i < b.N; i++ {
		reader := strings.NewReader(jsonInput)
		var s Sample
		err := structutil.DecodeJSON(reader, &s)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}
