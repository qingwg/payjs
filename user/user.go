package user

import (
	"github.com/qingwg/payjs/context"
	"net/http"
)

const getUserOpenIDURL = "https://payjs.cn/api/openid"

// User struct
type User struct {
	*context.Context
}

//NewUser init
func NewUser(context *context.Context) *User {
	user := new(User)
	user.Context = context
	return user
}

// GetUserOpenIDUrl 获取请求url
func (user *User) GetUserOpenIDUrl(callbackUrlReq string) (src string, err error) {
	return getUserOpenIDURL + "?callback_url=" + callbackUrlReq, nil
}

// GetUserOpenID 获取用户 OPENID
func (user *User) GetUserOpenID(req *http.Request) (openid string, err error) {
	//set openID
	if values, ok := req.URL.Query()["openid"]; ok && len(values) > 0 {
		openid = values[0]
	}
	return
}
