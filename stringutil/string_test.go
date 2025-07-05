package stringutil_test

import (
	"testing"

	"github.com/shoraid/stx-go-utils/stringutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringUtil_GenerateUUID(t *testing.T) {
	id, err := stringutil.GenerateUUID()

	require.NoError(t, err)
	require.NotEmpty(t, id)
}

func TestStringUtil_ToSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"TestCase", "test_case"},
		{"Already_Snake", "already_snake"},
		{"TestID", "test_id"},
		{"HTTPServer", "http_server"},
		{"snake_case_test", "snake_case_test"},
	}

	for _, test := range tests {
		actual := stringutil.ToSnakeCase(test.input)

		assert.Equal(t, test.expected, actual, "ToSnakeCase should convert %s to %s", test.input, test.expected)
	}
}

func BenchmarkStringUtil_GenerateUUID(b *testing.B) {
	for b.Loop() {
		stringutil.GenerateUUID()
	}
}

func BenchmarkStringUtil_ToSnakeCase(b *testing.B) {
	input := "TestCaseWithHTTPRequest"
	for b.Loop() {
		stringutil.ToSnakeCase(input)
	}
}
