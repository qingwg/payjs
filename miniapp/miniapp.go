package miniapp

import (
	"github.com/qingwg/payjs/context"
	"github.com/qingwg/payjs/util"
)

// MiniApp struct
type MiniApp struct {
	*context.Context
}

// OrderInfo 后端按照下面参数构造订单参数，返回给前端
type OrderInfo struct {
	MchID      string `json:"mch_id"`       //Y 商户号
	TotalFee   int64  `json:"total_fee"`    //Y 金额。单位：分
	OutTradeNo string `json:"out_trade_no"` //Y 用户端自主生成的订单号
	Body       string `json:"body"`         //N 订单标题
	Attach     string `json:"attach"`       //N 用户自定义数据，在notify的时候会原样返回
	NotifyUrl  string `json:"notify_url"`   //N 异步通知地址
	Nonce      string `json:"nonce"`        //Y 随机字符串
	Sign       string `json:"sign"`         //Y 数据签名 详见签名算法
}

//NewMiniApp init
func NewMiniApp(context *context.Context) *MiniApp {
	miniApp := new(MiniApp)
	miniApp.Context = context
	return miniApp
}

// GetOrderInfo 获取小程序跳转所需的参数
func (miniApp *MiniApp) GetOrderInfo(totalFeeReq int64, bodyReq, outTradeNoReq, attachReq string) (orderInfo OrderInfo, err error) {
	orderInfo.MchID = miniApp.MchID
	orderInfo.TotalFee = totalFeeReq
	orderInfo.OutTradeNo = outTradeNoReq
	orderInfo.Body = bodyReq
	orderInfo.Attach = attachReq
	orderInfo.NotifyUrl = miniApp.NotifyUrl
	orderInfo.Nonce = util.RandomStr(32)
	orderInfo.Sign = util.Signature(orderInfo, miniApp.Key)
	return
}
