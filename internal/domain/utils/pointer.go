package utils

func ToPointer[T any](x T) *T {
	return &x
}

func FromPointerOrZero[T any](x *T) T {
	var zero T
	if x == nil {
		return zero
	}
	return *x
}
