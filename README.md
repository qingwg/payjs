# PayJS SDK for Go


使用Golang开发的PayJS SDK，简单、易用。

[PayJS](https://payjs.cn/ref/DWPXBZ)是微信支付个人接口解决方案，感兴趣的可以去官网看下。

这里是SDK的演示地址：[https://payjs.qingwuguo.com](https://payjs.qingwuguo.com)

## TODO
- 商户资料接口签名验证失败BUG
- 付款码支付接口测试及演示
- JSAPI支付接口测试及演示

#### 和主流框架配合使用

主要是request和responseWriter在不同框架中获取方式可能不一样：

- Gin Framework: [./examples/gin/gin.go](./examples/gin/gin.go)

## 基本配置及初始化

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
    - 获取openid跳转url
	- 获取openid
	- ~~获取用户资料~~（即将废弃，SDK相应代码已删除）
- [商户资料](#商户)
- [银行编码查询](#银行编码查询)

## 扫码支付

```go
PayNative := Pay.GetNative()
//totalFee:金额。单位：分
//body:订单标题
//outTradeNo:用户端自主生成的订单号
//attach:用户自定义数据，在notify的时候会原样返回
PayNative.Create(totalFee, body, outTradeNo, attach)
```

官方文档：[扫码支付
](https://help.payjs.cn/api-lie-biao/sao-ma-zhi-fu.html)

## License

[MIT](LICENSE)
