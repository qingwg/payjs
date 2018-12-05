package micropay

import (
	"encoding/json"
	"fmt"
	"github.com/yuyan2077/payjs/context"
	"github.com/yuyan2077/payjs/util"
)

const getMicropayURL = "https://payjs.cn/api/cashier"

// Micropay struct
type Micropay struct {
	*context.Context
}

// MicropayRequest 请求参数
type MicropayRequest struct {
	MchID      string `json:"mchid"`        //Y	商户号
	TotalFee   int    `json:"total_fee"`    //Y	金额。单位：分
	OutTradeNo string `json:"out_trade_no"` //Y	用户端自主生成的订单号
	Body       string `json:"body"`         //N	订单标题
	Attach     string `json:"attach"`       //N	用户自定义数据，在notify的时候会原样返回
	AuthCode   string `json:"auth_code"`    //Y	扫码支付授权码，设备读取用户微信中的条码或者二维码信息(注：用户刷卡条形码规则：18位纯数字，以10、11、12、13、14、15开头)
	Sign       string `json:"sign"`         //Y	数据签名 详见签名算法
}

// MicropayResponse PayJS返回参数
type MicropayResponse struct {
	ReturnCode   int    `json:"return_code"`    //Y	1:请求成功，0:请求失败
	Status       int    `json:"status"`         //N	return_code为0时有status参数为0
	Msg          string `json:"msg"`            //N	return_code为0时返回的错误消息
	ReturnMsg    string `json:"return_msg"`     //Y	返回消息
	PayJSOrderID string `json:"payjs_order_id"` //Y	PAYJS 平台订单号
	OutTradeNo   string `json:"out_trade_no"`   //Y	用户生成的订单号原样返回
	TotalFee     string `json:"total_fee"`      //Y	金额。单位：分
	Sign         string `json:"sign"`           //Y	数据签名 详见签名算法
}

//NewMicropay init
func NewMicropay(context *context.Context) *Micropay {
	micropay := new(Micropay)
	micropay.Context = context
	return micropay
}

// GetMicropay 拿到扫码信息请求PayJS
func (micropay *Micropay) GetMicropay(micropayRequest *MicropayRequest) (micropayResponse MicropayResponse, err error) {
	sign := util.Signature(micropayRequest, micropay.Context.Key)
	micropayRequest.Sign = sign
	response, err := util.PostJSON(getMicropayURL, micropayRequest)
	if err != nil {
		return
	}
	err = json.Unmarshal(response, &micropayResponse)
	if err != nil {
		return
	}
	if micropayResponse.ReturnCode == 0 {
		err = fmt.Errorf("GetMicropay Error , errcode=%d , errmsg=%s", micropayResponse.Status, micropayResponse.Msg)
		return
	}
	// 检测sign
	msgSignature := micropayResponse.Sign
	msgSignatureGen := util.Signature(micropayResponse, micropay.Context.Key)
	if msgSignature != msgSignatureGen {
		err = fmt.Errorf("消息不合法，验证签名失败")
		return
	}
	return
}
