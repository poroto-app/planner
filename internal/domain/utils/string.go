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

func StrCopyPointerValue(s *string) *string {
	if s == nil {
		return nil
	}
	return StrPointer(*s)
}
