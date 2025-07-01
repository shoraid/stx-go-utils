package sliceutil

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSliceUtil_Difference(t *testing.T) {
	t.Run("Integer", func(t *testing.T) {
		tests := []struct {
			name     string
			base     []int
			exclude  []int
			expected []int
		}{
			{
				name:     "No overlap",
				base:     []int{1, 2, 3},
				exclude:  []int{11, 12, 13},
				expected: []int{1, 2, 3},
			},
			{
				name:     "Some overlap",
				base:     []int{1, 2, 3},
				exclude:  []int{2, 4},
				expected: []int{1, 3},
			},
			{
				name:     "All overlap",
				base:     []int{1, 2, 3},
				exclude:  []int{1, 2, 3},
				expected: []int{},
			},
			{
				name:     "With duplicates",
				base:     []int{1, 2, 1, 3, 2},
				exclude:  []int{2},
				expected: []int{1, 1, 3},
			},
			{
				name:     "Empty base",
				base:     []int{},
				exclude:  []int{1, 2},
				expected: []int{},
			},
			{
				name:     "Empty exclude",
				base:     []int{1, 2, 3},
				exclude:  []int{},
				expected: []int{1, 2, 3},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := Difference(tt.base, tt.exclude)
				assert.Equal(t, tt.expected, result)
			})
		}
	})

	t.Run("String", func(t *testing.T) {
		tests := []struct {
			name     string
			base     []string
			exclude  []string
			expected []string
		}{
			{
				name:     "No overlap",
				base:     []string{"a", "b", "c"},
				exclude:  []string{"x", "y", "z"},
				expected: []string{"a", "b", "c"},
			},
			{
				name:     "Some overlap",
				base:     []string{"a", "b", "c"},
				exclude:  []string{"b", "d"},
				expected: []string{"a", "c"},
			},
			{
				name:     "All overlap",
				base:     []string{"a", "b", "c"},
				exclude:  []string{"a", "b", "c"},
				expected: []string{},
			},
			{
				name:     "With duplicates",
				base:     []string{"a", "b", "a", "c", "b"},
				exclude:  []string{"b"},
				expected: []string{"a", "a", "c"},
			},
			{
				name:     "Empty base",
				base:     []string{},
				exclude:  []string{"a", "b"},
				expected: []string{},
			},
			{
				name:     "Empty exclude",
				base:     []string{"a", "b", "c"},
				exclude:  []string{},
				expected: []string{"a", "b", "c"},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := Difference(tt.base, tt.exclude)
				assert.Equal(t, tt.expected, result)
			})
		}
	})
}

func TestSliceUtil_Intersect(t *testing.T) {
	t.Run("Integer", func(t *testing.T) {
		tests := []struct {
			name     string
			source   []int
			target   []int
			expected []int
		}{
			{
				name:     "No overlap",
				source:   []int{1, 2, 3},
				target:   []int{4, 5},
				expected: []int{},
			},
			{
				name:     "Some overlap",
				source:   []int{1, 2, 3},
				target:   []int{2, 4},
				expected: []int{2},
			},
			{
				name:     "All overlap",
				source:   []int{1, 2, 3},
				target:   []int{1, 2, 3},
				expected: []int{1, 2, 3},
			},
			{
				name:     "With duplicates",
				source:   []int{1, 2, 1, 3},
				target:   []int{1, 3},
				expected: []int{1, 1, 3},
			},
			{
				name:     "Empty source",
				source:   []int{},
				target:   []int{1, 2},
				expected: []int{},
			},
			{
				name:     "Empty target",
				source:   []int{1, 2},
				target:   []int{},
				expected: []int{},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assert.Equal(t, tt.expected, Intersect(tt.source, tt.target))
			})
		}
	})

	t.Run("String", func(t *testing.T) {
		tests := []struct {
			name     string
			source   []string
			target   []string
			expected []string
		}{
			{
				name:     "No overlap",
				source:   []string{"a", "b", "c"},
				target:   []string{"x", "y", "z"},
				expected: []string{},
			},
			{
				name:     "Some overlap",
				source:   []string{"a", "b", "c"},
				target:   []string{"b", "d"},
				expected: []string{"b"},
			},
			{
				name:     "All overlap",
				source:   []string{"a", "b", "c"},
				target:   []string{"a", "b", "c"},
				expected: []string{"a", "b", "c"},
			},
			{
				name:     "With duplicates",
				source:   []string{"a", "b", "a", "c", "b"},
				target:   []string{"a", "c"},
				expected: []string{"a", "a", "c"},
			},
			{
				name:     "Empty source",
				source:   []string{},
				target:   []string{"a", "b"},
				expected: []string{},
			},
			{
				name:     "Empty target",
				source:   []string{"a", "b"},
				target:   []string{},
				expected: []string{},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assert.Equal(t, tt.expected, Intersect(tt.source, tt.target))
			})
		}
	})

}

func TestSliceUtil_Map(t *testing.T) {
	type User struct {
		ID   int
		Name string
	}

	tests := []struct {
		name     string
		input    []User
		selector func(User, int) string
		expected []string
	}{
		{
			name: "pluck names",
			input: []User{
				{ID: 1, Name: "Alice"},
				{ID: 2, Name: "Bob"},
			},
			selector: func(u User, _ int) string {
				return u.Name
			},
			expected: []string{"Alice", "Bob"},
		},
		{
			name:  "pluck from empty input",
			input: []User{},
			selector: func(u User, _ int) string {
				return u.Name
			},
			expected: []string{},
		},
		{
			name: "transform IDs to chars",
			input: []User{
				{ID: 3, Name: "Charlie"},
				{ID: 4, Name: "Dana"},
			},
			selector: func(u User, _ int) string {
				return string(rune(u.ID + 64)) // ID 3 → 'C', ID 4 → 'D'
			},
			expected: []string{"C", "D"},
		},
		{
			name: "combine index and name",
			input: []User{
				{ID: 10, Name: "Eve"},
				{ID: 11, Name: "Frank"},
			},
			selector: func(u User, i int) string {
				return fmt.Sprintf("%d-%s", i, u.Name)
			},
			expected: []string{"0-Eve", "1-Frank"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Map(tt.input, tt.selector)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSliceUtil_Unique(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{"With duplicates", []string{"a", "b", "a", "c", "b"}, []string{"a", "b", "c"}},
		{"All unique", []string{"x", "y", "z"}, []string{"x", "y", "z"}},
		{"Empty slice", []string{}, []string{}},
		{"One element", []string{"a"}, []string{"a"}},
		{"All same", []string{"a", "a", "a"}, []string{"a"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Unique(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func BenchmarkSliceUtil_Difference(b *testing.B) {
	type benchCase[T comparable] struct {
		name   string
		size   int
		gen    func(i int) T
		offset int
	}

	b.Run("Integer cases", func(b *testing.B) {
		intCases := []benchCase[int]{
			{"Int-10", 10, func(i int) int { return i }, 5},
			{"Int-100", 100, func(i int) int { return i }, 50},
			{"Int-500", 500, func(i int) int { return i }, 250},
			{"Int-1000", 1000, func(i int) int { return i }, 500},
		}

		for _, tc := range intCases {
			b.Run(tc.name, func(b *testing.B) {
				base := make([]int, tc.size)
				exclude := make([]int, tc.size)

				for i := range tc.size {
					base[i] = tc.gen(i)
					exclude[i] = tc.gen(i + tc.offset)
				}

				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					_ = Difference(base, exclude)
				}
			})
		}
	})

	b.Run("String cases", func(b *testing.B) {
		stringCases := []benchCase[string]{
			{"String-10", 10, func(i int) string { return "item" + strconv.Itoa(i%100) }, 5},
			{"String-100", 100, func(i int) string { return "item" + strconv.Itoa(i%100) }, 50},
			{"String-500", 500, func(i int) string { return "item" + strconv.Itoa(i%100) }, 250},
			{"String-1000", 1000, func(i int) string { return "item" + strconv.Itoa(i%100) }, 500},
		}

		for _, tc := range stringCases {
			b.Run(tc.name, func(b *testing.B) {
				base := make([]string, tc.size)
				exclude := make([]string, tc.size)

				for i := range tc.size {
					base[i] = tc.gen(i)
					exclude[i] = tc.gen(i + tc.offset)
				}

				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					_ = Difference(base, exclude)
				}
			})
		}
	})
}

func BenchmarkSliceUtil_Intersect(b *testing.B) {
	type benchCase[T comparable] struct {
		name   string
		size   int
		gen    func(i int) T
		offset int
	}

	b.Run("Integer cases", func(b *testing.B) {
		intCases := []benchCase[int]{
			{"Int-10", 10, func(i int) int { return i }, 5},
			{"Int-100", 100, func(i int) int { return i }, 50},
			{"Int-500", 500, func(i int) int { return i }, 250},
			{"Int-1000", 1000, func(i int) int { return i }, 500},
		}

		for _, tc := range intCases {
			b.Run(tc.name, func(b *testing.B) {
				source := make([]int, tc.size)
				target := make([]int, tc.size)

				for i := range tc.size {
					source[i] = tc.gen(i)
					target[i] = tc.gen(i + tc.offset)
				}

				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					_ = Intersect(source, target)
				}
			})
		}
	})

	b.Run("String cases", func(b *testing.B) {
		stringCases := []benchCase[string]{
			{"String-10", 10, func(i int) string { return "item" + strconv.Itoa(i%100) }, 5},
			{"String-100", 100, func(i int) string { return "item" + strconv.Itoa(i%100) }, 50},
			{"String-500", 500, func(i int) string { return "item" + strconv.Itoa(i%100) }, 250},
			{"String-1000", 1000, func(i int) string { return "item" + strconv.Itoa(i%100) }, 500},
		}

		for _, tc := range stringCases {
			b.Run(tc.name, func(b *testing.B) {
				source := make([]string, tc.size)
				target := make([]string, tc.size)

				for i := range tc.size {
					source[i] = tc.gen(i)
					target[i] = tc.gen(i + tc.offset)
				}

				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					_ = Intersect(source, target)
				}
			})
		}
	})
}

func BenchmarkSliceUtil_Map(b *testing.B) {
	type User struct {
		ID   int
		Name string
	}

	// Generate sample data
	users := make([]User, 1000)
	for i := range 1000 {
		users[i] = User{
			ID:   i,
			Name: "User" + strconv.Itoa(i),
		}
	}

	cases := []struct {
		name     string
		data     []User
		selector func(User, int) int
	}{
		{
			name:     "Pluck-ID-10",
			data:     users[:10],
			selector: func(u User, _ int) int { return u.ID },
		},
		{
			name:     "Pluck-ID-100",
			data:     users[:100],
			selector: func(u User, _ int) int { return u.ID },
		},
		{
			name:     "Pluck-ID-1000",
			data:     users,
			selector: func(u User, _ int) int { return u.ID },
		},
	}

	b.ResetTimer()

	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = Map(tc.data, tc.selector)
			}
		})
	}
}

func BenchmarkSliceUtil_Unique(b *testing.B) {
	type benchCase[T comparable] struct {
		name string
		size int
		gen  func(i int) T
	}

	b.Run("Integer cases", func(b *testing.B) {
		intCases := []benchCase[int]{
			{"Int-10", 10, func(i int) int { return i % 5 }},
			{"Int-100", 100, func(i int) int { return i % 50 }},
			{"Int-500", 500, func(i int) int { return i % 250 }},
			{"Int-1000", 1000, func(i int) int { return i % 500 }},
		}

		for _, tc := range intCases {
			b.Run(tc.name, func(b *testing.B) {
				input := make([]int, tc.size)
				for i := range input {
					input[i] = tc.gen(i)
				}

				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					_ = Unique(input)
				}
			})
		}
	})

	b.Run("String cases", func(b *testing.B) {
		stringCases := []benchCase[string]{
			{"String-10", 10, func(i int) string { return "item" + strconv.Itoa(i%5) }},
			{"String-100", 100, func(i int) string { return "item" + strconv.Itoa(i%50) }},
			{"String-500", 500, func(i int) string { return "item" + strconv.Itoa(i%250) }},
			{"String-1000", 1000, func(i int) string { return "item" + strconv.Itoa(i%500) }},
		}

		for _, tc := range stringCases {
			b.Run(tc.name, func(b *testing.B) {
				input := make([]string, tc.size)
				for i := range input {
					input[i] = tc.gen(i)
				}

				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					_ = Unique(input)
				}
			})
		}
	})
}
