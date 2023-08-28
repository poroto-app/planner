package utils

func StrPointer(s string) *string {
	return &s
}

func StrOmitEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
