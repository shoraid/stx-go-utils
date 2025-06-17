package sliceutil

import (
	"slices"
)

/*
 * Difference returns elements from the 'base' slice that are not present in the 'exclude' slice.
 * It compares each item in 'base' and includes it in the result only if it is not found in 'exclude'.
 *
 * Note: Duplicates from 'base' will be preserved if they are not in 'exclude'.
 *
 * Example:
 *
 *	base := []int{1, 2, 3, 4, 5}
 *	exclude := []int{2, 4}
 *
 *	result := Difference(base, exclude)
 *	// result: []int{1, 3, 5}
 */
func Difference[T comparable](base []T, exclude []T) []T {
	result := make([]T, 0)

	for _, item := range base {
		if !slices.Contains(exclude, item) {
			result = append(result, item)
		}
	}

	return result
}

/*
 * Intersect returns the intersection of two slices.
 * It compares the elements in 'source' and checks which ones also exist in 'target'.
 * The result includes all matching elements in the order they appear in 'source'.
 *
 * Note: Duplicates are not removed. If an item appears multiple times in 'source' and exists in 'target',
 * it will be included multiple times in the result.
 *
 * Example:
 *
 *	source := []int{1, 2, 2, 3, 4}
 *	target := []int{2, 4, 6}
 *
 *	result := Intersect(source, target)
 *	// result: []int{2, 2, 4}
 */
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
//	names := Map(users, func(u User, _ int) string {
//		return u.Name
//	})
//
//	names := Map(users, func(u User, index int) string {
//		return fmt.Sprintf("%d-%s", index, u.Name)
//	})
//
// This is a generic alternative to manual loops for transforming data.
func Map[T any, R any](items []T, selector func(item T, index int) R) []R {
	result := make([]R, 0, len(items))

	for i, item := range items {
		result = append(result, selector(item, i))
	}

	return result
}
