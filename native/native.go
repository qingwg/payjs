package native

import (
	"encoding/json"
	"fmt"
	"github.com/qingwg/payjs/context"
	"github.com/qingwg/payjs/util"
)

const getCreateURL = "https://payjs.cn/api/native"

// Native struct
type Native struct {
	*context.Context
}

// CreateRequest 请求参数
type CreateRequest struct {
	MchID      string `json:"mchid"`        //Y	商户号
	TotalFee   int64  `json:"total_fee"`    //Y	金额。单位：分
	OutTradeNo string `json:"out_trade_no"` //Y	用户端自主生成的订单号
	Type       string `json:"type"`         //N 留空表示微信支付。支付宝交易传值：alipay
	Body       string `json:"body"`         //N	订单标题
	Attach     string `json:"attach"`       //N	用户自定义数据，在notify的时候会原样返回
	NotifyUrl  string `json:"notify_url"`   //N	接收微信支付异步通知的回调地址。必须为可直接访问的URL，不能带参数、session验证、csrf验证。留空则不通知
	Sign       string `json:"sign"`         //Y	数据签名 详见签名算法
}

// CreateResponse PayJS返回参数
type CreateResponse struct {
	ReturnCode   int    `json:"return_code"`    //Y	1:请求成功，0:请求失败
	Msg          string `json:"msg"`            //N	return_code为0时返回的错误消息
	ReturnMsg    string `json:"return_msg"`     //Y	返回消息
	PayJSOrderID string `json:"payjs_order_id"` //Y	PAYJS 平台订单号
	OutTradeNo   string `json:"out_trade_no"`   //Y	用户生成的订单号原样返回
	TotalFee     int64  `json:"total_fee"`      //Y	金额。单位：分
	Qrcode       string `json:"qrcode"`         //Y	二维码图片地址
	CodeUrl      string `json:"code_url"`       //Y	可将该参数生成二维码展示出来进行扫码支付
	Sign         string `json:"sign"`           //Y	数据签名 详见签名算法
}

//NewNative init
func NewNative(context *context.Context) *Native {
	native := new(Native)
	native.Context = context
	return native
}

// Create 请求PayJS获取支付二维码
func (native *Native) Create(totalFeeReq int64, bodyReq, outTradeNoReq, attachReq, payType string) (createResponse CreateResponse, err error) {
	createRequest := CreateRequest{
		MchID:      native.MchID,
		TotalFee:   totalFeeReq,
		OutTradeNo: outTradeNoReq,
		Type:       payType,
		Body:       bodyReq,
		Attach:     attachReq,
		NotifyUrl:  native.NotifyUrl,
	}
	sign := util.Signature(createRequest, native.Key)
	createRequest.Sign = sign
	response, err := util.PostJSON(getCreateURL, createRequest)
	if err != nil {
		return
	}

	err = json.Unmarshal(response, &createResponse)
	if err != nil {
		return
	}
	if createResponse.ReturnCode != 1 {
		err = fmt.Errorf("NativeCreate Error , errcode=%v , errmsg=%s, errmsg=%s", createResponse.ReturnCode, createResponse.Msg, createResponse.ReturnMsg)
		return
	}
	// 检测sign
	msgSignature := createResponse.Sign
	msgSignatureGen := util.Signature(createResponse, native.Key)
	if msgSignature != msgSignatureGen {
		err = fmt.Errorf("消息不合法，验证签名失败")
		return
	}
	return
}
