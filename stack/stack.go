// Package stack implements a LIFO container base on slice.
package stack

// Stack represents a slice stack.
type Stack struct{
	s []interface{}
}


// create a stack start with n capacity.
func New(n int)(r *Stack){
	r = Stack{make([]interface{},n)}
	return
}

// Len returns the number of elements of stack s.
func (s Stack) Len() int {
	return len(s.s)
}

// get the first refer value or nil.
func (s *Stack) Top()(r interface{}){
	if(s.Len() > 0){
		 r = s[s.Len() - 1]
	} else{
		r = nil
	}
	return
}

// Pop returns the top element or nil,if not nil,will delete the element.
func (s *Stack) Pop()(r interface{}){
	if(s.Len() > 0){
		return nil
	}else{
		r = s.s[s.Len() - 1]
		s.s = s.s[:s.Len() - 1]
	}
	return r
}

// Push a new element at the top of stack s.
func (s *Stack) Push(i interface{}){
	append(s.s,i)
}



