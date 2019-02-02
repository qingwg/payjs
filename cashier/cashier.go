package cashier

import (
	"encoding/json"
	"fmt"
	"github.com/qingwg/payjs/context"
	"github.com/qingwg/payjs/util"
	"net/url"
	"sort"
	"strings"
)

const getCashierURL = "https://payjs.cn/api/cashier"

// Cashier struct
type Cashier struct {
	*context.Context
}

// CashierRequest 请求参数
type CashierRequest struct {
	MchID       string `json:"mchid"`        //Y	商户号
	TotalFee    int64  `json:"total_fee"`    //Y	金额。单位：分
	OutTradeNo  string `json:"out_trade_no"` //Y	用户端自主生成的订单号，在用户端要保证唯一性
	Body        string `json:"body"`         //N	订单标题
	Attach      string `json:"attach"`       //N	用户自定义数据，在notify的时候会原样返回
	NotifyUrl   string `json:"notify_url"`   //N	接收微信支付异步通知的回调地址。必须为可直接访问的URL，不能带参数、session验证、csrf验证。留空则不通知
	CallbackUrl string `json:"callback_url"` //N	用户支付成功后，前端跳转地址。留空则支付后关闭webview
	Auto        int    `json:"auto"`         //N	auto=1：无需点击支付按钮，自动发起支付。默认手动点击发起支付
	Hide        int    `json:"hide"`         //N	hide=1：隐藏收银台背景界面。默认显示背景界面（这里hide为1时，自动忽略auto参数）
	Sign        string `json:"sign"`         //Y	数据签名 详见签名算法
}

// CashierResponse PayJS返回参数
type CashierResponse struct {
	ReturnCode int    `json:"return_code"` //Y	0:提交失败
	Status     int    `json:"status"`      //Y	0:失败
	Msg        string `json:"msg"`         //Y	失败原因
	ReturnMsg  string `json:"return_msg"`  //Y	失败原因，同msg
}

//NewCashier init
func NewCashier(context *context.Context) *Cashier {
	cashier := new(Cashier)
	cashier.Context = context
	return cashier
}

func (cashier *Cashier) GetRequestUrl(totalFeeReq int64, bodyReq, outTradeNoReq, attachReq, callbackUrlReq string, auto, hide int) (src string, err error) {
	cashierRequest := CashierRequest{
		MchID:       cashier.MchID,
		TotalFee:    totalFeeReq,
		OutTradeNo:  outTradeNoReq,
		Body:        bodyReq,
		Attach:      attachReq,
		NotifyUrl:   cashier.NotifyUrl,
		CallbackUrl: callbackUrlReq,
		Auto:        auto,
		Hide:        hide,
	}
	sign := util.Signature(cashierRequest, cashier.Key)
	cashierRequest.Sign = sign

	var params = url.Values{}
	jsonbs, _ := json.Marshal(cashierRequest)
	jsonmap := make(map[string]interface{})
	json.Unmarshal(jsonbs, &jsonmap)
	for k, v := range jsonmap {
		params.Add(k, fmt.Sprintf("%v", v))
	}

	var keys = make([]string, 0, 0)
	for key := range params {
		if params.Get(key) != `` {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)

	var pList = make([]string, 0, 0)
	for _, key := range keys {
		var value = strings.TrimSpace(params.Get(key))
		if len(value) > 0 {
			pList = append(pList, key+"="+value)
		}
	}

	src = getCashierURL + "?"
	src += strings.Join(pList, "&")
	return
}
