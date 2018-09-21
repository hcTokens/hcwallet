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
	"strings"

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

		"omni_senddexsell":                       {handler: OmniSenddexsell},
		"omni_senddexaccept":                     {handler: OmniSenddexaccept},
		"omni_sendissuancecrowdsale":             {handler: OmniSendissuancecrowdsale},
		"omni_sendissuancemanaged":               {handler: OmniSendissuancemanaged},
		"omni_sendsto":                           {handler: OmniSendsto},
		"omni_sendgrant":                         {handler: OmniSendgrant},
		"omni_sendrevoke":                        {handler: OmniSendrevoke},
		"omni_sendclosecrowdsale":                {handler: OmniSendclosecrowdsale},
		"omni_sendtrade":                         {handler: OmniSendtrade},
		"omni_sendcanceltradesbyprice":           {handler: OmniSendcanceltradesbyprice},
		"omni_sendcanceltradesbypair":            {handler: OmniSendcanceltradesbypair},
		"omni_sendcancelalltrades":               {handler: OmniSendcancelalltrades},
		"omni_sendchangeissuer":                  {handler: OmniSendchangeissuer},
		"omni_sendall":                           {handler: OmniSendall},
		"omni_sendenablefreezing":                {handler: OmniSendenablefreezing},
		"omni_senddisablefreezing":               {handler: OmniSenddisablefreezing},
		"omni_sendfreeze":                        {handler: OmniSendfreeze},
		"omni_sendunfreeze":                      {handler: OmniSendunfreeze},
		"omni_sendrawtx":                         {handler: OmniSendrawtx},
		"omni_funded_send":                       {handler: OmniFundedSend},
		"omni_funded_sendall":                    {handler: OmniFundedSendall},
		"omni_getallbalancesforid":               {handler: OmniGetallbalancesforid},
		"omni_getallbalancesforaddress":          {handler: OmniGetallbalancesforaddress},
		"omni_getwalletbalances":                 {handler: OmniGetwalletbalances},
		"omni_getwalletaddressbalances":          {handler: OmniGetwalletaddressbalances},
		"omni_gettransaction":                    {handler: OmniGettransaction},
		"omni_listtransactions":                  {handler: OmniListtransactions},
		"omni_listblocktransactions":             {handler: OmniListblocktransactions},
		"omni_listpendingtransactions":           {handler: OmniListpendingtransactions},
		"omni_getactivedexsells":                 {handler: OmniGetactivedexsells},
		"omni_getproperty":                       {handler: OmniGetproperty},
		"omni_getactivecrowdsales":               {handler: OmniGetactivecrowdsales},
		"omni_getcrowdsale":                      {handler: OmniGetcrowdsale},
		"omni_getgrants":                         {handler: OmniGetgrants},
		"omni_getsto":                            {handler: OmniGetsto},
		"omni_gettrade":                          {handler: OmniGettrade},
		"omni_getorderbook":                      {handler: OmniGetorderbook},
		"omni_gettradehistoryforpair":            {handler: OmniGettradehistoryforpair},
		"omni_gettradehistoryforaddress":         {handler: OmniGettradehistoryforaddress},
		"omni_getactivations":                    {handler: OmniGetactivations},
		"omni_getpayload":                        {handler: OmniGetpayload},
		"omni_getseedblocks":                     {handler: OmniGetseedblocks},
		"omni_getcurrentconsensushash":           {handler: OmniGetcurrentconsensushash},
		"omni_decodetransaction":                 {handler: OmniDecodetransaction},
		"omni_createrawtx_opreturn":              {handler: OmniCreaterawtxOpreturn},
		"omni_createrawtx_multisig":              {handler: OmniCreaterawtxMultisig},
		"omni_createrawtx_input":                 {handler: OmniCreaterawtxInput},
		"omni_createrawtx_reference":             {handler: OmniCreaterawtxReference},
		"omni_createrawtx_change":                {handler: OmniCreaterawtxChange},
		"omni_createpayload_sendall":             {handler: OmniCreatepayloadSendall},
		"omni_createpayload_dexsell":             {handler: OmniCreatepayloadDexsell},
		"omni_createpayload_dexaccept":           {handler: OmniCreatepayloadDexaccept},
		"omni_createpayload_sto":                 {handler: OmniCreatepayloadSto},
		"omni_createpayload_issuancecrowdsale":   {handler: OmniCreatepayloadIssuancecrowdsale},
		"omni_createpayload_issuancemanaged":     {handler: OmniCreatepayloadIssuancemanaged},
		"omni_createpayload_closecrowdsale":      {handler: OmniCreatepayloadClosecrowdsale},
		"omni_createpayload_grant":               {handler: OmniCreatepayloadGrant},
		"omni_createpayload_revoke":              {handler: OmniCreatepayloadRevoke},
		"omni_createpayload_changeissuer":        {handler: OmniCreatepayloadChangeissuer},
		"omni_createpayload_trade":               {handler: OmniCreatepayloadTrade},
		"omni_createpayload_canceltradesbyprice": {handler: OmniCreatepayloadCanceltradesbyprice},
		"omni_createpayload_canceltradesbypair":  {handler: OmniCreatepayloadCanceltradesbypair},
		"omni_createpayload_cancelalltrades":     {handler: OmniCreatepayloadCancelalltrades},
		"omni_createpayload_enablefreezing":      {handler: OmniCreatepayloadEnablefreezing},
		"omni_createpayload_disablefreezing":     {handler: OmniCreatepayloadDisablefreezing},
		"omni_createpayload_freeze":              {handler: OmniCreatepayloadFreeze},
		"omni_createpayload_unfreeze":            {handler: OmniCreatepayloadUnfreeze},
		"omni_getfeecache":                       {handler: OmniGetfeecache},
		"omni_getfeetrigger":                     {handler: OmniGetfeetrigger},
		"omni_getfeeshare":                       {handler: OmniGetfeeshare},
		"omni_getfeedistribution":                {handler: OmniGetfeedistribution},
		"omni_getfeedistributions":               {handler: OmniGetfeedistributions},
		"omni_setautocommit":                     {handler: OmniSetautocommit},
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
	fmt.Printf(strReq)
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
		return nil, err
	}
	hexStr := strings.Trim(string(ret), "\"")
	payLoad, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}
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
		ChangeAddress: omniSendCmd.Fromaddress,
		ToAddress:     omniSendCmd.Toaddress,
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

func OmniSend(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniSenddexsell(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniSenddexaccept(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniSendissuancecrowdsale(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniSendissuancefixed(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniSendissuancemanaged(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniSendsto(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniSendgrant(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniSendrevoke(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniSendclosecrowdsale(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniSendtrade(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniSendcanceltradesbyprice(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniSendcanceltradesbypair(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniSendcancelalltrades(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniSendchangeissuer(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniSendall(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniSendenablefreezing(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniSenddisablefreezing(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniSendfreeze(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniSendunfreeze(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniSendrawtx(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniFundedSend(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniFundedSendall(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniGetinfo(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniGetbalance(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniGetallbalancesforid(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniGetallbalancesforaddress(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniGetwalletbalances(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniGetwalletaddressbalances(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniGettransaction(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniListtransactions(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniListblocktransactions(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniListpendingtransactions(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniGetactivedexsells(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniListproperties(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniGetproperty(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniGetactivecrowdsales(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniGetcrowdsale(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniGetgrants(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniGetsto(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniGettrade(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniGetorderbook(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniGettradehistoryforpair(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniGettradehistoryforaddress(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniGetactivations(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniGetpayload(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniGetseedblocks(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniGetcurrentconsensushash(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniDecodetransaction(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniCreaterawtxOpreturn(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniCreaterawtxMultisig(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniCreaterawtxInput(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniCreaterawtxReference(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniCreaterawtxChange(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniCreatepayloadSimplesend(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniCreatepayloadSendall(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniCreatepayloadDexsell(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniCreatepayloadDexaccept(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniCreatepayloadSto(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniCreatepayloadIssuancefixed(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniCreatepayloadIssuancecrowdsale(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniCreatepayloadIssuancemanaged(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniCreatepayloadClosecrowdsale(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniCreatepayloadGrant(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniCreatepayloadRevoke(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniCreatepayloadChangeissuer(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniCreatepayloadTrade(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniCreatepayloadCanceltradesbyprice(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniCreatepayloadCanceltradesbypair(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniCreatepayloadCancelalltrades(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniCreatepayloadEnablefreezing(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniCreatepayloadDisablefreezing(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniCreatepayloadFreeze(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniCreatepayloadUnfreeze(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniGetfeecache(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniGetfeetrigger(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniGetfeeshare(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniGetfeedistribution(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniGetfeedistributions(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}

func OmniSetautocommit(icmd interface{}, w *wallet.Wallet) (interface{}, error) {
	return omni_cmdReq(icmd, w)
}
