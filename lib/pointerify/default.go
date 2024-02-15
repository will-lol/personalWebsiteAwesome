package pointerify

func Pointer[T any](s T) *T {
	return &s
}

func DePointer[T any](s *T) T {
	return *s
}
