package lox

import "errors"

var (
	ErrStackEmpty = errors.New("stack empty")
	ErrOutOfBound = errors.New("out of bound")
)

type scopeStack struct {
	val []map[string]bool
}

func newScopeStack() *scopeStack {
	return &scopeStack{val: make([]map[string]bool, 0)}
}

func (s *scopeStack) push(v map[string]bool) {
	s.val = append(s.val, v)
}

func (s *scopeStack) pop() (map[string]bool, error) {
	v, err := s.peek()
	if err != nil {
		return nil, err
	}
	s.val = s.val[:len(s.val)-1]
	return v, nil
}

func (s *scopeStack) peek() (map[string]bool, error) {
	if len(s.val) == 0 {
		return nil, ErrStackEmpty
	}
	return s.val[len(s.val)-1], nil
}

func (s *scopeStack) isEmpty() bool {
	return len(s.val) == 0
}

func (s *scopeStack) clear() {
	s.val = nil
}

func (s *scopeStack) size() int {
	return len(s.val)
}

// get returns the scope at the given index from the top of the stack.
func (s *scopeStack) get(idx int) (map[string]bool, error) {
	if idx < 0 || idx > s.size()-1 {
		return nil, ErrOutOfBound
	}
	return s.val[len(s.val)-1-idx], nil
}
