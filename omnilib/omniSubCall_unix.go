// +build linux aix darwin dragonfly freebsd  netbsd openbsd solaris

package omnilib

// #include <stdio.h>
// #include <stdlib.h>
// #include "./omniproxy.h"
// #cgo CFLAGS: -I./
//#cgo LDFLAGS:-L./ -lomnicored -lbitcoin_server -lbitcoin_common -lunivalue -lbitcoin_util -lbitcoin_wallet  -lbitcoin_consensus -lbitcoin_crypto -lleveldb -lmemenv -lsecp256k1 -lboost_system -lboost_filesystem -lboost_program_options -lboost_thread -lboost_chrono -ldb_cxx -lssl -lcrypto  -levent_pthreads -levent -lm -lstdc++
import "C"
import (
	//"unsafe"
	//"time"
	"time"
)



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

func OmniCommunicate() {
//add by ycj 20180915
	LoadLibAndInit()
	go OmniStart("exeName -regtest -txindex")

	time.Sleep(time.Second*2)

}

