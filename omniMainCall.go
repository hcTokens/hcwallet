package main

/*
#cgo CFLAGS: -I./
#include <stdio.h>
#include <stdlib.h>
*/
import "C"
import (
//	"unsafe"
	"time"
	"github.com/HcashOrg/hcwallet/rpc/legacyrpc"
	"fmt"
)

/*
//export JsonCmdReqOmToHc
func JsonCmdReqOmToHc(pcReq *C.char) *C.char {

	if PtrLegacyRPCServer==nil ||  pcReq==(*C.char)(unsafe.Pointer(uintptr(0))) {
		return (*C.char)(unsafe.Pointer(uintptr(0)))
	}
	strRsp,err:=PtrLegacyRPCServer.JsonCmdReq(C.GoString(pcReq))
	if err!=nil {
		return (*C.char)(unsafe.Pointer(uintptr(0)))
	}

	//strRsp:="test"
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


/*
func CJsonCmdReqInGo(strReq string) string{
	strRsp:=C.GoString(C.CJsonCmdReq(C.CString(strReq)))
	return strRsp;
}
*/

//add by ycj 20180915
func omniCommunicate(){
	//omniCommunicate();
	time.Sleep(time.Second*6)
	legacyrpc.LoadLibAndInit()

	//fn:=printlnTest//GJsonCmdReq
	//omni.CSetCallback(1,unsafe.Pointer(&fn))

	//time.Sleep(time.Second*1000)
	go legacyrpc.OmniStart("exeName -regtest -txindex")

	time.Sleep(time.Second*9)
	strReq := "{\"method\":\"omni_getinfo\",\"params\":[],\"id\":1}\n"
	strRsp := legacyrpc.JsonCmdReqHcToOm(strReq)
	fmt.Println("in Go strRsp 1:", strRsp)

	//legacyrpc.JsonCmdReqOmToHc((*C.char)(unsafe.Pointer(uintptr(0))));
}





