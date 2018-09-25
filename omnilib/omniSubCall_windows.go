
// +build windows

package omnilib

// #include <stdio.h>
// #include <stdlib.h>
// #include "./omniproxy.h"
// #cgo CFLAGS: -I./
import "C"
import (
	//"unsafe"
	//"time"
	"time"
)

//var PtrLegacyRPCServer *Server=nil

func JsonCmdReqHcToOm(strReq string) string{
	strRsp:=C.GoString(C.CJsonCmdReq(C.CString(strReq)))
	return strRsp;
}
func LoadLibAndInit() {
	C.CLoadLibAndInit()
}

func OmniStart(strArgs string) {
	C.COmniStart(C.CString(strArgs))
}

func OmniCommunicate(netName string) {
//add by ycj 20180915
	LoadLibAndInit()
	OmniStart(netName)

	time.Sleep(time.Second*2)
	/*
	strReq := "{\"method\":\"omni_getinfo\",\"params\":[],\"id\":1}\n"
	strRsp := JsonCmdReqHcToOm(strReq)
	fmt.Println("in Go strRsp 1:", strRsp)
*/
	//legacyrpc.JsonCmdReqOmToHc((*C.char)(unsafe.Pointer(uintptr(0))));
}

/* abolish callback to LegacyRPCServer
//export JsonCmdReqOmToHc
func JsonCmdReqOmToHc(pcReq *C.char) *C.char {

	if PtrLegacyRPCServer==nil ||  pcReq==(*C.char)(unsafe.Pointer(uintptr(0))) {
		return (*C.char)(unsafe.Pointer(uintptr(0)))
	}
	strRsp,err:=PtrLegacyRPCServer.JsonCmdReq(C.GoString(pcReq))
	if err!=nil {
		return (*C.char)(unsafe.Pointer(uintptr(0)))
	}

	cs := C.CString(strRsp)

	defer func(){
		go func() {
			time.Sleep(time.Microsecond*1)
			C.free(unsafe.Pointer(cs))
		}()
	}()

	return cs
}
*/
