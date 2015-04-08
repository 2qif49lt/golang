package hashmap

import (
	. "fmt"
)

type KeyInterFace interface {
	Hash() int
	Equal(interface{}) bool
	//Equal(KeyInterFace) bool
}

type pair struct {
	key  KeyInterFace
	val  interface{}
	next *pair
}

type bucket struct {
	num  int
	link *pair
}

type HashMap struct {
	bkts []bucket
	num  int
}

func New() *HashMap {
	b := make([]bucket, 32, 32)
	return &HashMap{
		bkts: b,
	}
}

const growrate = 1.5

func (m *HashMap) grow() {
	nl := len(m.bkts) * 2
	na := make([]bucket, nl, nl)

	for _, bkt := range m.bkts {
		p := bkt.link
		for p != nil {
			hash := p.key.Hash()
			i := hash % nl

			tmp := p.next

			p.next = na[i].link
			na[i].link = p
			na[i].num++

			p = tmp
		}
	}

	m.bkts = na
}
func (m *HashMap) Put(k KeyInterFace, v interface{}) {
	if fnum := float64(m.num); fnum/float64(len(m.bkts)) > growrate {
		//	Printf("grow() fnum: %f len(m.bkts): %d new len: %d\n", fnum, len(m.bkts), 2*len(m.bkts))
		m.grow()
	}

	hash := k.Hash()
	i := hash % len(m.bkts)

	tmp := m.bkts[i].link
	for tmp != nil {
		if tmp.key.Equal(k) {
			return
		}
		tmp = tmp.next
	}

	p := &pair{
		key:  k,
		val:  v,
		next: m.bkts[i].link,
	}
	m.bkts[i].link = p
	m.bkts[i].num++

	m.num++
}
func (m *HashMap) Get(k KeyInterFace) interface{} {
	hash := k.Hash()
	i := hash % len(m.bkts)

	tmp := m.bkts[i].link
	for tmp != nil {
		if tmp.key.Equal(k) {
			return tmp.val
		}
		tmp = tmp.next
	}
	return nil
}

func (m *HashMap) Size() int {
	return m.num
}

func (m *HashMap) Erase(k KeyInterFace) {
	hash := k.Hash()
	i := hash % len(m.bkts)

	b := false
	beg := m.bkts[i].link

	if beg.key.Equal(k) {
		m.bkts[i].link = m.bkts[i].link.next
		b = true
	} else {
		prev := m.bkts[i].link
		cur := m.bkts[i].link.next

		for cur != nil {
			if cur.key.Equal(k) {
				prev.next = cur.next
				b = true
				break
			}
			cur = cur.next
		}
	}
	if b {
		m.bkts[i].num--
		m.num--
	}
	return
}
func (m *HashMap) Clear() {
	m.bkts = make([]bucket, 32, 32)
	m.num = 0
}
func (m *HashMap) Info() (k KeyInterFace, v interface{}) {
	Printf("size: %d ", m.num)
	count := 0
	bignum := 0
	bigidx := 0
	for i, bkt := range m.bkts {
		count += bkt.num
		if bkt.num > bignum {
			bigidx = i
			bignum = bkt.num
		}
	}

	Printf("count: %d bigidx: %d bignum: %d ", count, bigidx, bignum)
	if m.bkts[bigidx].link != nil {
		tmp := m.bkts[bigidx].link

		i := 0
		for tmp != nil {
			Printf("pair: %v ", tmp)
			tmp = tmp.next

			if i++; i == 3 {
				k = tmp.key
				v = tmp.val
			}
		}
	}
	Println()
	return
}
