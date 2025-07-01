package sliceutil

import (
	"slices"
)

// Difference returns elements from the 'base' slice that are not present in the 'exclude' slice.
// It compares each item in 'base' and includes it in the result only if it is not found in 'exclude'.
//
// Note: Duplicates from 'base' will be preserved if they are not in 'exclude'.
//
// Example:
//
//	base := []int{1, 2, 3, 4, 5}
//	exclude := []int{2, 4}
//
//	result := Difference(base, exclude)
//	// result: []int{1, 3, 5}
func Difference[T comparable](base []T, exclude []T) []T {
	result := make([]T, 0)

	for _, item := range base {
		if !slices.Contains(exclude, item) {
			result = append(result, item)
		}
	}

	return result
}

// Intersect returns the intersection of two slices.
// It compares the elements in 'source' and checks which ones also exist in 'target'.
// The result includes all matching elements in the order they appear in 'source'.
//
// Note: Duplicates are not removed. If an item appears multiple times in 'source' and exists in 'target',
// it will be included multiple times in the result.
//
// Example:
//
//	source := []int{1, 2, 2, 3, 4}
//	target := []int{2, 4, 6}
//
//	result := Intersect(source, target)
//	// result: []int{2, 2, 4}
func Intersect[T comparable](source, target []T) []T {
	result := make([]T, 0)

	for _, item := range source {
		if slices.Contains(target, item) {
			result = append(result, item)
		}
	}

	return result
}

// Map applies a transformation function to each element of a slice and
// returns a new slice containing the results.
//
// It takes a slice of type T and a selector function that maps each
// element of type T to a result of type R. The resulting slice will
// have the same length and contain the transformed elements.
//
// Example:
//
//	names := Map(users, func(u User, _ int) []string {
//		return u.Name
//	})
//	// names: []string{"Alice", "Bob", "Charlie"}
//
//	names := Map(users, func(u User, index int) []string {
//		return fmt.Sprintf("%d-%s", index, u.Name)
//	})
//	// names: []string{"0-Alice", "1-Bob", "2-Charlie"}
//
// This is a generic alternative to manual loops for transforming data.
func Map[T any, R any](items []T, selector func(item T, index int) R) []R {
	result := make([]R, 0, len(items))

	for i, item := range items {
		result = append(result, selector(item, i))
	}

	return result
}

// Unique returns a new slice containing only the unique elements from the input slice.
// It removes duplicates while preserving the order of first occurrence.
//
// This function works with any slice of comparable types (e.g., string, int, float64).
//
// Note: If an element appears multiple times, only its first occurrence will be included in the result.
//
// Example:
//
//	input := []string{"a", "b", "a", "c", "b"}
//
//	result := Unique(input)
//	// result: []string{"a", "b", "c"}
func Unique[T comparable](input []T) []T {
	seen := make(map[T]struct{})
	result := make([]T, 0)

	for _, val := range input {
		if _, exists := seen[val]; !exists {
			seen[val] = struct{}{}
			result = append(result, val)
		}
	}
	return result
}
