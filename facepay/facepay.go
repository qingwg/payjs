package facepay

import (
	"encoding/json"
	"fmt"
	"github.com/qingwg/payjs/context"
	"github.com/qingwg/payjs/util"
)

const getCreateURL = "https://payjs.cn/api/facepay"

// Facepay struct
type Facepay struct {
	*context.Context
}

// CreateRequest 请求参数
type CreateRequest struct {
	MchID      string `json:"mchid"`        //Y	商户号
	TotalFee   int64  `json:"total_fee"`    //Y	金额。单位：分
	OutTradeNo string `json:"out_trade_no"` //Y	用户端自主生成的订单号
	Body       string `json:"body"`         //N	订单标题
	Attach     string `json:"attach"`       //N	用户自定义数据，在notify的时候会原样返回
	Openid     string `json:"openid"`       //Y	OPENID
	FaceCode   string `json:"face_code"`    //Y	人脸支付识别码
	Sign       string `json:"sign"`         //Y	数据签名 详见签名算法
}

// CreateResponse PayJS返回参数
type CreateResponse struct {
	ReturnCode   int    `json:"return_code"`    //Y	1:请求成功，0:请求失败
	Msg          string `json:"msg"`            //N	return_code为0时返回的错误消息
	ReturnMsg    string `json:"return_msg"`     //Y	返回消息
	PayJSOrderID string `json:"payjs_order_id"` //Y	PAYJS 平台订单号
	OutTradeNo   string `json:"out_trade_no"`   //Y	用户生成的订单号原样返回
	TotalFee     string `json:"total_fee"`      //Y	金额。单位：分
	Sign         string `json:"sign"`           //Y	数据签名 详见签名算法
}

//NewFacepay init
func NewFacepay(context *context.Context) *Facepay {
	facepay := new(Facepay)
	facepay.Context = context
	return facepay
}

// Create
func (facepay *Facepay) Create(totalFeeReq int64, bodyReq, outTradeNoReq, attachReq, openidReq, faceCode string) (createResponse CreateResponse, err error) {
	facepayRequest := CreateRequest{
		MchID:      facepay.MchID,
		TotalFee:   totalFeeReq,
		OutTradeNo: outTradeNoReq,
		Body:       bodyReq,
		Attach:     attachReq,
		Openid:     openidReq,
		FaceCode:   faceCode,
	}
	sign := util.Signature(facepayRequest, facepay.Key)
	facepayRequest.Sign = sign
	response, err := util.PostJSON(getCreateURL, facepayRequest)
	if err != nil {
		return
	}
	err = json.Unmarshal(response, &createResponse)
	if err != nil {
		return
	}
	if createResponse.ReturnCode == 0 {
		err = fmt.Errorf("FacepayCreate Error , errcode=%d , errmsg=%s, errmsg=%s", createResponse.ReturnCode, createResponse.Msg, createResponse.ReturnMsg)
		return
	}
	// 检测sign
	msgSignature := createResponse.Sign
	msgSignatureGen := util.Signature(createResponse, facepay.Key)
	if msgSignature != msgSignatureGen {
		err = fmt.Errorf("消息不合法，验证签名失败")
		return
	}
	return
}
