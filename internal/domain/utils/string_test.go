package utils

import "testing"

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
