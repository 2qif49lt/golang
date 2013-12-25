// Package cstruct tramsforms golang struct into c style struct for network transfer to c/c++ programs.
package cstruct


type msg_net_inner struct{
	uinum int64
}
type msg_net_example struct{
	uiid uint32
	iall [10]int32

	szname string `32`	// as char szname[32] which can store '32' bytes most.tag must be a positive integer.
  	fixarr [32]byte		// it is analogous to szname

	inmsg msg_net_inner		// inner struct
  	inall [3]msg_net_inner	// struct array

	uisize uint32	// the size of chbuff
  	chbuff []byte `uisize`	// uncertain byte array.usually in the end of msg,buf not necessary.
}


func Trans(i interface{})([]byte, error){
	s := make([]byte,0)

	switch i.(type){
		case bool,int,uint,uintptr,i

	}

	return s
}

func invalidtype(i)
