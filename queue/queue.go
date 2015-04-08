package queue

import (
	. "fmt"
	"sync"
)

const (
	CHUNK_SIZE = 64
)

type chunk struct {
	items       [CHUNK_SIZE]interface{}
	first, last int
	next        *chunk
}

func (q chunk) String() string {
	return Sprintf("first: %p last: %p next: %p", q.first, q.last, q.next)
}

type Queue struct {
	head, tail *chunk
	size       int
	lock       sync.Mutex
}

func (q Queue) String() string {
	return Sprintf("head: %v  tail: %v size: %d", q.head, q.tail, q.size)
}
func New() (q *Queue) {
	item := new(chunk)
	return &Queue{
		head: item,
		tail: item,
	}
}
func (q *Queue) Len() (length int) {
	q.lock.Lock()
	defer q.lock.Unlock()

	length = q.size
	return
}
func (q *Queue) Size() (length int) {
	q.lock.Lock()
	defer q.lock.Unlock()

	length = q.size
	return
}
func (q *Queue) Push(item interface{}) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if item == nil {
		return
	}

	if q.tail.last >= CHUNK_SIZE {
		q.tail.next = new(chunk)
		q.tail = q.tail.next
	}

	q.tail.items[q.tail.last] = item
	q.tail.last++
	q.size++
}

func (q *Queue) Pop() (item interface{}) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.size == 0 {
		return nil
	}
	if q.head.first >= q.head.last {
		return nil
	}
	item = q.head.items[q.head.first]
	q.head.first++
	q.size--

	if q.head.first >= q.head.last {
		if q.size == 0 {
			q.head.first = 0
			q.head.last = 0
			q.head.next = nil
		} else {
			q.head = q.head.next
		}
	}
	return
}
func (q *Queue) Info() {
	Println(q.size, q.head, q.tail)
}
