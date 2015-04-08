package algo

import (
//	. "fmt"
)

func InsertSort(arr []Lesser) {
	for i := 0; i < len(arr); i++ {
		t := arr[i]
		j := 0
		for j = i - 1; j >= 0; j-- {
			if t.Less(arr[j]) {
				arr[j+1] = arr[j]
			} else {
				break
			}
		}
		arr[j+1] = t
	}
}

func SelectSort(arr []Lesser) {
	for i := 0; i < len(arr); i++ {
		min := arr[i]
		for j := i + 1; j < len(arr); j++ {
			if arr[j].Less(min) {
				min, arr[j] = arr[j], min
			}
		}
		arr[i] = min
	}
}

func partition(arr []Lesser) int {

	t := arr[0]

	i, j := 0, len(arr)-1
	for i < j {
		for ; j > i; j-- {
			if arr[j].Less(t) {
				arr[i] = arr[j]
				i++
				break
			}
		}
		for ; i < j; i++ {
			if t.Less(arr[i]) {
				arr[j] = arr[i]
				j--
				break
			}
		}
	}

	arr[i] = t
	return i
}
func QuickSort(arr []Lesser) {
	if len(arr) > 1 {
		idx := partition(arr)
		QuickSort(arr[:idx])
		QuickSort(arr[idx+1:])
	}
}
