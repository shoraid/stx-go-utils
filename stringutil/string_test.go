package stringutil_test

import (
	"testing"

	"github.com/shoraid/stx-go-utils/stringutil"

	"github.com/stretchr/testify/assert"
)

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

func BenchmarkStringUtil_ToSnakeCase(b *testing.B) {
	input := "TestCaseWithHTTPRequest"
	for i := 0; i < b.N; i++ {
		_ = stringutil.ToSnakeCase(input)
	}
}
