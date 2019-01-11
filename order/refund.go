package order

import (
	"encoding/json"
	"fmt"
	"github.com/qingwg/payjs/util"
)

const getRefundURL = "https://payjs.cn/api/refund"

// RefundRequest 请求参数
type RefundRequest struct {
	PayJSOrderID string `json:"payjs_order_id"` //Y	PAYJS 平台订单号
	Sign         string `json:"sign"`           //Y	数据签名 详见签名算法
}

// RefundResponse PayJS返回参数
type RefundResponse struct {
	ReturnCode    int    `json:"return_code"`    //Y	1:请求成功 0:请求失败
	ReturnMsg     string `json:"return_msg"`     //Y	返回消息
	PayJSOrderID  string `json:"payjs_order_id"` //Y	PAYJS 平台订单号
	OutTradeNo    string `json:"out_trade_no"`   //N	用户侧订单号
	TransactionID string `json:"transaction_id"` //N	微信支付订单号
	Sign          string `json:"sign"`           //Y	数据签名 详见签名算法
}

// Refund 对已经支付的订单发起退款
func (order *Order) Refund(payJSOrderID string) (refundResponse RefundResponse, err error) {
	refundRequest := RefundRequest{
		PayJSOrderID: payJSOrderID,
	}
	sign := util.Signature(refundRequest, order.Key)
	refundRequest.Sign = sign
	response, err := util.PostJSON(getRefundURL, refundRequest)
	if err != nil {
		return
	}
	err = json.Unmarshal(response, &refundResponse)
	if err != nil {
		return
	}
	if refundResponse.ReturnCode == 0 {
		err = fmt.Errorf("OrderRefund Error , errcode=%d , errmsg=%s", refundResponse.ReturnCode, refundResponse.ReturnMsg)
		return
	}
	// 检测sign
	msgSignature := refundResponse.Sign
	msgSignatureGen := util.Signature(refundResponse, order.Key)
	if msgSignature != msgSignatureGen {
		err = fmt.Errorf("消息不合法，验证签名失败")
		return
	}
	return
}
