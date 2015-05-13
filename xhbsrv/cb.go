package xhbsrv

import (
	//"errors"
	"net"
	"unsafe"
)

type CallBack interface {
	OnNewConn(addr net.Addr, connid int) bool
	OnRecvMsg(addr net.Addr, connid int, msg []byte) bool
	OnConnClosed(addr net.Addr, connid int, err error)
	OnConnectRemote(addr net.Addr, connid int, err error) bool
}

func (s *Srv) RegCb(cb CallBack) {
	s.cb = cb
}

type defaultcb struct {
	recv int
	send int
	conn int
}

func (cb *defaultcb) OnNewConn(addr net.Addr, connid int) bool {
	return true
}
func (cb *defaultcb) OnRecvMsg(addr net.Addr, connid int, msg []byte) bool {
	return true
}
func (cb *defaultcb) OnConnClosed(addr net.Addr, connid int, err error) {

}
func (cb *defaultcb) OnConnectRemote(addr net.Addr, connid int, err error) bool {
	return true
}

type cmdheader struct {
	size uint32
	cmd  uint32
}

func msgunpack(buff []byte) (*cmdheader, []byte) {
	header := (*cmdheader)(unsafe.Pointer(&buff[0]))
	body := buff[unsafe.Sizeof(header):]
	return header, body
}
