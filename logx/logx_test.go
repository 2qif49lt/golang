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
	l := Newx("log", "base.log", Linfo)
	if l == nil {
		t.Error("l is nil")
	}
	for i := 0; i < 1000*1000; i++ {
		l.Log(i%Lmax, "%d", i)
	}
}

func TestAutoDel(t *testing.T) {
	l := Newx("log", "multi.log", Ldebug)
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
