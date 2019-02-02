package mch

import (
	"encoding/json"
	"fmt"
	"github.com/qingwg/payjs/context"
	"github.com/qingwg/payjs/util"
)

const getMchInfoURL = "https://payjs.cn/api/info"

// Mch struct
type Mch struct {
	*context.Context
}

// MchInfoRequest 请求参数
type MchInfoRequest struct {
	MchID string `json:"mchid"` //Y	商户号
	Sign  string `json:"sign"`  //Y	数据签名 详见签名算法
}

// MchInfoResponse PayJS返回参数
type MchInfoResponse struct {
	ReturnCode int    `json:"return_code"` //Y	1:请求成功 0:请求失败
	ReturnMsg  string `json:"return_msg"`  //Y	返回消息
	Doudou     int64  `json:"doudou"`      //Y	用户豆豆数
	Name       string `json:"name"`        //Y	商户名称
	Username   string `json:"username"`    //Y	用户姓名
	IDcardNo   string `json:"idcardno"`    //Y	身份证号
	JsApiPath  string `json:"jsapi_path"`  //Y	JSAPI 支付目录
	Phone      string `json:"phone"`       //Y	客服电话
	MchID      string `json:"mchid"`       //Y	商户号
	Sign       string `json:"sign"`        //Y	数据签名 详见签名算法
}

//NewMch init
func NewMch(context *context.Context) *Mch {
	mch := new(Mch)
	mch.Context = context
	return mch
}

// GetMchInfo 请求PayJS获取支付二维码
func (mch *Mch) GetMchInfo() (mchInfoResponse MchInfoResponse, err error) {
	mchInfoRequest := MchInfoRequest{
		MchID: mch.MchID,
	}
	sign := util.Signature(mchInfoRequest, mch.Key)
	mchInfoRequest.Sign = sign
	response, err := util.PostJSON(getMchInfoURL, mchInfoRequest)
	if err != nil {
		return
	}

	err = json.Unmarshal(response, &mchInfoResponse)
	if err != nil {
		return
	}
	if mchInfoResponse.ReturnCode != 1 {
		err = fmt.Errorf("GetMchInfo Error , errcode=%d , errmsg=%s", mchInfoResponse.ReturnCode, mchInfoResponse.ReturnMsg)
		return
	}
	// 检测sign
	msgSignature := mchInfoResponse.Sign
	msgSignatureGen := util.Signature(mchInfoResponse, mch.Key)
	if msgSignature != msgSignatureGen {
		err = fmt.Errorf("消息不合法，验证签名失败")
		return
	}
	return
}
