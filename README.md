# PayJS SDK for Go


使用Golang开发的PayJS SDK，简单、易用。

[PayJS](https://payjs.cn/ref/DWPXBZ)是微信支付个人接口解决方案，感兴趣的可以去官网看下。

这里是SDK的演示地址：[https://payjs.qingwuguo.com](https://payjs.qingwuguo.com)

## TODO
- examples未完成
- 商户资料签名验证失败BUG
- JSAPI支付签名验证失败BUG
- 付款码支付测试及演示
- JSAPI支付测试及演示
- 小程序支付测试及演示
- 人脸支付测试及演示
- 订单-撤销测试及演示
- 异步通知演示
- 商户资料测试及演示
- 银行编码查询测试及演示

#### 和主流框架配合使用

主要是request和responseWriter在不同框架中获取方式可能不一样：

- Gin Framework: [./examples/gin](./examples/gin)

## 基本配置及初始化
下面的是伪代码，请自行理解
```go
payjsConfig := &payjs.Config{
    Key:       "PayJS的通信密钥",
    MchID:     "PayJS的商户号",
    NotifyUrl: "异步通知的路由",
}
Pay = payjs.New(payjsConfig)
```


## 基本API使用

- [扫码支付](#扫码支付)
- [付款码支付](#付款码支付)
- [收银台支付](#收银台支付)
- [JSAPI支付](#JSAPI支付)
- [小程序支付](#小程序支付)
- [人脸支付](#人脸支付)
- [订单](#订单)
	- 查询
	- 关闭
	- 撤销
	- 退款
- [异步通知](#异步通知)
- [用户](#用户)
    - 获取浏览器跳转的url
	- 获取openid
	- ~~获取用户资料~~（即将废弃，SDK相应代码已删除）
- [商户资料](#商户)
- [银行编码查询](#银行编码查询)

## 扫码支付
下面的是伪代码，请自行理解
```go
PayNative := Pay.GetNative()
type Request struct {
	TotalFee     int    `json:"total_fee"`       //Y	金额。单位：分
	Body         string `json:"body"`            //N	订单标题
	Attach       string `json:"attach"`          //N	用户自定义数据，在notify的时候会原样返回
	OutTradeNo   string `json:"out_trade_no"`    //Y	用户端自主生成的订单号
}
PayNative.Create(Request.TotalFee, Request.Body, Request.OutTradeNo, Request.Attach)
```

官方文档：[扫码支付
](https://help.payjs.cn/api-lie-biao/sao-ma-zhi-fu.html)

## 付款码支付（未测试）
下面的是伪代码，请自行理解
```go
PayMicropay := Pay.GetMicropay()
type Request struct {
	TotalFee     int    `json:"total_fee"`       //Y	金额。单位：分
	Body         string `json:"body"`            //N	订单标题
	Attach       string `json:"attach"`          //N	用户自定义数据，在notify的时候会原样返回
	OutTradeNo   string `json:"out_trade_no"`    //Y	用户端自主生成的订单号
	AuthCode     string `json:"auth_code"`       //Y	扫码支付授权码，设备读取用户微信中的条码或者二维码信息(注：用户刷卡条形码规则：18位纯数字，以10、11、12、13、14、15开头)
}
PayMicropay.Create(Request.TotalFee, Request.Body, Request.OutTradeNo, Request.Attach, Request.AuthCode)
```

官方文档：[付款码支付
](https://help.payjs.cn/api-lie-biao/shua-qia-zhi-fu.html)

## 收银台支付
下面的是伪代码，请自行理解
```go
PayCashier := Pay.GetCashier()
type Request struct {
	TotalFee     int    `json:"total_fee"`       //Y	金额。单位：分
	Body         string `json:"body"`            //N	订单标题
	Attach       string `json:"attach"`          //N	用户自定义数据，在notify的时候会原样返回
	OutTradeNo   string `json:"out_trade_no"`    //Y	用户端自主生成的订单号
	CallbackUrl  string `json:"callback_url"`    //N	用户支付成功后，前端跳转地址。留空则支付后关闭webview
	Auto         bool   `json:"auto"`            //N	auto=1：无需点击支付按钮，自动发起支付。默认手动点击发起支付
	Hide         bool   `json:"hide"`            //N	hide=1：隐藏收银台背景界面。默认显示背景界面（这里hide为1时，自动忽略auto参数）
}
PayCashier.GetRequestUrl(Request.TotalFee, Request.Body, Request.OutTradeNo, Request.Attach, Request.CallbackUrl, Request.Auto, Request.Hide)
```

官方文档：[收银台支付
](https://help.payjs.cn/api-lie-biao/shou-yin-tai-zhi-fu.html)

## JSAPI支付（未测试）（有bug）
下面的是伪代码，请自行理解
```go
PayJS := Pay.GetJs()
type Request struct {
	TotalFee   int    `json:"total_fee"`    //Y	金额。单位：分
    OutTradeNo string `json:"out_trade_no"` //Y	用户端自主生成的订单号，在用户端要保证唯一性
    Body       string `json:"body"`         //N	订单标题
    Attach     string `json:"attach"`       //N	用户自定义数据，在notify的时候会原样返回
    Openid     string `json:"openid"`       //Y	用户openid
}
PayJS.Create(Request.TotalFee, Request.Body, Request.OutTradeNo, Request.Attach, Request.Openid)
```

官方文档：[JSAPI支付
](https://help.payjs.cn/api-lie-biao/jsapiyuan-sheng-zhi-fu.html)

## 小程序支付（未测试）
下面的是伪代码，请自行理解

小程序发起支付的解决方案有两种，仅供测试使用

- 方案一：使用小程序消息，结合收银台模式，可以解决小程序支付
- 方案二：使用小程序跳转到 PAYJS 小程序，支付后返回（下面代码是方案二）
```go
PayMiniApp := Pay.GetMiniApp()
type Request struct {
	TotalFee   int    `json:"total_fee"`    //Y	金额。单位：分
    OutTradeNo string `json:"out_trade_no"` //Y	用户端自主生成的订单号，在用户端要保证唯一性
    Body       string `json:"body"`         //N	订单标题
    Attach     string `json:"attach"`       //N	用户自定义数据，在notify的时候会原样返回
    Nonce      string `json:"nonce"`        //Y 随机字符串
}
// 获取小程序跳转所需的参数
PayMiniApp.GetOrderInfo(Request.TotalFee, Request.Body, Request.OutTradeNo, Request.Attach, Request.Nonce)
```

官方文档：[小程序支付
](https://help.payjs.cn/api-lie-biao/xiao-cheng-xu-zhi-fu.html)

## 人脸支付（未测试）
下面的是伪代码，请自行理解
```go
PayFacepay := Pay.GetFacepay()
type Request struct {
	TotalFee   int    `json:"total_fee"`    //Y	金额。单位：分
    OutTradeNo string `json:"out_trade_no"` //Y	用户端自主生成的订单号，在用户端要保证唯一性
    Body       string `json:"body"`         //N	订单标题
    Attach     string `json:"attach"`       //N	用户自定义数据，在notify的时候会原样返回
    Openid     string `json:"openid"`       //Y	OPENID
    FaceCode   string `json:"face_code"`    //Y	人脸支付识别码
}
PayFacepay.Create(Request.TotalFee, Request.Body, Request.OutTradeNo, Request.Attach, Request.Openid, Request.FaceCode)
```

官方文档：[人脸支付
](https://help.payjs.cn/api-lie-biao/ren-lian-zhi-fu.html)

## 订单
下面的是伪代码，请自行理解
```go
// 初始化
PayOrder := Pay.GetOrder()
```

####查询
```go
type Request struct {
	PayJSOrderID string `json:"payjs_order_id"` //Y	PAYJS 平台订单号
}
PayOrder.Check(Request.PayJSOrderID)
```
官方文档：[订单-查询
](https://help.payjs.cn/api-lie-biao/ding-dan-cha-xun.html)

####关闭
```go
type Request struct {
	PayJSOrderID string `json:"payjs_order_id"` //Y	PAYJS 平台订单号
}
PayOrder.Close(Request.PayJSOrderID)
```
官方文档：[订单-关闭
](https://help.payjs.cn/api-lie-biao/guan-bi-ding-dan.html)

####撤销（未测试）
```go
type Request struct {
	PayJSOrderID string `json:"payjs_order_id"` //Y	PAYJS 平台订单号
}
PayOrder.Reverse(Request.PayJSOrderID)
```
官方文档：[订单-撤销
](https://help.payjs.cn/api-lie-biao/che-xiao-ding-dan.html)

####退款
```go
type Request struct {
	PayJSOrderID string `json:"payjs_order_id"` //Y	PAYJS 平台订单号
}
PayOrder.Refund(Request.PayJSOrderID)
```
官方文档：[订单-退款
](https://help.payjs.cn/api-lie-biao/tui-kuan.html)

## 异步通知
下面的是伪代码，请自行理解
```go
// 传入request和responseWriter
PayNotify := Pay.GetNotify(request, responseWriter)

//设置接收消息的处理方法
PayNotify.SetMessageHandler(func(msg notify.Message) {
    //这里处理支付成功回调，一般是修改数据库订单信息等等
    //msg即为支付成功异步通知过来的内容
})

//处理消息接收以及回复
err := PayNotify.Serve()
if err != nil {
    fmt.Println(err)
    return
}

//发送回复的消息
PayNotify.SendResponseMsg()
```

官方文档：[异步通知
](https://help.payjs.cn/api-lie-biao/jiao-yi-xin-xi-tui-song.html)

## 用户
下面的是伪代码，请自行理解
```go
// 初始化
PayUser := Pay.GetUser()
```

####获取浏览器跳转的url
```go
type Request struct {
	CallbackUrl string `json:"callback_url"` //Y	接收 openid 的 url。必须为可直接访问的url，不能带session验证、csrf验证。url 可携带最多1个参数，多个参数会自动忽略
}
PayUser.GetUserOpenIDUrl(Request.CallbackUrl)
```
官方文档：[用户-获取浏览器跳转的url
](https://help.payjs.cn/api-lie-biao/huo-qu-openid.html)

####获取openid
```go
// 在callback_url方法内，传入request
PayUser.GetUserOpenID(request)
```
官方文档：[用户-获取openid
](https://help.payjs.cn/api-lie-biao/huo-qu-openid.html)

## 商户资料（未测试）（有bug）
下面的是伪代码，请自行理解
```go
PayMch := Pay.GetMch()
PayMch.GetMchInfo()
```

官方文档：[商户资料
](https://help.payjs.cn/api-lie-biao/shang-hu-zi-liao.html)

## 银行编码查询（未测试）
下面的是伪代码，请自行理解
```go
PayBank := Pay.GetBank()
type Request struct {
	Bank string `json:"bank"` //Y	银行简写
}
PayBank.GetBankInfo(Request.Bank)
```

官方文档：[银行编码查询
](https://help.payjs.cn/api-lie-biao/yin-xing-bian-ma-cha-xun.html)

## License

[MIT](LICENSE)
