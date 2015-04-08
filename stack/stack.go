package stack

type node struct {
	value interface{}
	next  *node
}
type Stack struct {
	top  *node
	size int
}

func New() *Stack {
	return &Stack{
		top:  nil,
		size: 0,
	}
}

func (s *Stack) Size() int {
	return s.size
}

func (s *Stack) Push(v interface{}) {
	n := &node{
		value: v,
		next:  s.top,
	}
	s.top = n
	s.size++
}

func (s *Stack) Pop() interface{} {
	if s.size == 0 {
		return nil
	}
	n := s.top
	s.size--
	s.top = s.top.next
	return n.value
}

func (s *Stack) Peek() interface{} {
	if s.size == 0 {
		return nil
	}

	return s.top.value
}
