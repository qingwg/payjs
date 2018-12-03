package miniapp

import "github.com/yuyan2077/payjs/context"

// MiniApp struct
type MiniApp struct {
	*context.Context
}

// OrderInfo 后端按照下面参数构造订单参数，返回给前端
type OrderInfo struct {
	MchID      string `json:"mch_id"`       // 必填 商户号
	TotalFee   int    `json:"total_fee"`    // 必填 金额。单位：分
	OutTradeNo string `json:"out_trade_no"` // 必填 用户端自主生成的订单号
	Body       string `json:"body"`         // 非必填 订单标题
	Attach     string `json:"attach"`       // 非必填 用户自定义数据，在notify的时候会原样返回
	NotifyUrl  string `json:"notify_url"`   // 非必填 异步通知地址
	Nonce      string `json:"nonce"`        // 必填 随机字符串
	Sign       string `json:"sign"`         // 必填 数据签名 详见签名算法
}

//NewMiniApp init
func NewMiniApp(context *context.Context) *MiniApp {
	miniApp := new(MiniApp)
	miniApp.Context = context
	return miniApp
}
