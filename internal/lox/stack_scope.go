package lox

type scopeStack struct {
	stack *stack
}

func newScopeStack() *scopeStack {
	// HACK: perhaps generic is good here
	stack := newStack()
	stack.push(make(map[string]bool, 0))
	return &scopeStack{stack: stack}
}

func (s *scopeStack) push(v map[string]bool) {
	s.stack.push(v)
}

func (s *scopeStack) pop() (map[string]bool, error) {
	val, err := s.stack.pop()
	return val.(map[string]bool), err
}

func (s *scopeStack) peek() (map[string]bool, error) {
	val, err := s.stack.peek()
	return val.(map[string]bool), err
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
	return val.(map[string]bool), err
}
