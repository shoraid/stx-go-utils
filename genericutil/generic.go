package genericutil

func Ptr[T any](v T) *T {
	return &v
}
