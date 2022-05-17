package object_test

import (
	"monkey/object"
	"testing"
)

func TestStringHashKey(t *testing.T) {
	h1 := &object.String{Value: "Hello"}
	h2 := &object.String{Value: "Hello"}
	diff := &object.String{Value: "World"}

	if h1.HashKey() != h2.HashKey() {
		t.Errorf("Strings %v and %v should have same HashKey.", h1, h2)
	}

	if h1.HashKey() == diff.HashKey() {
		t.Errorf("Strings %v and %v should have different HashKey.", h1, diff)
	}
}
