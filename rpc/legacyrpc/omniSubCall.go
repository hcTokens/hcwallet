package legacyrpc


// #include <stdio.h>
// #include <stdlib.h>
// #include "./omniproxy.h"
// #cgo CFLAGS: -I./
///* #cgo LDFLAGS:  -L./libomni -lomni*/
import "C"
import (
	"unsafe"
	"time"
)

var PtrLegacyRPCServer *Server=nil

func JsonCmdReqHcToOm(strReq string) string{
	strRsp:=C.GoString(C.CJsonCmdReq(C.CString(strReq)))
	return strRsp;
	//C.CJsonCmdReq(C.CString("abc"));
	//return main.CJsonCmdReqInGo("abc")
	//C.getchar();
	//return ""
}
func LoadLibAndInit(){
	C.CLoadLibAndInit()
}

func OmniStart(strArgs string){
	C.COmniStart(C.CString(strArgs))
}

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




