package array

func IsContain(array []string, value string) bool {
	for _, elem := range array {
		if value == elem {
			return true
		}
	}
	return false
}

// Check whether including or not
func HasIntersection(a, b []string) bool {
	for _, a_value := range a {
		if IsContain(b, a_value) {
			return true
		}
	}
	return false
}
