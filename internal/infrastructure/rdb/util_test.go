package rdb

func arrayFromPointerSlice[T any](values []*T) []T {
	var array []T
	for _, value := range values {
		array = append(array, *value)
	}
	return array
}
