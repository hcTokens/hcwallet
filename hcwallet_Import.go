// Copyright (c) 2018-2020 The Hc developers

package main

/*
#cgo CFLAGS: -I./
#include <ImportDll.h>
#include <stdlib.h>
#include <stdio.h>
*/
import "C"
import (
	"fmt"
	"github.com/HcashOrg/hcd/hcjson"
	"github.com/HcashOrg/hcwallet/rpc/legacyrpc"
	"github.com/HcashOrg/hcwallet/wallet"
	"strings"
	"sync"
	"time"
	"unsafe"
)

const CBINDEX_PROCESS_PAYLOAD = 4

const CBINDEX_GETHEIGHT = 10
const CBINDEX_GETHASH = 11
const CBINDEX_CREATETX = 12
const CBINDEX_VALIDATEADDR = 13

var LegacyServer *legacyrpc.Server

//export GoCallback
func GoCallback(nType C.int) C.int {
	switch nType {
	case CBINDEX_GETHEIGHT:
		{
			var height int64
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				/*
					best := serverCall.blockManager.chain.BestSnapshot()
					height = best.Height
				*/
				wg.Done()
			}()
			wg.Wait()
			return C.int(height)
		}
	}
	return nType
}

//export GoCallbackChar
func GoCallbackChar(nType C.int, content *C.char) *C.char {
	var retCallback *C.char = nil
	defer func() {
		go func() {
			time.Sleep(time.Millisecond * 50)
			C.free(unsafe.Pointer(retCallback))
		}()
	}()

	switch nType {
	case CBINDEX_GETHASH:
		{
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				/*
					best := serverCall.blockManager.chain.BestSnapshot()
						retCallback = (*C.char)(unsafe.Pointer(C.CString(best.Hash.String())))
				*/
				wg.Done()
			}()
			wg.Wait()
			return retCallback
		}
	case CBINDEX_CREATETX: //CreateRawTransaction
		{
			var output []byte
			p := uintptr(unsafe.Pointer(content))
			value := *(*byte)(unsafe.Pointer(p))
			for value != 0 {
				output = append(output, value)
				p += unsafe.Sizeof(value)
				value = *(*byte)(unsafe.Pointer(p))
			}
		}
	}
	//	}
	return nil
}

//export GoCallbackCharEx
func GoCallbackCharEx(nType C.int, content *C.char, length C.int) *C.char {
	var retCallback *C.char = nil
	defer func() {
		go func() {
			time.Sleep(time.Millisecond * 50)
			C.free(unsafe.Pointer(retCallback))
		}()
	}()

	var output []byte
	p := uintptr(unsafe.Pointer(content))
	value := *(*byte)(unsafe.Pointer(p))
	for i := int(0); i < int(length); i++ {
		output = append(output, value)
		p += unsafe.Sizeof(value)
		value = *(*byte)(unsafe.Pointer(p))
		fmt.Println(i)
	}

	loader := LegacyServer.GetWallet()
	if loader == nil {
		return retCallback
	}
	w, _ := loader.LoadedWallet()
	if w == nil {
		return retCallback
	}

	var ret string
	var err error
	switch nType {
	case CBINDEX_CREATETX:
		{
			ret, err = cbCreateRawTransaction(string(output), w)
			if err != nil {
				fmt.Println(err)
			}
		}
	case CBINDEX_VALIDATEADDR:
		{
			ret, err = legacyrpc.DllCallValidateAddress(string(output), w)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	retCallback = (*C.char)(unsafe.Pointer(C.CString(ret)))
	return retCallback
}

func cbCreateRawTransaction(text string, w *wallet.Wallet) (string, error) {
	con := text
	fmt.Println(con)

	conSet := strings.Split(con, ";;;")
	fmt.Println(conSet)
	var fromAddr, payLoad []byte
	for i := 0; i < len(conSet); i++ {
		item := strings.Split(conSet[i], "===")
		if len(item) == 2 {
			if item[0] == "fromaddress" {
				fromAddr = []byte(item[1])
			} else if item[0] == "payload" {
				payLoad = []byte(item[1])
			}
		}
	}
	fmt.Println(fromAddr)
	fmt.Println(payLoad)

	cmd := &hcjson.SendToAddressCmd{
		Address: string(fromAddr),
		Amount:  100,
	}
	fmt.Println(cmd)

	fmt.Println("11111111111111111111111")
	return legacyrpc.DllCallsendToAddress(cmd, w, payLoad)
}

func GoCallCpp(nType int32, str string) {
	retCallback := (*C.char)(unsafe.Pointer(C.CString(str)))
	C.CallCpp(C.int(nType), unsafe.Pointer(retCallback))
	C.free(unsafe.Pointer(retCallback))
}

func GoCallCppInit(recieve <-chan string) {
	for {
		select {
		case value, ok := <-recieve:
			if !ok {
				return
			}
			fmt.Println(value)
			GoCallCpp(CBINDEX_PROCESS_PAYLOAD, value)
		}
	}
}
