package stringutil

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestStringUtil_IsValidUUID(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "valid UUID v4",
			input:    uuid.NewString(), // generate valid UUID v4
			expected: true,
		},
		{
			name:     "valid UUID v7",
			input:    uuid.Must(uuid.NewV7()).String(), // generate valid UUID v7
			expected: true,
		},
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "malformed UUID",
			input:    "1234-invalid-uuid",
			expected: false,
		},
		{
			name:     "almost valid UUID",
			input:    "550e8400-e29b-41d4-a716-44665544000", // missing one digit
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := IsValidUUID(tt.input)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func BenchmarkIsValidUUID_Valid(b *testing.B) {
	validUUID := uuid.NewString()

	for b.Loop() {
		IsValidUUID(validUUID)
	}
}

func BenchmarkIsValidUUID_Invalid(b *testing.B) {
	invalidUUID := "not-a-valid-uuid"

	for b.Loop() {
		IsValidUUID(invalidUUID)
	}
}

func BenchmarkIsValidUUID_Valid2(b *testing.B) {
	validUUID := uuid.NewString()

	for b.Loop() {
		IsValidUUID(validUUID)
	}
}

func BenchmarkIsValidUUID_Invalid2(b *testing.B) {
	invalidUUID := "not-a-valid-uuid"

	for b.Loop() {
		IsValidUUID(invalidUUID)
	}
}
