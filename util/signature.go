package util

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

//Signature 签名
func Signature(params url.Values, privKey string) string {
	params.Del(`sign`)
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
	var src = strings.Join(pList, "&")
	src += "&key=" + privKey

	md5bs := md5.Sum([]byte(src))
	md5res := hex.EncodeToString(md5bs[:])
	return strings.ToUpper(md5res)
}

func (pj *PayJS) CreateTrade(param TradeParam) (res string, err error) {
	var p = url.Values{}
	jsonbs, _ := json.Marshal(param)
	jsonmap := make(map[string]interface{})
	json.Unmarshal(jsonbs, &jsonmap)
	for k, v := range jsonmap {
		p.Add(k, fmt.Sprintf("%v", v))
	}
	p.Add("mchid", pj.mchid)

	p.Add("sign", sign(p, pj.privKey))

	cli := http.Client{}
	r, err := cli.PostForm(pj.apiUrl, p)
	if err != nil {
		return ``, err
	}
	defer r.Body.Close()
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return ``, err
	}
	return string(bs), nil
}
