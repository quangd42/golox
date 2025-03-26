package lox

type loopStack struct {
	stack *stack
}

func newLoopStack() *loopStack {
	return &loopStack{stack: newStack()}
}

func (s *loopStack) push(v string) {
	s.stack.push(v)
}

func (s *loopStack) pop() (string, error) {
	val, err := s.stack.pop()
	return val.(string), err
}

func (s *loopStack) peek() (string, error) {
	val, err := s.stack.peek()
	return val.(string), err
}

func (s *loopStack) isEmpty() bool {
	return s.stack.isEmpty()
}

func (s *loopStack) clear() {
	s.stack.clear()
}

func (s *loopStack) size() int {
	return s.stack.size()
}

func (s *loopStack) get(idx int) (string, error) {
	val, err := s.stack.get(idx)
	return val.(string), err
}

func (s *loopStack) contains(val string) bool {
	for _, item := range s.stack.val {
		if str, ok := item.(string); ok && str == val {
			return true
		}
	}
	return false
}
