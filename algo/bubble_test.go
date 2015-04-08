package algo

import (
	. "fmt"
	"testing"
)

type INT int

func (i INT) More(r interface{}) bool {
	if j, ok := r.(INT); ok {
		return int(i) > int(j)
	} else {
		Println("fuck")
		return false
	}
}
func TestBubbleSort(t *testing.T) {

	arr := []INT{4, 1, 6, 3, 1, 5, 7, 9, 3, 2}

	morearr := make([]Morer, len(arr))
	for i := range arr {
		morearr[i] = arr[i]
	}
	BubbleSort(morearr)
	Println(morearr)
}
