// Copyright (c) 2013-2016 The btcsuite developers
// Copyright (c) 2015-2017 The Decred developers
// Copyright (c) 2018-2020 The Hc developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package legacyrpc

import (
	"encoding/hex"
	"encoding/json"

	"github.com/HcashOrg/hcd/hcjson"
	"github.com/HcashOrg/hcd/hcutil"
	"github.com/HcashOrg/hcwallet/apperrors"
	"github.com/HcashOrg/hcwallet/omnilib"
	"github.com/HcashOrg/hcwallet/wallet"
	"github.com/HcashOrg/hcwallet/wallet/txrules"
	"github.com/HcashOrg/hcwallet/wallet/udb"
)

const (
	MininumAmount = 1000000
)

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

//add by ycj 20180915
//commonly used cmd request
func omni_cmdReq(icmd interface{}, w *wallet.Wallet) (json.RawMessage, error) {
	byteCmd, err := hcjson.MarshalCmd(1, icmd)
	if err != nil {
		return nil, err
	}
	strReq := string(byteCmd)
	strRsp := omnilib.JsonCmdReqHcToOm(strReq)

	var response hcjson.Response
	_ = json.Unmarshal([]byte(strRsp), &response)
	return response.Result, nil
}

//
func omni_getinfo(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
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
}

func omni_createpayload_issuancefixed(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func omni_listproperties(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func omniSendIssuanceFixed(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	txIdBytes, err := omni_cmdReq(icmd, w)
	sendIssueCmd := icmd.(*hcjson.OmniSendissuancefixedCmd)
	if err != nil {
		return err, nil
	}

	txidStr := ""
	err = json.Unmarshal(txIdBytes, &txidStr)
	if err != nil {
		return err, nil
	}

	payLoad, err := hex.DecodeString(txidStr)
	if err != nil {
		return err, nil
	}

	sendParams := &SendFromAddressToAddress{
		FromAddress:   sendIssueCmd.Fromaddress,
		ToAddress:     sendIssueCmd.Fromaddress,
		ChangeAddress: sendIssueCmd.Fromaddress,
	}
	return omniSendToAddress(sendParams, w, payLoad)
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
	return sendPairsWithPayLoad(w, pairs, account, 1, changeAddr, payLoad, "")
}

func omniGetBalance(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func omniSend(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	omniSendCmd := icmd.(*hcjson.OmniSendCmd)
	ret, err := omni_cmdReq(icmd, w)
	if err != nil {
		return ret, err
	}

	payLoad, err := hex.DecodeString(string(ret))

	_, err = decodeAddress(omniSendCmd.Fromaddress, w.ChainParams())
	if err != nil {
		return nil, err
	}

	_, err = decodeAddress(omniSendCmd.Toaddress, w.ChainParams())
	if err != nil {
		return nil, err
	}

	cmd := &SendFromAddressToAddress{
		FromAddress:   omniSendCmd.Fromaddress,
		ToAddress:     omniSendCmd.Toaddress,
		ChangeAddress: omniSendCmd.Fromaddress,
		Amount:        1,
	}
	return omniSendToAddress(cmd, w, payLoad)
}

type SendFromAddressToAddress struct {
	FromAddress   string
	ToAddress     string
	ChangeAddress string
	Amount        float64
	Comment       *string
	CommentTo     *string
}

func omniSendToAddress(cmd *SendFromAddressToAddress, w *wallet.Wallet, payLoad []byte) (string, error) {
	// Transaction comments are not yet supported.  Error instead of
	// pretending to save them.
	if !isNilOrEmpty(cmd.Comment) || !isNilOrEmpty(cmd.CommentTo) {
		return "", &hcjson.RPCError{
			Code:    hcjson.ErrRPCUnimplemented,
			Message: "Transaction comments are not yet supported",
		}
	}

	account := uint32(udb.DefaultAccountNum)

	// Mock up map of address and amount pairs.
	pairs := map[string]hcutil.Amount{
		cmd.ToAddress: MininumAmount,
	}

	return sendPairsWithPayLoad(w, pairs, account, 1, cmd.ChangeAddress, payLoad, cmd.FromAddress)
}

// sendPairsWithPayLoad creates and sends payment transactions.
// It returns the transaction hash in string format upon success
// All errors are returned in hcjson.RPCError format
func sendPairsWithPayLoad(w *wallet.Wallet, amounts map[string]hcutil.Amount, account uint32, minconf int32, changeAddr string, payLoad []byte, fromAddress string) (string, error) {
	outputs, err := makeOutputs(amounts, w.ChainParams())
	if err != nil {
		return "", err
	}
	payloadNullDataOutput, err := w.MakeNulldataOutput(payLoad)
	if err != nil {
		return "", err
	}

	outputs = append(outputs, payloadNullDataOutput)

	txSha, err := w.SendOutputs(outputs, account, minconf, changeAddr, fromAddress)
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
