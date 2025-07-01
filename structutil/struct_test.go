package structutil

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/shoraid/stx-go-utils/apperror"
	"github.com/stretchr/testify/assert"
)

func TestStructUtil_BindJSON(t *testing.T) {
	type TestPayload struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	type testCase struct {
		name        string
		body        io.Reader
		expected    TestPayload
		expectError error
	}

	tests := []testCase{
		{
			name: "valid JSON",
			body: bytes.NewBufferString(`{"name":"John","age":30}`),
			expected: TestPayload{
				Name: "John",
				Age:  30,
			},
			expectError: nil,
		},
		{
			name: "unknown field",
			body: bytes.NewBufferString(`{"name":"Alice","age":25,"extra":"field"}`),
			expected: TestPayload{
				Name: "Alice",
				Age:  25,
			},
			expectError: assert.AnError,
		},
		{
			name:        "nil body",
			body:        nil,
			expected:    TestPayload{},
			expectError: apperror.Err400InvalidBody,
		},
		{
			name:        "invalid JSON format",
			body:        bytes.NewBufferString(`{invalid json}`),
			expected:    TestPayload{},
			expectError: assert.AnError,
		},
		{
			name:        "type mismatch",
			body:        bytes.NewBufferString(`{"name":"Alice","age":"old}`),
			expected:    TestPayload{},
			expectError: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{
				Body: toReadCloser(tt.body),
			}

			var result TestPayload
			err := BindJSON(req, &result)

			if tt.expectError != nil {
				assert.Error(t, err)
				if tt.expectError != assert.AnError {
					assert.ErrorIs(t, err, tt.expectError)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func BenchmarkStructUtil_BindJSON(b *testing.B) {
	type TestPayload struct {
		Name     string `json:"name"`
		Age      int    `json:"age"`
		IsActive bool   `json:"isActive"`
	}

	for i := 0; i < b.N; i++ {
		body := bytes.NewBufferString(`{"name":"John","age":30, "isActive": true}`)

		req := &http.Request{
			Body: toReadCloser(body),
		}

		var payload TestPayload
		err := BindJSON(req, &payload)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

// Helper to convert io.Reader to io.ReadCloser
func toReadCloser(r io.Reader) io.ReadCloser {
	if r == nil {
		return nil
	}
	return io.NopCloser(r)
}
