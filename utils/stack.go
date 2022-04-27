package utils

type Stack[T any] interface {
	Push(T)
	Pop() (T, bool)
	Peek() (T, bool)
	IsEmpty() bool
	List() []T
}

type stack[T any] struct {
	arr []T
	top int
}

func NewStack[T any]() *stack[T] {
	return &stack[T]{
		arr: make([]T, 0, 100),
		top: -1,
	}
}

func (s *stack[T]) Push(el T) {
	if s.top >= len(s.arr)-1 {
		s.arr = append(s.arr, el)
		s.top++
	} else {
		s.top++
		s.arr[s.top] = el
	}
}

func (s *stack[T]) Pop() (T, bool) {
	if s.IsEmpty() {
		return *new(T), false
	}

	el := s.arr[s.top]
	s.top--
	return el, true
}

func (s *stack[T]) Peek() (T, bool) {
	if s.IsEmpty() {
		return *new(T), false
	}

	el := s.arr[s.top]
	return el, true
}

func (s *stack[T]) List() []T {
	return s.arr[:s.top+1]
}

func (s *stack[T]) IsEmpty() bool {
	return s.top == -1
}
