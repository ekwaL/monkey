package utils_test

import (
	"monkey/utils"
	"reflect"
	"testing"
)

func TestStack(t *testing.T) {
	s := utils.NewStack[int]()

	if !s.IsEmpty() {
		t.Error("Stack should be empty but it's not.")
	}

	if _, ok := s.Pop(); ok {
		t.Error("Expected error on Pop, but got none.")
	}

	if _, ok := s.Peek(); ok {
		t.Error("Expected error on Peek, but got none.")
	}

	if !s.IsEmpty() {
		t.Error("Stack should be empty but it's not.")
	}

	if l := s.List(); !reflect.DeepEqual(l, []int{}) {
		t.Errorf("Got wrong List, want %v, got %v.", []int{}, l)
	}

	push := []int{1, 2, 3, 4, 5}
	pop := []int{5, 4, 3, 2, 1}

	for _, p := range push {
		s.Push(p)
	}

	if l := s.List(); !reflect.DeepEqual(l, push) {
		t.Errorf("Got wrong List, want %v, got %v.", push, l)
	}

	if s.IsEmpty() {
		t.Error("Stack should not be empty but it is.")
	}

	got, ok := s.Peek()
	if !ok {
		t.Error("Unexpected error on Peek")
	}

	want := 5
	if got != want {
		t.Errorf("Got wrong element on Peek. Want %v got %v.", want, got)
	}

	for _, want := range pop {
		got, ok := s.Pop()
		if !ok {
			t.Error("Unexpected error on Pop")
		}
		if got != want {
			t.Errorf("Got wrong element on Pop. Want %v got %v.", want, got)
		}
	}

	if l := s.List(); !reflect.DeepEqual(l, []int{}) {
		t.Errorf("Got wrong List, want %v, got %v.", []int{}, l)
	}

	if !s.IsEmpty() {
		t.Error("Stack should be empty but it's not.")
	}

	if _, ok = s.Peek(); ok {
		t.Error("Expected error on Peek, but got none.")
	}

	i := 0
	for i = 0; i < 10000; i++ {
		s.Push(i)
	}

	for i := 10000 - 1; i >= 0; i-- {
		got, ok := s.Pop()
		if !ok {
			t.Fatal("Unexpected error on Pop.")
		}
		if got != i {
			t.Fatalf("Got wrong element on Pop. Want %v got %v.", i, got)
		}
	}

	if !s.IsEmpty() {
		t.Error("Stack should be empty but it's not.")
	}

	if _, ok = s.Peek(); ok {
		t.Error("Expected error on Peek, but got none.")
	}
}
