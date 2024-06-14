package utils

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestStrPinter(t *testing.T) {
	object := struct {
		value string
	}{
		value: "test",
	}

	result := StrPointer(object.value)
	if *result != object.value {
		t.Errorf("StrPointer(%v) = %v, want %v", object.value, *result, object.value)
	}

	object.value = "changed"
	if *result != "test" {
		t.Errorf("StrPointer(%v) = %v, want %v", object.value, *result, "test")
	}
}

func TestStrOmitEmpty(t *testing.T) {
	result := StrOmitEmpty("test")
	if *result != "test" {
		t.Errorf("StrOmitEmpty(%v) = %v, want %v", "test", result, "test")
	}

	result = StrOmitEmpty("")
	if result != nil {
		t.Errorf("StrOmitEmpty(%v) = %v, want %v", "", result, nil)
	}
}

func TestStrOmitWhitespace(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected *string
	}{
		{
			name:     "return original value if not empty",
			input:    "test",
			expected: StrOmitEmpty("test"),
		},
		{
			name:     "return nil if empty",
			input:    "",
			expected: nil,
		},
		{
			name:     "return nil if whitespace",
			input:    "  ",
			expected: nil,
		},
		{
			name:     "return nil if tab",
			input:    "\t",
			expected: nil,
		},
		{
			name:     "return nil if newline",
			input:    "\n",
			expected: nil,
		},
		{
			name:     "return nil if carriage return",
			input:    "\r",
			expected: nil,
		},
		{
			name:     "return nil if mixed whitespace",
			input:    " \t\n\r",
			expected: nil,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			result := StrOmitWhitespace(c.input)
			if diff := cmp.Diff(c.expected, result); diff != "" {
				t.Errorf("unexpected result (-want +got):\n%s", diff)
			}
		})
	}
}

func TestStrEmptyIfNil(t *testing.T) {
	input := "test"
	result := StrEmptyIfNil(&input)
	if result != "test" {
		t.Errorf("StrEmptyIfNil(%v) = %v, want %v", "test", result, "test")
	}

	result = StrEmptyIfNil(nil)
	if result != "" {
		t.Errorf("StrEmptyIfNil(%v) = %v, want %v", "nil", result, "")
	}
}

func TestStrCopyPointerValue(t *testing.T) {
	input := "test"
	result := StrCopyPointerValue(&input)
	if *result != "test" {
		t.Errorf("StrCopyPointerValue(%v) = %v, want %v", "test", result, "test")
	}

	result = StrCopyPointerValue(nil)
	if result != nil {
		t.Errorf("StrCopyPointerValue(%v) = %v, want %v", "nil", result, nil)
	}
}
