package payjs

import (
	"github.com/qingwg/payjs/cashier"
	"github.com/qingwg/payjs/context"
	"github.com/qingwg/payjs/facepay"
	"github.com/qingwg/payjs/js"
	"github.com/qingwg/payjs/mch"
	"github.com/qingwg/payjs/micropay"
	"github.com/qingwg/payjs/miniapp"
	"github.com/qingwg/payjs/native"
	"github.com/qingwg/payjs/order"
	"github.com/qingwg/payjs/server"
	"github.com/qingwg/payjs/user"
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
	NotifyUrl string
}

func New(cfg *Config) *PayJS {
	context := new(context.Context)
	copyConfigToContext(cfg, context)
	return &PayJS{context}
}

func copyConfigToContext(cfg *Config, context *context.Context) {
	context.MchID = cfg.MchID
	context.Key = cfg.Key
	context.NotifyUrl = cfg.NotifyUrl
}

// GetNative 扫码支付，主扫
func (payjs *PayJS) GetNative() *native.Native {
	return native.NewNative(payjs.Context)
}

// GetMicropay 扫码支付，主扫
func (payjs *PayJS) GetMicropay() *micropay.Micropay {
	return micropay.NewMicropay(payjs.Context)
}

// GetCashier 收银台支付 收银台方式同样是通过 JSAPI 方式发起的支付，只是简化了开发步骤和流程。适用于微信webview环境
func (payjs *PayJS) GetCashier() *cashier.Cashier {
	return cashier.NewCashier(payjs.Context)
}

// GetJs JSAPI 接口
func (payjs *PayJS) GetJs(req *http.Request, writer http.ResponseWriter) *js.Js {
	payjs.Context.Request = req
	payjs.Context.Writer = writer
	return js.NewJs(payjs.Context)
}

// GetMiniApp 微信小程序支付
func (payjs *PayJS) GetMiniApp() *miniapp.MiniApp {
	return miniapp.NewMiniApp(payjs.Context)
}

// GetFacepay 人脸支付
func (payjs *PayJS) GetFacepay() *facepay.Facepay {
	return facepay.NewFacepay(payjs.Context)
}

// GetOrder 订单 订单查询、订单关闭、订单退款
func (payjs *PayJS) GetOrder() *order.Order {
	return order.NewOrder(payjs.Context)
}

// GetServer 异步通知消息管理
func (payjs *PayJS) GetServer(req *http.Request, writer http.ResponseWriter) *server.Server {
	payjs.Context.Request = req
	payjs.Context.Writer = writer
	return server.NewServer(payjs.Context)
}

// GetUser 用户 用户详情
func (payjs *PayJS) GetUser() *user.User {
	return user.NewUser(payjs.Context)
}

// GetMch 商户 商户详情
func (payjs *PayJS) GetMch() *mch.Mch {
	return mch.NewMch(payjs.Context)
}
