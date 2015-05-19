package xhbsrv

import (
	"unsafe"
)

// 拆包

// size类型固定为4字节
const sizebyte int = 4
const sizeheader int = 8

func (s *Srv) netunpack(buff []byte) ([]byte, error) {

	if len(buff) < sizebyte+s.PackSizePos {
		return nil, nil
	}
	msgsize := *(*uint32)(unsafe.Pointer(&buff[s.PackSizePos]))

	if msgsize > uint32(s.PackMaxLen) || msgsize < uint32(s.PackSizePos+sizebyte) {
		return nil, ErrSizeInvalid
	}
	if len(buff) < int(msgsize) {
		return nil, nil
	}
	retbyte := make([]byte, msgsize)
	copy(retbyte, buff[:msgsize])

	return retbyte, nil
}

func (s *Srv) netpack(cmd int, body []byte) ([]byte, error) {
	if sizeheader+len(body) > s.PackMaxLen {
		return nil, ErrSizeBodyTooLong
	}

	buff := make([]byte, sizeheader+len(body))
	phead := (*cmdheader)(unsafe.Pointer(&buff[0]))
	phead.cmd = 1
	phead.size = uint32(sizeheader + len(body))
	copy(buff[sizeheader:], body)
	return buff[:phead.size], nil
}
