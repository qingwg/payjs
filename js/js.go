package js

import (
	"encoding/json"
	"fmt"
	"github.com/qingwg/payjs/context"
	"github.com/qingwg/payjs/util"
)

const getJsApiURL = "https://payjs.cn/api/jsapi"

// Js struct
type Js struct {
	*context.Context
}

// JsApiRequest
type JsApiRequest struct {
	MchID      string `json:"mchid"`        //Y	商户号
	TotalFee   int64  `json:"total_fee"`    //Y	金额。单位：分
	OutTradeNo string `json:"out_trade_no"` //Y	用户端自主生成的订单号，在用户端要保证唯一性
	Body       string `json:"body"`         //N	订单标题
	Attach     string `json:"attach"`       //N	用户自定义数据，在notify的时候会原样返回
	NotifyUrl  string `json:"notify_url"`   //N	接收微信支付异步通知的回调地址。必须为可直接访问的URL，不能带参数、session验证、csrf验证。留空则不通知
	Openid     string `json:"openid"`       //Y	用户openid
	Sign       string `json:"sign"`         //Y	数据签名 详见签名算法
}

// JsApiResponse
type JsApiResponse struct {
	ReturnCode   int    `json:"return_code"`    //Y	0:失败 1:成功
	ReturnMsg    string `json:"return_msg"`     //Y	失败原因
	PayJSOrderID string `json:"payjs_order_id"` //Y	PAYJS 侧订单号
	JsApi        JsApi  `json:"jsapi"`          //N	用于发起支付的支付参数
	Sign         string `json:"sign"`           //Y	数据签名
}

// JsApi
type JsApi struct {
	AppID     string `json:"appId"`
	TimeStamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign"`
}

//NewJs init
func NewJs(context *context.Context) *Js {
	js := new(Js)
	js.Context = context
	return js
}

// Create 获取发起支付所需要的参数
func (js *Js) Create(totalFeeReq int64, bodyReq, outTradeNoReq, attachReq, openid string) (jsApiResponse JsApiResponse, err error) {
	jsApiRequest := JsApiRequest{
		MchID:      js.MchID,
		TotalFee:   totalFeeReq,
		OutTradeNo: outTradeNoReq,
		Body:       bodyReq,
		Attach:     attachReq,
		NotifyUrl:  js.NotifyUrl,
		Openid:     openid,
	}
	sign := util.Signature(jsApiRequest, js.Key)
	jsApiRequest.Sign = sign
	response, err := util.PostJSON(getJsApiURL, jsApiRequest)
	if err != nil {
		return
	}
	//fmt.Println("=====response", string(response))

	err = json.Unmarshal(response, &jsApiResponse)
	if err != nil {
		return
	}
	if jsApiResponse.ReturnCode == 0 {
		err = fmt.Errorf("GetJsApi Error , errcode=%d , errmsg=%s", jsApiResponse.ReturnCode, jsApiResponse.ReturnMsg)
		return
	}
	//todo:解决多维结构造成签名验证失败bug
	//// 检测sign
	//msgSignature := jsApiResponse.Sign
	//msgSignatureGen := util.Signature(jsApiResponse, js.Key)
	//if msgSignature != msgSignatureGen {
	//	err = fmt.Errorf("消息不合法，验证签名失败")
	//	return
	//}
	return
}
