package micropay

import (
	"encoding/json"
	"fmt"
	"github.com/qingwg/payjs/context"
	"github.com/qingwg/payjs/util"
)

const getCreateURL = "https://payjs.cn/api/micropay"

// Micropay struct
type Micropay struct {
	*context.Context
}

// CreateRequest 请求参数
type CreateRequest struct {
	MchID      string `json:"mchid"`        //Y	商户号
	TotalFee   int64  `json:"total_fee"`    //Y	金额。单位：分
	OutTradeNo string `json:"out_trade_no"` //Y	用户端自主生成的订单号
	Body       string `json:"body"`         //N	订单标题
	Attach     string `json:"attach"`       //N	用户自定义数据，在notify的时候会原样返回
	AuthCode   string `json:"auth_code"`    //Y	扫码支付授权码，设备读取用户微信中的条码或者二维码信息(注：用户刷卡条形码规则：18位纯数字，以10、11、12、13、14、15开头)
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
	Status       int    `json:"status"`         //Y	0：未支付，1：支付成功（以后会取消）
	Sign         string `json:"sign"`           //Y	数据签名 详见签名算法
}

//NewMicropay init
func NewMicropay(context *context.Context) *Micropay {
	micropay := new(Micropay)
	micropay.Context = context
	return micropay
}

// Create 拿到扫码信息请求PayJS
func (micropay *Micropay) Create(totalFeeReq int64, bodyReq, outTradeNoReq, attachReq, autoCodeReq string) (createResponse CreateResponse, err error) {
	createRequest := CreateRequest{
		MchID:      micropay.MchID,
		TotalFee:   totalFeeReq,
		OutTradeNo: outTradeNoReq,
		Body:       bodyReq,
		Attach:     attachReq,
		AuthCode:   autoCodeReq,
	}
	sign := util.Signature(createRequest, micropay.Key)
	createRequest.Sign = sign
	response, err := util.PostJSON(getCreateURL, createRequest)
	if err != nil {
		return
	}

	err = json.Unmarshal(response, &createResponse)
	if err != nil {
		return
	}

	if createResponse.ReturnCode == 0 {
		err = fmt.Errorf("MicropayCreate Error , errcode=%d , errmsg=%s, errmsg=%s", createResponse.ReturnCode, createResponse.Msg, createResponse.ReturnMsg)
		return
	}

	// 检测sign
	msgSignature := createResponse.Sign
	msgSignatureGen := util.Signature(createResponse, micropay.Key)
	if msgSignature != msgSignatureGen {
		err = fmt.Errorf("消息不合法，验证签名失败")
		return
	}
	return
}
