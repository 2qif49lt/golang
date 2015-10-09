package logx

import (
	"fmt"
	"testing"
)

func TestGetProcAbsDir(t *testing.T) {
	s, e := getProcAbsDir()
	fmt.Println(s, e)
}
func TestGetTimeStr(t *testing.T) {
	s := getTimeStr()
	fmt.Println(s)
}

func TestLog(t *testing.T) {
	l := NewLog("log", "base.log", Linfo)
	if l == nil {
		t.Error("l is nil")
	}
	l.SetFile(5, 1024*1024*5)
	for i := 0; i < 2000*1000; i++ {
		l.Log(i%Lmax, "%d", i)
	}
}

func TestAutoDel(t *testing.T) {
	l := NewLog("log", "multi.log", Ldebug)
	if l == nil {
		t.Error("l is nil")
	}

	l.SetFile(5, 1024*1024*5)

	for i := 0; i < 1000*1000; i++ {
		l.Log(i%Lmax, "%d", i)
	}
}

func TestConsole(t *testing.T) {
	SetLevel(Ldebug)
	for i := 0; i < 1000*1000; i++ {
		Log(i%Lmax, "%d", i)
	}
}

type countfile int

func (c *countfile) Dofile(fpath string) error {
	*c = *c + 1
	fmt.Println(*c)
	return nil
}
func TestHandler(t *testing.T) {
	l := NewLog("log", "base.log", Linfo)
	if l == nil {
		t.Error("l is nil")
	}
	l.SetFile(5, 1024*1024*5)
	h := countfile(0)
	l.SetHandler(&h)
	for i := 0; i < 2000*1000; i++ {
		l.Log(i%Lmax, "%d", i)
	}
}
