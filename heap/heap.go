package heap

import (
	. "fmt"
)

type Interface interface {
	Less(r interface{}) bool
}

type Heap struct {
	arr []Interface
}

// 最小堆
/*
	位置:   	1  2  3  4  5  6  7  8  9 10 11
	最小堆:  	1  2  3  4  5  6  7  8  9 10 11
	最大堆: 11  9 10  5  6  7  8  1  2  3  4

	父节点： (i - 1)/2
	子节点： 2i + 1,2i + 2
*/

func New() *Heap {
	Println(4)
	a := make([]Interface, 0)
	return &Heap{
		arr: a,
	}
}

func (h *Heap) Size() int {
	return len(h.arr)
}

func (h *Heap) Push(v Interface) {
	h.arr = append(h.arr, v)

	i := h.Size() - 1
	j := (i - 1) / 2

	for i != j && j >= 0 && i > 0 && h.arr[i].Less(h.arr[j]) {
		h.arr[i], h.arr[j] = h.arr[j], h.arr[i]
		i = j
		j = (i - 1) / 2

	}
}

func (h *Heap) Pop() Interface {
	if h.Size() == 0 {
		return nil
	}

	v := h.arr[0]
	h.arr[0] = h.arr[h.Size()-1]

	var pos int
	if h.Size() > 1 {
		pos = h.Size() - 1
	} else {
		pos = 0
	}

	h.arr = h.arr[:pos]

	i := 0
	j := 2*i + 1
	for j <= h.Size()-1 {
		if j+1 <= h.Size()-1 && h.arr[j].Less(h.arr[j+1]) == false {
			j++
		}
		if h.arr[i].Less(h.arr[j]) == false {
			h.arr[i], h.arr[j] = h.arr[j], h.arr[i]
		} else {
			break
		}
		i = j
		j = 2*i + 1
	}
	return v
}
