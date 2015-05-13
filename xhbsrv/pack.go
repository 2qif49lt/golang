package xhbsrv

import (
	"unsafe"
)

// 拆包

// size类型固定为4字节

func (s *Srv) netunpack(buff []byte) ([]byte, error) {
	const sizebyte int = 4
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
