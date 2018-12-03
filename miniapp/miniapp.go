package miniapp

import "github.com/yuyan2077/payjs/context"

// MiniApp struct
type MiniApp struct {
	*context.Context
}

// Config 返回给用户jssdk配置信息
type Config struct {
	AppID     string `json:"app_id"`
	Timestamp int64  `json:"timestamp"`
	NonceStr  string `json:"nonce_str"`
	Signature string `json:"signature"`
}

//NewMiniApp init
func NewMiniApp(context *context.Context) *MiniApp {
	miniApp := new(MiniApp)
	miniApp.Context = context
	return miniApp
}
