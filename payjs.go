package payjs

import (
	"github.com/yuyan2077/payjs/context"
	"github.com/yuyan2077/payjs/docs/wechat-develop/server"
	"github.com/yuyan2077/payjs/miniapp"
	"net/http"
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

// GetServer 异步通知消息管理
func (payjs *PayJS) GetServer(req *http.Request, writer http.ResponseWriter) *server.Server {
	payjs.Context.Request = req
	payjs.Context.Writer = writer
	return server.NewServer(payjs.Context)
}

// GetMiniAppPay 微信小程序支付配置
func (payjs *PayJS) GetMiniApp() *miniapp.MiniApp {
	return miniapp.NewMiniApp(payjs.Context)
}
