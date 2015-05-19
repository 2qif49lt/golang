package xhbsrv

import (
	. "fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"sync/atomic"
	"testing"
	"time"
)

func TestSrv(t *testing.T) {
	s := New()
	err := s.RunOn(":9888")
	if err != nil {
		Println(err)
		return
	}
	for {
		Println(s.ConnNum())
		time.Sleep(time.Second)
	}
}

type testcb struct {
	*Srv
}

var sendcounter int64 = 0

func (cb *testcb) OnNewConn(addr net.Addr, connid int) bool {
	Println("OnNewConn", addr, connid)
	atomic.AddInt64(&sendcounter, 1)
	lj := Sprintf("hello are you ok! %d", atomic.LoadInt64(&sendcounter))
	body := make([]byte, len(lj)+1)
	copy(body, lj)
	if err := cb.SendCmd(connid, 1, body); err != nil {
		Println("OnNewConn SendCmd", err)
	}
	return true
}

var recvcounter int = 0

func (cb *testcb) OnRecvMsg(addr net.Addr, connid int, msg []byte) bool {
	header, body := msgunpack(msg)
	Println("OnRecvMsg", addr, connid, "msglen", len(msg), header, len(body), string(body))

	atomic.AddInt64(&sendcounter, 1)
	lj := Sprintf("hello are you ok! %d", atomic.LoadInt64(&sendcounter))
	sbody := make([]byte, len(lj)+1)
	copy(sbody, lj)
	if err := cb.SendCmd(connid, 1, sbody); err != nil {
		Println("OnNewConn SendCmd", err)
	}
	return true
}
func (cb *testcb) OnConnClosed(addr net.Addr, connid int, err error) {
	Println("OnConnClosed", addr, connid, err)
}
func (cb *testcb) OnConnectRemote(addr net.Addr, connid int, err error) bool {
	Println("OnConnectRemote", addr, connid, err)
	return true
}

var waitch = make(chan int, 1)

func TestConnect(t *testing.T) {
	s := New()
	s.RegCb(&testcb{})

	connid, err := s.Connect("10.1.9.34:9889", time.Second*10)
	Println(connid, err)
	if err != nil {
		return
	}

	<-waitch
}

func TestConnectAndSendRecv(t *testing.T) {
	s := New()
	s.RegCb(&testcb{s})

	_, err := s.Connect("10.1.9.34:9889", time.Second*10)
	if err != nil {
		return
	}
	<-waitch
	time.Sleep(time.Second)
}
func TestConnectAndSendRecvTwiceForUnpack(t *testing.T) {
	s := New()
	s.RegCb(&testcb{s})

	_, err := s.Connect("10.1.9.34:9889", time.Second*10)
	if err != nil {
		return
	}
	<-waitch
	time.Sleep(time.Second)
}

func TestActiveSend(t *testing.T) {
	s := New()
	s.RegCb(&testcb{s})

	connid, err := s.ConnectBlock("10.1.9.34:9889", time.Second*10)
	if err != nil {
		return
	}

	err = s.SendCmd(connid, 1, []byte("hello c++!"))
	if err != nil {
		return
	}

	<-waitch
	time.Sleep(time.Second)
}

func TestActiveClose(t *testing.T) {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	s := New()
	s.RegCb(&testcb{s})

	p := make([]runtime.StackRecord, 100)
	fn, _ := runtime.GoroutineProfile(p)
	Println(p[:fn])
	Println("runtime.NumGoroutine()", runtime.NumGoroutine())

	connid, err := s.ConnectBlock("10.1.9.34:9889", time.Second*10)
	if err != nil {
		return
	}
	time.Sleep(time.Second * 50)
	fn, _ = runtime.GoroutineProfile(p)
	Println(p[:fn])
	Println("runtime.NumGoroutine()", runtime.NumGoroutine())

	s.CloseConn(connid)

	time.Sleep(time.Second * 50)
	fn, _ = runtime.GoroutineProfile(p)
	Println(p[:fn])
	Println("runtime.NumGoroutine()", runtime.NumGoroutine())
}

func TestSrvClientSR(t *testing.T) {
	s := New()
	err := s.RunOn(":9888")
	if err != nil {
		Println(err)
		return
	}
	s.RegCb(&testcb{s})
	for {
		Println(s.ConnNum(), atomic.LoadInt64(&sendcounter))
		time.Sleep(time.Second)
	}
}
