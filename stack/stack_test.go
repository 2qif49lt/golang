/**
 * Created with IntelliJ IDEA.
 * User: xuhaibin
 * Date: 13-12-25
 * Time: 下午3:20
 * To change this template use File | Settings | File Templates.
 */
package stack

import "testing"

func Teststack(t *testing.T){
	s := New(10)
	for i:= 0; i != 100; i++{
		s.Push(i)
	}
	if s.Len() != 100 {
		t.Errorf("Len() Error.s.len: %d != 100",s.Len())
	}
	if s.Pop().(int) != 99{
		t.Errorf("Pop Error")
	}
}
