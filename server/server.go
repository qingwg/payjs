package server

import (
	"fmt"
	"github.com/qingwg/payjs/context"
	"github.com/qingwg/payjs/util"
	"strconv"
)

//Server struct
type Server struct {
	*context.Context

	debug bool

	openID string

	messageHandler func(message Message)

	requestMsg  Message
	responseMsg interface{}
}

// Message PayJS支付成功异步通知过来的内容
type Message struct {
	ReturnCode    int    `json:"return_code"`    // 必填	1：支付成功
	TotalFee      int    `json:"total_fee"`      // 必填	金额。单位：分
	OutTradeNo    string `json:"out_trade_no"`   // 必填	用户端自主生成的订单号
	PayJSOrderID  string `json:"payjs_order_id"` // 必填	PAYJS 订单号
	TransactionID string `json:"transaction_id"` // 必填	微信用户手机显示订单号
	TimeEnd       string `json:"time_end"`       // 必填	支付成功时间
	Openid        string `json:"openid"`         // 必填	用户OPENID标示，本参数没有实际意义，旨在方便用户端区分不同用户
	Attach        string `json:"attach"`         // 非必填 用户自定义数据
	MchID         string `json:"mchid"`          // 必填	PAYJS 商户号
	Sign          string `json:"sign"`           // 必填	数据签名 详见签名算法
}

// NewServer init
func NewServer(context *context.Context) *Server {
	srv := new(Server)
	srv.Context = context
	return srv
}

// SetDebug set debug field
func (srv *Server) SetDebug(debug bool) {
	srv.debug = debug
}

//SetMessageHandler 设置用户自定义处理PayJS支付成功推送消息的方法
func (srv *Server) SetMessageHandler(handler func(message Message)) {
	srv.messageHandler = handler
}

//Serve 处理PayJS支付成功推送的消息
func (srv *Server) Serve() error {
	err := srv.handleRequest()
	if err != nil {
		return err
	}

	return srv.SendResponseMsg()
}

//HandleRequest 处理微信的请求。消息有可能反复推送，所以要去重；可能会增加数据，所以要取最新的；检测金额是否与自己订单相同；
func (srv *Server) handleRequest() (err error) {
	message, err := srv.getMessage()
	if err != nil {
		return
	}

	//set openID
	srv.openID = srv.Query("openid")

	srv.messageHandler(message)
	return
}

//getMessage 解析PayJS支付成功推送的消息
func (srv *Server) getMessage() (message Message, err error) {
	message.ReturnCode, _ = strconv.Atoi(srv.Request.PostFormValue("return_code"))
	message.TotalFee, _ = strconv.Atoi(srv.Request.PostFormValue("total_fee"))
	message.OutTradeNo = srv.Request.PostFormValue("out_trade_no")
	message.PayJSOrderID = srv.Request.PostFormValue("payjs_order_id")
	message.TransactionID = srv.Request.PostFormValue("transaction_id")
	message.TimeEnd = srv.Request.PostFormValue("time_end")
	message.Openid = srv.Request.PostFormValue("openid")
	message.Attach = srv.Request.PostFormValue("attach")
	message.MchID = srv.Request.PostFormValue("mchid")
	message.Sign = srv.Request.PostFormValue("sign")

	//if err = json.NewDecoder(srv.Request.Body).Decode(&message); err != nil {
	//	return message, fmt.Errorf("从body中解析json失败,err=%v", err)
	//}

	//验证消息签名
	msgSignature := message.Sign
	msgSignatureGen := util.Signature(message, srv.Context.Key)
	if msgSignature != msgSignatureGen {
		return message, fmt.Errorf("消息不合法，验证签名失败")
	}

	if message.ReturnCode != 1 {
		return message, fmt.Errorf("支付失败")
	}

	return message, err
}

func (srv *Server) SendResponseMsg() (err error) {
	srv.Writer.Header().Set("Content-Type", "application/json")
	srv.Writer.Write([]byte("ok"))
	return
}
