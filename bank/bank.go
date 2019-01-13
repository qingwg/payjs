package bank

import (
	"encoding/json"
	"fmt"
	"github.com/qingwg/payjs/context"
	"github.com/qingwg/payjs/util"
)

const getBankInfoURL = "https://payjs.cn/api/bank"

// Bank struct
type Bank struct {
	*context.Context
}

// BankInfoRequest 请求参数
type BankInfoRequest struct {
	MchID string `json:"mchid"` //Y	商户号
	Bank  string `json:"bank"`  //Y	银行简写
	Sign  string `json:"sign"`  //Y	数据签名 详见签名算法
}

// BankInfoResponse PayJS返回参数
type BankInfoResponse struct {
	ReturnCode int    `json:"return_code"` //Y	1:请求成功 0:请求失败
	ReturnMsg  string `json:"return_msg"`  //Y	返回消息
	Bank       string `json:"bank"`        //Y	银行名称
	Sign       string `json:"sign"`        //Y	数据签名 详见签名算法
}

//NewBank init
func NewBank(context *context.Context) *Bank {
	bank := new(Bank)
	bank.Context = context
	return bank
}

// GetBankInfo 根据银行简称查询银行详细名称。银行数据库会随时更新
func (bank *Bank) GetBankInfo(bankReq string) (bankInfoResponse BankInfoResponse, err error) {
	bankInfoRequest := BankInfoRequest{
		MchID: bank.MchID,
		Bank:  bankReq,
	}
	sign := util.Signature(bankInfoRequest, bank.Key)
	bankInfoRequest.Sign = sign
	response, err := util.PostJSON(getBankInfoURL, bankInfoRequest)
	if err != nil {
		return
	}
	err = json.Unmarshal(response, &bankInfoResponse)
	if err != nil {
		return
	}
	if bankInfoResponse.ReturnCode == 0 {
		err = fmt.Errorf("GetBankInfo Error , errcode=%d , errmsg=%s", bankInfoResponse.ReturnCode, bankInfoResponse.ReturnMsg)
		return
	}
	// 检测sign
	msgSignature := bankInfoResponse.Sign
	msgSignatureGen := util.Signature(bankInfoResponse, bank.Key)
	if msgSignature != msgSignatureGen {
		err = fmt.Errorf("消息不合法，验证签名失败")
		return
	}
	return
}
