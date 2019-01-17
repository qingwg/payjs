package order

import (
	"encoding/json"
	"fmt"
	"github.com/qingwg/payjs/util"
)

const getReverseURL = "https://payjs.cn/api/reverse"

// ReverseRequest 请求参数
type ReverseRequest struct {
	PayJSOrderID string `json:"payjs_order_id"` //Y	PAYJS 平台订单号
	Sign         string `json:"sign"`           //Y	数据签名 详见签名算法
}

// ReverseResponse PayJS返回参数
type ReverseResponse struct {
	ReturnCode   int    `json:"return_code"`    //Y	1:请求成功 0:请求失败
	ReturnMsg    string `json:"return_msg"`     //Y	返回消息
	PayJSOrderID string `json:"payjs_order_id"` //Y	PAYJS 平台订单号
	Sign         string `json:"sign"`           //Y	数据签名 详见签名算法
}

// Reverse 撤销订单
func (order *Order) Reverse(payJSOrderID string) (reverseResponse ReverseResponse, err error) {
	reverseRequest := ReverseRequest{
		PayJSOrderID: payJSOrderID,
	}
	sign := util.Signature(reverseRequest, order.Key)
	reverseRequest.Sign = sign
	response, err := util.PostJSON(getReverseURL, reverseRequest)
	if err != nil {
		return
	}

	err = json.Unmarshal(response, &reverseResponse)
	if err != nil {
		return
	}
	if reverseResponse.ReturnCode == 0 {
		err = fmt.Errorf("OrderReverse  Error , errcode=%d , errmsg=%s", reverseResponse.ReturnCode, reverseResponse.ReturnMsg)
		return
	}
	// 检测sign
	msgSignature := reverseResponse.Sign
	msgSignatureGen := util.Signature(reverseResponse, order.Key)
	if msgSignature != msgSignatureGen {
		err = fmt.Errorf("消息不合法，验证签名失败")
		return
	}
	return
}
