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

func DistinctBy[T any, U comparable](slice []T, selector func(T) U) []T {
	var distinct []T
	seen := make(map[U]bool)
	for _, v := range slice {
		key := selector(v)
		if _, ok := seen[key]; !ok {
			seen[key] = true
			distinct = append(distinct, v)
		}
	}
	return distinct
}
