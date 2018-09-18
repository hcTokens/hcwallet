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
	"fmt"
	"github.com/HcashOrg/hcwallet/omnilib"
)



//add by ycj 20180915
func omniCommunicate(){

	time.Sleep(time.Second*6)
	omnilib.LoadLibAndInit()


	//time.Sleep(time.Second*1000)
	go omnilib.OmniStart("exeName -regtest -txindex")

	time.Sleep(time.Second*9)
	strReq := "{\"method\":\"omni_getinfo\",\"params\":[],\"id\":1}\n"
	strRsp := omnilib.JsonCmdReqHcToOm(strReq)
	fmt.Println("in Go strRsp 1:", strRsp)

	//legacyrpc.JsonCmdReqOmToHc((*C.char)(unsafe.Pointer(uintptr(0))));
}





