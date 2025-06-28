package genericutil

// FirstNonNil returns the value of the first pointer that is not nil.
//
// It accepts a variadic list of pointers to any type T and returns the dereferenced value
// of the first pointer that is not nil.
//
// If all pointers are nil, it returns the zero value of T.
//
// This is useful when you have multiple optional sources of data and want to pick
// the first one that is available (not nil).
//
// Example:
//
//	a := "hello"
//	b := "world"
//	FirstNonNil(nil, &a, &b) // returns "hello"
//
//	FirstNonNil[int](nil, nil) // returns 0
func FirstNonNil[T any](values ...*T) T {
	for _, v := range values {
		if v != nil {
			return *v
		}
	}

	var zero T

	return zero
}

// FirstNonZero returns the first value from the input list that is not the zero value of type T.
//
// It works with any comparable type (e.g. string, int, bool, struct with comparable fields, etc).
// If all values are equal to the zero value (e.g. 0 for int, "" for string), it returns the zero value.
//
// This is useful when you want to pick the first meaningful value from a list of fallbacks.
//
// Example:
//
//	FirstNonZero(0, 0, 5)        // returns 5
//	FirstNonZero("", "", "go")   // returns "go"
//	FirstNonZero(false, true)    // returns true
//
// If all values are zero, it returns the zero value of T.
func FirstNonZero[T comparable](values ...T) T {
	var zero T

	for _, v := range values {
		if v != zero {
			return v
		}
	}

	return zero
}

// Ptr is a generic helper function that returns a pointer to the given value.
// It is useful for creating pointers to literal values in a concise and readable way.
//
// Example:
//
//	Ptr("hello") returns *string pointing to "hello"
//	Ptr(42) returns *int pointing to 42
//
// Useful in tests, optional parameters, or working with pointer-based APIs.
func Ptr[T any](v T) *T {
	return &v
}
