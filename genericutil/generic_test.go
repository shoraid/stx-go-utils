package genericutil_test

import (
	"testing"
	"time"

	"github.com/shoraid/stx-go-utils/genericutil"

	"github.com/stretchr/testify/assert"
)

func TestGenericUtil_FirstNonNil(t *testing.T) {
	type testCase[T any] struct {
		name     string
		input    []*T
		expected T
	}

	t.Run("string tests", func(t *testing.T) {
		str := func(s string) *string { return &s }

		stringCases := []testCase[string]{
			{
				name:     "all nil",
				input:    []*string{nil, nil},
				expected: "",
			},
			{
				name:     "first non-nil",
				input:    []*string{str("hello"), str("world")},
				expected: "hello",
			},
			{
				name:     "first is nil, second non-nil",
				input:    []*string{nil, str("second")},
				expected: "second",
			},
		}

		for _, tc := range stringCases {
			t.Run(tc.name, func(t *testing.T) {
				result := genericutil.FirstNonNil(tc.input...)
				assert.Equal(t, tc.expected, result)
			})
		}
	})

	t.Run("int tests", func(t *testing.T) {
		num := func(n int) *int { return &n }

		intCases := []testCase[int]{
			{
				name:     "all nil",
				input:    []*int{nil, nil},
				expected: 0,
			},
			{
				name:     "first non-nil",
				input:    []*int{num(42), num(99)},
				expected: 42,
			},
			{
				name:     "first is nil, second non-nil",
				input:    []*int{nil, num(99)},
				expected: 99,
			},
		}

		for _, tc := range intCases {
			t.Run(tc.name, func(t *testing.T) {
				result := genericutil.FirstNonNil(tc.input...)
				assert.Equal(t, tc.expected, result)
			})
		}
	})
}

func TestGenericUtil_FirstNonZero(t *testing.T) {
	type testCase[T comparable] struct {
		name     string
		input    []T
		expected T
	}

	t.Run("string tests", func(t *testing.T) {
		stringCases := []testCase[string]{
			{
				name:     "all zero (empty strings)",
				input:    []string{"", "", ""},
				expected: "",
			},
			{
				name:     "first non-zero string",
				input:    []string{"hello", "", "world"},
				expected: "hello",
			},
			{
				name:     "non-zero string in the middle",
				input:    []string{"", "world", ""},
				expected: "world",
			},
		}

		for _, tc := range stringCases {
			t.Run(tc.name, func(t *testing.T) {
				result := genericutil.FirstNonZero(tc.input...)
				assert.Equal(t, tc.expected, result)
			})
		}
	})

	t.Run("int tests", func(t *testing.T) {
		intCases := []testCase[int]{
			{
				name:     "all zero (0s)",
				input:    []int{0, 0, 0},
				expected: 0,
			},
			{
				name:     "first non-zero int",
				input:    []int{42, 0, 0},
				expected: 42,
			},
			{
				name:     "non-zero int in the middle",
				input:    []int{0, 100, 0},
				expected: 100,
			},
		}

		for _, tc := range intCases {
			t.Run(tc.name, func(t *testing.T) {
				result := genericutil.FirstNonZero(tc.input...)
				assert.Equal(t, tc.expected, result)
			})
		}
	})
}

func TestGenericUtil_Ptr(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		s := "hello"
		ptr := genericutil.Ptr(s)
		assert.NotNil(t, ptr)
		assert.Equal(t, s, *ptr)
	})

	t.Run("Int", func(t *testing.T) {
		i := 42
		ptr := genericutil.Ptr(i)
		assert.NotNil(t, ptr)
		assert.Equal(t, i, *ptr)
	})

	t.Run("Float64", func(t *testing.T) {
		f := 3.14
		ptr := genericutil.Ptr(f)
		assert.NotNil(t, ptr)
		assert.Equal(t, f, *ptr)
	})

	t.Run("Bool", func(t *testing.T) {
		b := true
		ptr := genericutil.Ptr(b)
		assert.NotNil(t, ptr)
		assert.Equal(t, b, *ptr)
	})

	t.Run("Time", func(t *testing.T) {
		now := time.Now()
		ptr := genericutil.Ptr(now)
		assert.NotNil(t, ptr)
		assert.Equal(t, now, *ptr)
	})

	t.Run("Struct", func(t *testing.T) {
		type Example struct {
			Name string
		}
		ex := Example{Name: "Go"}
		ptr := genericutil.Ptr(ex)
		assert.NotNil(t, ptr)
		assert.Equal(t, ex, *ptr)
	})
}

func BenchmarkGenericUtil_FirstNonNil(b *testing.B) {
	str := func(s string) *string { return &s }

	value := str("data")

	cases := []struct {
		name  string
		input []*string
	}{
		{
			name:  "first non-nil (best case)",
			input: []*string{value, nil, nil, nil, nil},
		},
		{
			name:  "middle non-nil",
			input: []*string{nil, nil, value, nil, nil},
		},
		{
			name:  "last non-nil",
			input: []*string{nil, nil, nil, nil, value},
		},
		{
			name:  "all nil (worst case)",
			input: []*string{nil, nil, nil, nil, nil},
		},
	}

	for _, cs := range cases {
		b.Run(cs.name, func(b *testing.B) {
			for b.Loop() {
				genericutil.FirstNonNil(cs.input...)
			}
		})
	}
}

func BenchmarkGenericUtil_FirstNonZero(b *testing.B) {
	stringCases := []struct {
		name  string
		input []string
	}{
		{
			name:  "first non-zero (best case)",
			input: []string{"hello", "", "", "", ""},
		},
		{
			name:  "middle non-zero",
			input: []string{"", "", "hello", "", ""},
		},
		{
			name:  "last non-zero",
			input: []string{"", "", "", "", "hello"},
		},
		{
			name:  "all zero (worst case)",
			input: []string{"", "", "", "", ""},
		},
	}

	for _, cs := range stringCases {
		b.Run(cs.name, func(b *testing.B) {
			for b.Loop() {
				genericutil.FirstNonZero(cs.input...)
			}
		})
	}
}

func BenchmarkGenericUtil_Ptr(b *testing.B) {
	b.Run("String", func(b *testing.B) {
		s := "benchmark"

		for b.Loop() {
			genericutil.Ptr(s)
		}
	})

	b.Run("Int", func(b *testing.B) {
		data := 123

		for b.Loop() {
			genericutil.Ptr(data)
		}
	})

	b.Run("Struct", func(b *testing.B) {
		type Example struct {
			A int
			B string
		}

		ex := Example{A: 1, B: "data"}

		for b.Loop() {
			genericutil.Ptr(ex)
		}
	})
}
