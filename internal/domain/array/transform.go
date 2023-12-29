package array

func Map[T any, U any](slice []T, transform func(T) U) []U {
	var mapped []U
	for _, v := range slice {
		mapped = append(mapped, transform(v))
	}
	return mapped
}

func MapAndFilter[T any, U any](slice []T, transform func(T) (U, bool)) []U {
	var mapped []U
	for _, v := range slice {
		if u, ok := transform(v); ok {
			mapped = append(mapped, u)
		}
	}
	return mapped
}
