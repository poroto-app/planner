package array

func Find[T any](slice []T, condition func(T) bool) (T, bool) {
	var zero T
	for _, v := range slice {
		if condition(v) {
			return v, true
		}
	}
	return zero, false
}
