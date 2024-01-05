package array

func Filter[T any](slice []T, condition func(T) bool) []T {
	var filtered []T
	for _, v := range slice {
		if condition(v) {
			filtered = append(filtered, v)
		}
	}
	return filtered
}
