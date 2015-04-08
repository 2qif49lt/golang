package list

import (
	. "fmt"
)

type node struct {
	Value      interface{}
	prev, next *node
}

type Iter *node

func (n node) Next() Iter {
	return n.next
}
func (n node) Prev() Iter {
	return n.prev
}

type List struct {
	head *node
	tail *node
	size int
}

func (l *List) Info() {
	Println(l.head)
	Println(l.tail)
	Println(l.size)
}

func New() (l *List) {
	l = &List{
		head: nil,
		tail: nil,
		size: 0,
	}
	return
}

func (l *List) Size() (size int) {
	size = l.size
	return
}

func (l *List) Erase(iter Iter) Iter {
	if iter == nil {
		return nil
	}
	cur := (*node)(iter)
	prev := cur.prev

	if prev != nil {
		prev.next = cur.next
	} else {
		l.head = cur.next
	}

	if cur.next != nil {
		cur.next.prev = prev
	} else {
		l.tail = nil
	}

	l.size--
	return cur.next
}
func (l *List) Front() Iter {
	return l.head
}
func (l *List) Back() Iter {
	return l.tail
}
func (l *List) Find(v interface{}) Iter {
	for n := l.head; n != nil; n = n.next {
		if n.Value == v {
			return n
		}
	}
	return nil
}
func (l *List) PushBack(v interface{}) Iter {
	iter := &node{
		Value: v,
		prev:  l.tail,
		next:  nil,
	}
	if l.head == nil {
		l.head = iter
	}
	if l.tail != nil {
		l.tail.next = iter
	}
	l.tail = iter
	l.size++
	return iter
}
func (l *List) PushFront(v interface{}) Iter {
	iter := &node{
		Value: v,
		prev:  nil,
		next:  l.head,
	}
	if l.tail == nil {
		l.tail = iter
	}
	if l.head != nil {
		l.head.prev = iter
	}
	l.head = iter
	l.size++
	return iter
}

func (l *List) InsertAfter(v interface{}, pos Iter) Iter {
	if pos == nil {
		return nil
	}

	iter := &node{
		Value: v,
		prev:  pos,
		next:  pos.next,
	}

	if pos.next == nil {
		l.tail = iter
	} else {
		pos.next.prev = iter
	}

	pos.next = iter

	l.size++
	return iter
}
func (l *List) InsertBefore(v interface{}, pos Iter) Iter {
	if pos == nil {
		return nil
	}
	iter := &node{
		Value: v,
		prev:  pos.prev,
		next:  pos,
	}
	if pos.prev == nil {
		l.head = iter
	} else {
		pos.prev.next = iter
	}

	pos.prev = iter

	l.size++
	return iter
}
