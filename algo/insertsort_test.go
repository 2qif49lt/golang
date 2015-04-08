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

func (i INT) Less(r interface{}) bool {
	if j, ok := r.(INT); ok {
		return int(i) < int(j)
	} else {
		Println("fuck")
		return false
	}
}

func TestInsertSort(t *testing.T) {
	arr := []INT{4, 1, 6, 3, 1, 5, 7, 9, 3, 2}

	morearr := make([]Lesser, len(arr))
	for i := range arr {
		morearr[i] = arr[i]
	}
	InsertSort(morearr)
	Println(morearr)
}
func TestSelectSort(t *testing.T) {
	arr := []INT{4, 1, 6, 3, 1, 5, 7, 9, 3, 2}

	morearr := make([]Lesser, len(arr))
	for i := range arr {
		morearr[i] = arr[i]
	}
	SelectSort(morearr)

	Println(morearr)
}

func TestQuickSort(t *testing.T) {
	//arr := []INT{4, 1, 6, 3, 1, 5, 7, 9, 3, 2}
	arr := []INT{9, 8, 7, 6, 5, 4, 3, 2, 1, 0}
	morearr := make([]Lesser, len(arr))
	for i := range arr {
		morearr[i] = arr[i]
	}

	QuickSort(morearr)
	Println(morearr)
}
