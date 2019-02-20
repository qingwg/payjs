package user

import (
	"encoding/json"
	"fmt"
	"github.com/qingwg/payjs/context"
	"github.com/qingwg/payjs/util"
	"net/http"
)

const getUserInfoURL = "https://payjs.cn/api/user"
const getUserOpenIDURL = "https://payjs.cn/api/openid"

// User struct
type User struct {
	*context.Context
}

// UserInfoRequest 请求参数
type UserInfoRequest struct {
	MchID  string `json:"mchid"`  //Y	商户号
	Openid string `json:"openid"` //Y	openid
	Sign   string `json:"sign"`   //Y	数据签名 详见签名算法
}

// UserInfoResponse PayJS返回参数
type UserInfoResponse struct {
	ReturnCode int      `json:"return_code"` //Y	1:请求成功 0:请求失败
	ReturnMsg  string   `json:"return_msg"`  //Y	返回消息
	User       UserInfo `json:"user"`        //N	用户资料
	Sign       string   `json:"sign"`        //Y	数据签名 详见签名算法
}

// UserInfo 用户参数说明(同微信官方文档)
type UserInfo struct {
	Subscribe      int    `json:"subscribe"`       //	用户是否订阅 PAYJS 公众号标识，值为0时，代表此用户没有关注该公众号，拉取不到其余信息，只有openid和UnionID
	Openid         string `json:"openid"`          //	用户的标识，对公众号唯一
	Nickname       string `json:"nickname"`        //   用户的昵称
	Sex            int    `json:"sex"`             //   用户的性别，值为1时是男性，值为2时是女性，值为0时是未知
	City           string `json:"city"`            //	用户所在城市
	Country        string `json:"country"`         //   用户所在国家
	Province       string `json:"province"`        //   用户所在省份
	Language       string `json:"language"`        //   用户的语言，简体中文为zh_CN
	Headimgurl     string `json:"headimgurl"`      //   用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），用户没有头像时该项为空。若用户更换头像，原有头像URL将失效。
	SubscribeTime  int    `json:"subscribe_time"`  //   用户关注时间，为时间戳。如果用户曾多次关注，则取最后关注时间
	Remark         string `json:"remark"`          //   公众号运营者对粉丝的备注，公众号运营者可在微信公众平台用户管理界面对粉丝添加备注
	Groupid        int    `json:"groupid"`         //   用户所在的分组ID
	TagidList      []int  `json:"tagid_list"`      //   用户被打上的标签ID列表
	SubscribeScene string `json:"subscribe_scene"` //   返回用户关注的渠道来源，ADD_SCENE_SEARCH 公众号搜索，ADD_SCENE_ACCOUNT_MIGRATION 公众号迁移，ADD_SCENE_PROFILE_CARD 名片分享，ADD_SCENE_QR_CODE 扫描二维码，ADD_SCENE_PROFILE_LINK 图文页内名称点击，ADD_SCENE_PROFILE_ITEM 图文页右上角菜单，ADD_SCENE_PAID 支付后关注，ADD_SCENE_OTHERS 其他
	QrScene        int    `json:"qr_scene"`        //   二维码扫码场景
	QrSceneStr     string `json:"qr_scene_str"`    //   二维码扫码场景描述
}

//NewUser init
func NewUser(context *context.Context) *User {
	user := new(User)
	user.Context = context
	return user
}

// GetUserOpenIDUrl 获取请求url
func (user *User) GetUserOpenIDUrl(callbackUrlReq string) (src string, err error) {
	return getUserOpenIDURL + "?mchid=" + user.MchID + "&callback_url=" + callbackUrlReq, nil
}

// GetUserOpenID 获取用户 OPENID
func (user *User) GetUserOpenID(req *http.Request) (openid string, err error) {
	//set openID
	if values, ok := req.URL.Query()["openid"]; ok && len(values) > 0 {
		openid = values[0]
	}
	return
}

// GetUserInfo 根据支付订单中的 openid 获取用户更多资料，例如昵称、头像、性别、地区等信息
//提示：openid必须是有在 PAYJS 支付过的
//特别说明：此请求同时返回了用户在微信公众号中的 unionid ，进一步阅读详情
func (user *User) GetUserInfo(openid string) (userInfoResponse UserInfoResponse, err error) {
	userInfoRequest := UserInfoRequest{
		MchID:  user.MchID,
		Openid: openid,
	}
	sign := util.Signature(userInfoRequest, user.Key)
	userInfoRequest.Sign = sign
	response, err := util.PostJSON(getUserInfoURL, userInfoRequest)
	if err != nil {
		return
	}

	err = json.Unmarshal(response, &userInfoResponse)
	if err != nil {
		return
	}

	if userInfoResponse.ReturnCode == 0 {
		err = fmt.Errorf("GetUserInfo Error , errcode=%d , errmsg=%s", userInfoResponse.ReturnCode, userInfoResponse.ReturnMsg)
		return
	}

	//todo:解决多维结构造成签名验证失败bug
	//// 检测sign
	//msgSignature := userInfoResponse.Sign
	//msgSignatureGen := util.Signature(userInfoResponse, user.Key)
	//if msgSignature != msgSignatureGen {
	//	err = fmt.Errorf("消息不合法，验证签名失败")
	//	return
	//}

	return
}
