package array

func Map[T any, U any](slice []T, transform func(T) U) []U {
	var mapped []U
	for _, v := range slice {
		mapped = append(mapped, transform(v))
	}
	return mapped
}

func MapWithErr[T any, U any](slice []T, transform func(T) (*U, error)) (*[]U, error) {
	mapped := make([]U, 0, len(slice))
	for _, v := range slice {
		u, err := transform(v)
		if err != nil {
			return nil, err
		}

		if u == nil {
			continue
		}

		mapped = append(mapped, *u)
	}
	return &mapped, nil
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

func FlatMap[T any, U any](slice []T, transform func(T) []U) []U {
	var mapped []U
	for _, v := range slice {
		mapped = append(mapped, transform(v)...)
	}
	return mapped
}

func Flatten[T any](slice [][]T) []T {
	var flattened []T
	for _, v := range slice {
		flattened = append(flattened, v...)
	}
	return flattened
}
