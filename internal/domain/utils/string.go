package utils

import "strings"

// TODO: ToPointer 関数に置き換える
func StrPointer(s string) *string {
	return &s
}

func StrOmitEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func StrOmitWhitespace(s string) *string {
	// すべて空白の場合はnilを返す
	if strings.TrimSpace(s) == "" {
		return nil
	}

	return &s
}

func StrEmptyIfNil(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func StrCopyPointerValue(s *string) *string {
	if s == nil {
		return nil
	}
	return StrPointer(*s)
}
