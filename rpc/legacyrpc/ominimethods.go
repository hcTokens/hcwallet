// Copyright (c) 2013-2016 The btcsuite developers
// Copyright (c) 2015-2017 The Decred developers
// Copyright (c) 2018-2020 The Hc developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package legacyrpc

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/HcashOrg/hcd/hcjson"
	"github.com/HcashOrg/hcd/hcutil"
	"github.com/HcashOrg/hcwallet/apperrors"
	"github.com/HcashOrg/hcwallet/omnilib"
	"github.com/HcashOrg/hcwallet/wallet"
	"github.com/HcashOrg/hcwallet/wallet/txrules"
	"github.com/HcashOrg/hcwallet/wallet/udb"
)

//add by ycj 20180915
//commonly used cmd request
func OmnCmdReq(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	byteCmd, err := hcjson.MarshalCmd(1, icmd)
	if err != nil {
		return err, nil
	}
	strReq := string(byteCmd)
	strRsp := omnilib.JsonCmdReqHcToOm(strReq)

	payLoad, err := hex.DecodeString(strRsp)
	if err == nil {
		/*
			{"jsonrpc":"1.0","method":"omni_sendissuancefixed","params":["Tsk6gAJ7X9wjihFPo4nt5HHa9GNZysTyugn",2,1,0,"Companies","Bitcoin Mining","Quantum Miner","","","1000000"],"id":1}

		*/
		var req hcjson.Request
		err = json.Unmarshal(byteCmd, &req)
		addr := req.Params[0]
		cmd := &SendFromAddressToAddressCmd{
			FromAddress:   string(addr[1 : len(addr)-1]),
			ToAddress:     string(addr[1 : len(addr)-1]),
			ChangeAddress: string(addr[1 : len(addr)-1]),
			Amount:        10,
		}
		fmt.Println(cmd)
		return omniSendToAddress(cmd, w, payLoad)
	}
	var response hcjson.Response
	_ = json.Unmarshal([]byte(strRsp), &response)
	//strResult:=string(response.Result);
	return response.Result, nil
}

func getOminiMethod() map[string]LegacyRpcHandler {
	return map[string]LegacyRpcHandler{

		"omni_getinfo":                     {handler: omni_getinfo}, //by ycj 20180915
		"omni_createpayload_simplesend":    {handler: omni_createpayload_simplesend},
		"omni_createpayload_issuancefixed": {handler: omni_createpayload_issuancefixed},
		"omni_listproperties":              {handler: omni_listproperties},

		"omni_sendissuancefixed": {handler: omniSendIssuanceFixed},
		"omni_getbalance":        {handler: omniGetBalance},
		"omni_send":              {handler: omniSend},
	}
}

func omni_getinfo(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return OmnCmdReq(icmd, w)
}

func omni_createpayload_simplesend(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	cmd := icmd.(*hcjson.OmniCreatepayloadSimplesendCmd)
	byteCmd, err := hcjson.MarshalCmd(1, cmd)
	if err != nil {
		return err, nil
	}
	strReq := string(byteCmd)
	strRsp := omnilib.JsonCmdReqHcToOm(strReq)

	var response hcjson.Response
	_ = json.Unmarshal([]byte(strRsp), &response)

	return response.Result, nil
	//return w.Locked(), nil
}

func omni_createpayload_issuancefixed(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return OmnCmdReq(icmd, w)
}

func omni_listproperties(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return OmnCmdReq(icmd, w)
}

func omniSendIssuanceFixed(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return OmnCmdReq(icmd, w)

	/*
		if err != nil {
				return "", err
			}

			switch v := msg.(type) {
			case json.RawMessage:
				payload, err := v.MarshalJSON()
				if err != nil {
					return "", err
				}
				fmt.Println("omniSendIssuanceFixed:", string(payload))
				payload = payload[1:len(payload)-1]
				fmt.Println("omniSendIssuanceFixed:", string(payload))

				return sendIssuanceFixed(w, []byte(payload))
			default:
				fmt.Printf("%T", msg)
				return "", fmt.Errorf("data from omni err type:%T", msg)
			}
	*/
}

//
func sendIssuanceFixed(w *wallet.Wallet, payLoad []byte) (string, error) {
	account := uint32(udb.DefaultAccountNum)

	var changeAddr string
	addr, err := w.FirstAddr(account)
	if err != nil {
		return "", err
	}
	changeAddr = addr.String()
	dstAddr := changeAddr

	amt, err := hcutil.NewAmount(20)
	if err != nil {
		return "", err
	}
	// Mock up map of address and amount pairs.
	pairs := map[string]hcutil.Amount{
		dstAddr: hcutil.Amount(amt),
	}

	// sendtoaddress always spends from the default account, this matches bitcoind
	return sendPairsWithPayLoad(w, pairs, account, 1, changeAddr, payLoad)
}

// sendPairsWithPayLoad creates and sends payment transactions.
// It returns the transaction hash in string format upon success
// All errors are returned in hcjson.RPCError format
func sendPairsWithPayLoad(w *wallet.Wallet, amounts map[string]hcutil.Amount,
	account uint32, minconf int32, changeAddr string, payLoad []byte) (string, error) {
	outputs, err := makeOutputs(amounts, w.ChainParams())
	if err != nil {
		return "", err
	}
	payloadNullDataOutput, err := w.MakeNulldataOutput(payLoad)
	if err != nil {
		return "", err
	}
	outputs = append(outputs, payloadNullDataOutput)

	txSha, err := w.SendOutputs(outputs, account, minconf, changeAddr, "")
	if err != nil {
		if err == txrules.ErrAmountNegative {
			return "", ErrNeedPositiveAmount
		}
		if apperrors.IsError(err, apperrors.ErrLocked) {
			return "", &ErrWalletUnlockNeeded
		}
		switch err.(type) {
		case hcjson.RPCError:
			return "", err
		}

		return "", &hcjson.RPCError{
			Code:    hcjson.ErrRPCInternal.Code,
			Message: err.Error(),
		}
	}

	return txSha.String(), err
}

type SendFromAddressToAddressCmd struct {
	FromAddress   string
	ToAddress     string
	ChangeAddress string
	Amount        float64
	Comment       *string
	CommentTo     *string
}

func omniSendToAddress(cmd *SendFromAddressToAddressCmd, w *wallet.Wallet, payLoad []byte) (string, error) {

	// Transaction comments are not yet supported.  Error instead of
	// pretending to save them.
	if !isNilOrEmpty(cmd.Comment) || !isNilOrEmpty(cmd.CommentTo) {
		return "", &hcjson.RPCError{
			Code:    hcjson.ErrRPCUnimplemented,
			Message: "Transaction comments are not yet supported",
		}
	}

	account := uint32(udb.DefaultAccountNum)
	amt, err := hcutil.NewAmount(cmd.Amount)
	if err != nil {
		return "", err
	}

	// Check that signed integer parameters are positive.
	if amt < 0 {
		return "", ErrNeedPositiveAmount
	}

	// Mock up map of address and amount pairs.
	pairs := map[string]hcutil.Amount{
		cmd.ToAddress: amt,
	}

	// sendtoaddress always spends from the default account, this matches bitcoind
	return sendPairs(w, pairs, account, 1, cmd.ChangeAddress, payLoad, cmd.FromAddress)
}

func omniGetBalance(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return OmnCmdReq(icmd, w)
}

func omniSend(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	msg, err := OmnCmdReq(icmd, w)
	if err != nil {
		return nil, err
	}
	switch v := msg.(type) {
	case json.RawMessage:
		payload, err := v.MarshalJSON()
		if err != nil {
			return "", err
		}

		payload = payload[1 : len(payload)-1]

		//

		return sendIssuanceFixed(w, []byte(payload))
	default:
		fmt.Printf("%T", msg)
		return "", fmt.Errorf("data from omni err type:%T", msg)
	}
}
