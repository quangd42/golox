package lox

type scopeStack struct {
	stack *stack
}

func newScopeStack() *scopeStack {
	return &scopeStack{stack: newStack()}
}

func (s *scopeStack) push(v map[string]bool) {
	s.stack.push(v)
}

func (s *scopeStack) pop() (map[string]bool, error) {
	val, err := s.stack.pop()
	if err != nil {
		return nil, err
	}
	return val.(map[string]bool), nil
}

func (s *scopeStack) peek() (map[string]bool, error) {
	val, err := s.stack.peek()
	if err != nil {
		return nil, err
	}
	return val.(map[string]bool), nil
}

func (s *scopeStack) isEmpty() bool {
	return s.stack.isEmpty()
}

func (s *scopeStack) clear() {
	s.stack.clear()
}

func (s *scopeStack) size() int {
	return s.stack.size()
}

func (s *scopeStack) get(idx int) (map[string]bool, error) {
	val, err := s.stack.get(idx)
	if err != nil {
		return nil, err
	}
	return val.(map[string]bool), nil
}
