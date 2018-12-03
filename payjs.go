package payjs

import (
	"github.com/yuyan2077/payjs/context"
	"github.com/yuyan2077/payjs/miniapp"
)

// PayJS struct
type PayJS struct {
	Context *context.Context
}

// config for PayJS
type Config struct {
	Key       string
	MchID     string
	NotifyURL string
}

func New(cfg *Config) *PayJS {
	context := new(context.Context)
	copyConfigToContext(cfg, context)
	return &PayJS{context}
}

func copyConfigToContext(cfg *Config, context *context.Context) {
	context.MchID = cfg.MchID
	context.Key = cfg.Key
	context.NotifyURL = cfg.NotifyURL
}

// GetMiniAppPay 微信小程序支付配置
func (payjs *PayJS) GetMiniApp() *miniapp.MiniApp {
	return miniapp.NewMiniApp(payjs.Context)
}
