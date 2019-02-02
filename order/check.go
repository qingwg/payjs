package order

import (
	"encoding/json"
	"fmt"
	"github.com/qingwg/payjs/util"
)

const getCheckURL = "https://payjs.cn/api/check"

// CheckRequest 请求参数
type CheckRequest struct {
	PayJSOrderID string `json:"payjs_order_id"` //Y	PAYJS 平台订单号
	Sign         string `json:"sign"`           //Y	数据签名 详见签名算法
}

// CheckResponse PayJS返回参数
type CheckResponse struct {
	ReturnCode    int    `json:"return_code"`    //Y	1:请求成功 0:请求失败
	ReturnMsg     string `json:"return_msg"`     //Y	返回消息
	MchID         string `json:"mchid"`          //Y	PAYJS 平台商户号
	OutTradeNo    string `json:"out_trade_no"`   //Y	用户端订单号
	PayJSOrderID  string `json:"payjs_order_id"` //Y	PAYJS 订单号
	TransactionID string `json:"transaction_id"` //N	微信显示订单号
	Status        int    `json:"status"`         //Y	0：未支付，1：支付成功
	Openid        string `json:"openid"`         //N	用户 OPENID
	TotalFee      int64  `json:"total_fee"`      //N	订单金额
	PaidTime      string `json:"paid_time"`      //N	订单支付时间
	Attach        string `json:"attach"`         //N	用户自定义数据
	Sign          string `json:"sign"`           //Y	数据签名 详见签名算法
}

// Check 用户发起支付后，可通过本接口发起订单查询来确认订单状态
func (order *Order) Check(payJSOrderID string) (checkResponse CheckResponse, err error) {
	checkRequest := CheckRequest{
		PayJSOrderID: payJSOrderID,
	}
	sign := util.Signature(checkRequest, order.Key)
	checkRequest.Sign = sign
	response, err := util.PostJSON(getCheckURL, checkRequest)
	if err != nil {
		return
	}

	err = json.Unmarshal(response, &checkResponse)
	if err != nil {
		return
	}
	if checkResponse.ReturnCode == 0 {
		err = fmt.Errorf("OrderCheck Error , errcode=%d, errmsg=%s", checkResponse.ReturnCode, checkResponse.ReturnMsg)
		return
	}
	// 检测sign
	msgSignature := checkResponse.Sign
	msgSignatureGen := util.Signature(checkResponse, order.Key)
	if msgSignature != msgSignatureGen {
		err = fmt.Errorf("消息不合法，验证签名失败")
		return
	}
	return
}
