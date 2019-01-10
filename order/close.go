package order

import (
	"encoding/json"
	"fmt"
	"github.com/qingwg/payjs/util"
)

const getCloseURL = "https://payjs.cn/api/close"

// CloseRequest 请求参数
type CloseRequest struct {
	PayJSOrderID string `json:"payjs_order_id"` //Y	PAYJS 平台订单号
	Sign         string `json:"sign"`           //Y	数据签名 详见签名算法
}

// CloseResponse PayJS返回参数
type CloseResponse struct {
	ReturnCode   int    `json:"return_code"`    //Y	1:请求成功 0:请求失败
	ReturnMsg    string `json:"return_msg"`     //Y	返回消息
	PayJSOrderID string `json:"payjs_order_id"` //Y	PAYJS 平台订单号
	Sign         string `json:"sign"`           //Y	数据签名 详见签名算法
}

// Close 关闭已经发起的订单
func (order *Order) Close(payJSOrderID string) (closeResponse CloseResponse, err error) {
	closeRequest := CloseRequest{
		PayJSOrderID: payJSOrderID,
	}
	sign := util.Signature(closeRequest, order.Key)
	closeRequest.Sign = sign
	response, err := util.PostJSON(getCloseURL, closeRequest)
	if err != nil {
		return
	}
	err = json.Unmarshal(response, &closeResponse)
	if err != nil {
		return
	}
	if closeResponse.ReturnCode == 0 {
		err = fmt.Errorf("OrderClose Error , errcode=%d , errmsg=%s", closeResponse.ReturnCode, closeResponse.ReturnMsg)
		return
	}
	// 检测sign
	msgSignature := closeResponse.Sign
	msgSignatureGen := util.Signature(closeResponse, order.Key)
	if msgSignature != msgSignatureGen {
		err = fmt.Errorf("消息不合法，验证签名失败")
		return
	}
	return
}
