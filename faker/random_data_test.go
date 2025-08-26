package faker_test

import (
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/shoraid/stx-go-utils/faker"

	"github.com/stretchr/testify/assert"
)

func TestFaker_PickRandom(t *testing.T) {
	tests := []struct {
		name     string
		elements []any
	}{
		{
			name:     "Pick from strings",
			elements: []any{"a", "b", "c", "d"},
		},
		{
			name:     "Pick from numbers",
			elements: []any{1, 2, 3, 4, 5},
		},
		{
			name:     "Pick from booleans",
			elements: []any{true, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			random := faker.PickRandom(tt.elements...)
			assert.Contains(t, tt.elements, random)
		})
	}
}

func TestFaker_RandBool(t *testing.T) {
	t.Run("should generate random boolean values", func(t *testing.T) {
		var trueCount, falseCount int

		for range 100 {
			if faker.RandBool() {
				trueCount++
			} else {
				falseCount++
			}
		}

		assert.Greater(t, trueCount, 0, "should return true at least once")
		assert.Greater(t, falseCount, 0, "should return false at least once")
	})
}

func TestFaker_RandBoolPtr(t *testing.T) {
	t.Run("should generate pointer to random boolean", func(t *testing.T) {
		ptr := faker.RandBoolPtr()

		assert.NotNil(t, ptr, "Expected pointer, got nil")
		assert.True(t, *ptr == true || *ptr == false, "Expected true or false")
	})
}

func TestFaker_RandInt(t *testing.T) {
	tests := []struct {
		name     string
		min      int
		max      int
		validate func(t *testing.T, result int, min int, max int)
	}{
		{
			name: "in-range 0-10",
			min:  0,
			max:  10,
			validate: func(t *testing.T, result, min, max int) {
				assert.GreaterOrEqual(t, result, min)
				assert.LessOrEqual(t, result, max)
			},
		},
		{
			name: "in-range -10 to 10",
			min:  -10,
			max:  10,
			validate: func(t *testing.T, result, min, max int) {
				assert.GreaterOrEqual(t, result, min)
				assert.LessOrEqual(t, result, max)
			},
		},
		{
			name: "same min and max",
			min:  5,
			max:  5,
			validate: func(t *testing.T, result, min, _ int) {
				assert.Equal(t, min, result)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := faker.RandInt(tt.min, tt.max)
			tt.validate(t, result, tt.min, tt.max)
		})
	}
}

func TestFaker_RandIntPtr(t *testing.T) {
	tests := []struct {
		name string
		min  int
		max  int
	}{
		{name: "range 1-3", min: 1, max: 3},
		{name: "range -5 to 5", min: -5, max: 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ptr := faker.RandIntPtr(tt.min, tt.max)
			assert.NotNil(t, ptr)
			assert.GreaterOrEqual(t, *ptr, tt.min)
			assert.LessOrEqual(t, *ptr, tt.max)
		})
	}
}

func TestFaker_RandSentence(t *testing.T) {
	tests := []struct {
		name      string
		wordCount int
	}{
		{name: "1-word", wordCount: 1},
		{name: "5-words", wordCount: 5},
		{name: "10-words", wordCount: 10},
		{name: "50-words", wordCount: 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sentence := faker.RandSentence(tt.wordCount)
			words := strings.Fields(sentence)

			assert.Len(t, words, tt.wordCount, "should generate correct number of words")
		})
	}
}

func TestFaker_RandSentencePtr(t *testing.T) {
	t.Run("10-words", func(t *testing.T) {
		ptr := faker.RandSentencePtr(10)
		assert.NotNil(t, ptr, "should not return nil")

		words := strings.Fields(*ptr)
		assert.Len(t, words, 10, "should contain 10 words")
	})
}

func TestFaker_RandString(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{name: "length-0", length: 0},
		{name: "length-1", length: 1},
		{name: "length-10", length: 10},
		{name: "length-50", length: 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := faker.RandString(tt.length)
			assert.Equal(t, tt.length, utf8.RuneCountInString(result), "string should have correct length")
		})
	}
}

func TestFaker_RandStringPtr(t *testing.T) {
	t.Run("length-16", func(t *testing.T) {
		ptr := faker.RandStringPtr(16)
		assert.NotNil(t, ptr, "should not be nil")
		assert.Equal(t, 16, utf8.RuneCountInString(*ptr), "should have correct length")
	})
}

func TestFaker_RandTime(t *testing.T) {
	now := time.Now()
	later := now.Add(24 * time.Hour)

	tests := []struct {
		name  string
		start time.Time
		end   time.Time
	}{
		{
			name:  "Valid range",
			start: now,
			end:   later,
		},
		{
			name:  "Reverse range",
			start: later,
			end:   now,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := faker.RandTime(tt.start, tt.end)
			start := tt.start
			end := tt.end
			if start.After(end) {
				start, end = end, start
			}
			assert.True(t, result.Equal(start) || result.After(start), "should be >= start")
			assert.True(t, result.Equal(end) || result.Before(end), "should be <= end")
		})
	}
}

func TestFaker_RandTimePtr(t *testing.T) {
	start := time.Now()
	end := start.Add(1 * time.Hour)

	ptr := faker.RandTimePtr(start, end)
	assert.NotNil(t, ptr, "pointer should not be nil")
	assert.True(t, (*ptr).After(start) || (*ptr).Equal(start))
	assert.True(t, (*ptr).Before(end) || (*ptr).Equal(end))
}

func TestFaker_RandURL(t *testing.T) {
	t.Run("should generate a valid random URL", func(t *testing.T) {
		url := faker.RandURL()

		assert.Contains(t, url, "https://", "URL should start with https://")
		assert.Contains(t, url, ".", "URL should contain a domain")
	})
}

func TestFaker_RandURLPtr(t *testing.T) {
	t.Run("should generate pointer to random URL", func(t *testing.T) {
		ptr := faker.RandURLPtr()

		assert.NotNil(t, ptr, "Expected pointer, got nil")
		assert.Contains(t, *ptr, "https://", "URL should start with https://")
	})
}

func TestFaker_UUID(t *testing.T) {
	tests := []struct {
		name string
		num  int
	}{
		{
			name: "generate 1 UUID",
			num:  1,
		},
		{
			name: "generate 100 UUIDs and ensure uniqueness",
			num:  100,
		},
		{
			name: "generate 1000 UUIDs and ensure uniqueness",
			num:  1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seen := make(map[string]bool, tt.num)

			for i := 0; i < tt.num; i++ {
				id := faker.UUID()

				parsed, err := uuid.Parse(id)
				assert.NoError(t, err, "should be valid UUID")
				assert.Equal(t, uuid.Version(7), parsed.Version(), "should be UUIDv7")

				_, exists := seen[id]
				assert.False(t, exists, "UUID must be unique")

				seen[id] = true
			}
		})
	}
}

func TestFaker_UUIDPtr(t *testing.T) {
	tests := []struct {
		name string
		num  int
	}{
		{
			name: "generate 1 UUID",
			num:  1,
		},
		{
			name: "generate 100 UUIDs and ensure uniqueness",
			num:  100,
		},
		{
			name: "generate 1000 UUIDs and ensure uniqueness",
			num:  1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seen := make(map[string]bool, tt.num)

			for i := 0; i < tt.num; i++ {
				id := faker.UUIDPtr()

				parsed, err := uuid.Parse(*id)
				assert.NoError(t, err, "should be valid UUID")
				assert.Equal(t, uuid.Version(7), parsed.Version(), "should be UUIDv7")

				_, exists := seen[*id]
				assert.False(t, exists, "UUID must be unique")

				seen[*id] = true
			}
		})
	}
}

func BenchmarkFaker_PickRandom(b *testing.B) {
	elements := []any{"apple", 123, true, 4.5, "banana", struct{}{}, []int{1, 2, 3}}

	for b.Loop() {
		faker.PickRandom(elements...)
	}
}

func BenchmarkFaker_RandBool(b *testing.B) {
	for b.Loop() {
		faker.RandBool()
	}
}

func BenchmarkFaker_RandBoolPtr(b *testing.B) {
	for b.Loop() {
		faker.RandBoolPtr()
	}
}

func BenchmarkFaker_RandInt(b *testing.B) {
	for b.Loop() {
		faker.RandInt(1, 100)
	}
}

func BenchmarkFaker_RandIntPtr(b *testing.B) {
	for b.Loop() {
		faker.RandIntPtr(1, 100)
	}
}

func BenchmarkFaker_RandSentence(b *testing.B) {
	for b.Loop() {
		faker.RandSentence(10)
	}
}

func BenchmarkFaker_RandSentencePtr(b *testing.B) {
	for b.Loop() {
		faker.RandSentencePtr(10)
	}
}

func BenchmarkFaker_RandString(b *testing.B) {
	for b.Loop() {
		faker.RandString(16)
	}
}

func BenchmarkFaker_RandStringPtr(b *testing.B) {
	for b.Loop() {
		faker.RandStringPtr(16)
	}
}

func BenchmarkFaker_RandTime(b *testing.B) {
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)

	for b.Loop() {
		faker.RandTime(start, end)
	}
}

func BenchmarkFaker_RandTimePtr(b *testing.B) {
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)

	for b.Loop() {
		faker.RandTimePtr(start, end)
	}
}

func BenchmarkFaker_RandURL(b *testing.B) {
	for b.Loop() {
		faker.RandURL()
	}
}

func BenchmarkFaker_RandURLPtr(b *testing.B) {
	for b.Loop() {
		faker.RandURLPtr()
	}
}

func BenchmarkFaker_UUID(b *testing.B) {
	for b.Loop() {
		faker.UUID()
	}
}

func BenchmarkFaker_UUIDPtr(b *testing.B) {
	for b.Loop() {
		faker.UUIDPtr()
	}
}
