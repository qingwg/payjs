package notify

import (
	"fmt"
	"github.com/qingwg/payjs/context"
	"net/http"
	"strconv"
)

//Notify struct
type Notify struct {
	*context.Context

	debug bool

	ServerWriter  http.ResponseWriter
	ServerRequest *http.Request

	messageHandler func(message Message)

	requestMsg  Message
	responseMsg interface{}
}

// Message PayJS支付成功异步通知过来的内容
type Message struct {
	ReturnCode    int    `json:"return_code"`    // 必填	1：支付成功
	TotalFee      int64  `json:"total_fee"`      // 必填	金额。单位：分
	OutTradeNo    string `json:"out_trade_no"`   // 必填	用户端自主生成的订单号
	PayJSOrderID  string `json:"payjs_order_id"` // 必填	PAYJS 订单号
	TransactionID string `json:"transaction_id"` // 必填	微信用户手机显示订单号
	TimeEnd       string `json:"time_end"`       // 必填	支付成功时间
	Openid        string `json:"openid"`         // 必填	用户OPENID标示，本参数没有实际意义，旨在方便用户端区分不同用户
	Attach        string `json:"attach"`         // 非必填 用户自定义数据
	MchID         string `json:"mchid"`          // 必填	PAYJS 商户号
	Sign          string `json:"sign"`           // 必填	数据签名 详见签名算法
}

// Query returns the keyed url query value if it exists
func (notify *Notify) Query(key string) string {
	value, _ := notify.GetQuery(key)
	return value
}

// GetQuery is like Query(), it returns the keyed url query value
func (notify *Notify) GetQuery(key string) (string, bool) {
	req := notify.ServerRequest
	if values, ok := req.URL.Query()[key]; ok && len(values) > 0 {
		return values[0], true
	}
	return "", false
}

// NewNotify init
func NewNotify(context *context.Context, req *http.Request, writer http.ResponseWriter) *Notify {
	notify := new(Notify)
	notify.Context = context
	notify.ServerWriter = writer
	notify.ServerRequest = req
	return notify
}

// SetDebug set debug field
func (notify *Notify) SetDebug(debug bool) {
	notify.debug = debug
}

//SetMessageHandler 设置用户自定义处理PayJS支付成功推送消息的方法
func (notify *Notify) SetMessageHandler(handler func(message Message)) {
	notify.messageHandler = handler
}

//Serve 处理PayJS支付成功推送的消息
func (notify *Notify) Serve() error {
	err := notify.handleRequest()
	if err != nil {
		return err
	}

	return notify.SendResponseMsg()
}

//HandleRequest 处理微信的请求。消息有可能反复推送，所以要去重；可能会增加数据，所以要取最新的；检测金额是否与自己订单相同；
func (notify *Notify) handleRequest() (err error) {
	message, err := notify.getMessage()
	if err != nil {
		return
	}

	notify.messageHandler(message)
	return
}

//getMessage 解析PayJS支付成功推送的消息
func (notify *Notify) getMessage() (message Message, err error) {
	message.ReturnCode, _ = strconv.Atoi(notify.ServerRequest.PostFormValue("return_code"))
	message.TotalFee, _ = strconv.ParseInt(notify.ServerRequest.PostFormValue("total_fee"), 10, 64)
	message.OutTradeNo = notify.ServerRequest.PostFormValue("out_trade_no")
	message.PayJSOrderID = notify.ServerRequest.PostFormValue("payjs_order_id")
	message.TransactionID = notify.ServerRequest.PostFormValue("transaction_id")
	message.TimeEnd = notify.ServerRequest.PostFormValue("time_end")
	message.Openid = notify.ServerRequest.PostFormValue("openid")
	message.Attach = notify.ServerRequest.PostFormValue("attach")
	message.MchID = notify.ServerRequest.PostFormValue("mchid")
	message.Sign = notify.ServerRequest.PostFormValue("sign")

	//if err = json.NewDecoder(srv.Request.Body).Decode(&message); err != nil {
	//	return message, fmt.Errorf("从body中解析json失败,err=%v", err)
	//}

	//todo:解决多维结构造成签名验证失败bug
	////验证消息签名
	//msgSignature := message.Sign
	//msgSignatureGen := util.Signature(message, notify.Key)
	//if msgSignature != msgSignatureGen {
	//	return message, fmt.Errorf("消息不合法，验证签名失败")
	//}

	if message.ReturnCode != 1 {
		return message, fmt.Errorf("支付失败")
	}

	return message, err
}

func (notify *Notify) SendResponseMsg() (err error) {
	notify.ServerWriter.Header().Set("Content-Type", "application/json")
	notify.ServerWriter.Write([]byte("success"))
	return
}
