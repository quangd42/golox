package lox

import "errors"

var (
	ErrStackEmpty = errors.New("stack empty")
	ErrOutOfBound = errors.New("out of bound")
)

type stack struct {
	val []any
}

func newStack() *stack {
	return &stack{val: make([]any, 0)}
}

func (s *stack) push(v any) {
	s.val = append(s.val, v)
}

func (s *stack) pop() (any, error) {
	v, err := s.peek()
	if err != nil {
		return nil, err
	}
	s.val = s.val[:len(s.val)-1]
	return v, nil
}

func (s *stack) peek() (any, error) {
	if len(s.val) == 0 {
		return nil, ErrStackEmpty
	}
	return s.val[len(s.val)-1], nil
}

func (s *stack) isEmpty() bool {
	return len(s.val) == 0
}

func (s *stack) clear() {
	s.val = make([]any, 0)
}

func (s *stack) size() int {
	return len(s.val)
}

// get returns the item at the given index from the top of the stack.
func (s *stack) get(idx int) (any, error) {
	if idx < 0 || idx > s.size()-1 {
		return nil, ErrOutOfBound
	}
	return s.val[len(s.val)-1-idx], nil
}
