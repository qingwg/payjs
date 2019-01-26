package ip

import (
	"encoding/json"
	"fmt"
	"github.com/qingwg/payjs/context"
	"github.com/qingwg/payjs/util"
)

const getIPListURL = "https://payjs.cn/api/iplist"

// IP struct
type IP struct {
	*context.Context
}

// IPListRequest 请求参数
type IPListRequest struct {
	MchID string `json:"mchid"` //Y	商户号
	Sign  string `json:"sign"`  //Y	数据签名 详见签名算法
}

// IPListResponse PayJS返回参数
type IPListResponse struct {
	ReturnCode int      `json:"return_code"` //Y	1:请求成功 0:请求失败
	ReturnMsg  string   `json:"return_msg"`  //Y	返回消息
	IPList     []string `json:"iplist"`      //Y	ip地址列表
	Sign       string   `json:"sign"`        //Y	数据签名 详见签名算法
}

//NewIP init
func NewIP(context *context.Context) *IP {
	ip := new(IP)
	ip.Context = context
	return ip
}

// GetIPList 根据商户号查询异步通知服务器的IP列表
func (ip *IP) GetIPList() (ipListResponse IPListResponse, err error) {
	ipListRequest := IPListRequest{
		MchID: ip.MchID,
	}
	sign := util.Signature(ipListRequest, ip.Key)
	ipListRequest.Sign = sign
	response, err := util.PostJSON(getIPListURL, ipListRequest)
	if err != nil {
		return
	}

	err = json.Unmarshal(response, &ipListResponse)
	if err != nil {
		return
	}
	if ipListResponse.ReturnCode != 1 {
		err = fmt.Errorf("GetIPList Error , errcode=%d , errmsg=%s", ipListResponse.ReturnCode, ipListResponse.ReturnMsg)
		return
	}
	//// 检测sign
	//msgSignature := ipListResponse.Sign
	//msgSignatureGen := util.Signature(ipListResponse, mch.Key)
	//if msgSignature != msgSignatureGen {
	//	err = fmt.Errorf("消息不合法，验证签名失败")
	//	return
	//}
	return
}
