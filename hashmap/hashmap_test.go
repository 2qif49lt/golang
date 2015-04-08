package hashmap

import (
	. "fmt"
	"math/rand"
	"testing"
)

type INT int

func (i INT) Hash() int {
	return int(i)
}

/*
func (i INT) Equal(k KeyInterFace) bool {
	j := k.(INT)
	return int(i) == int(j)
}
*/

func (i INT) Equal(j interface{}) bool {
	if r, ok := j.(INT); ok {
		return int(i) == int(r)
	} else {
		Println("fuck where")
		return false
	}

}

func Test_HashMap(t *testing.T) {
	m := New()

	for i := 0; i < 1000000; i++ {
		m.Put(INT(rand.Int()), i)
	}
	k, v := m.Info()
	Println(k, v)

	v = m.Get(k)
	Println(v)

	m.Erase(k)
	v = m.Get(k)
	Println(v)

	k, v = m.Info()
	Println(k, v)
}

func Benchmark_Cycle_Memory(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m := New()

		for i := 0; i < 1000000; i++ {
			m.Put(INT(rand.Int()), i)
		}

		m.Clear()
	}
}
