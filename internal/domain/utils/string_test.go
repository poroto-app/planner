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
