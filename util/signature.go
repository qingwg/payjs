package util

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

//Signature 签名
func Signature(message interface{}, privKey string) (sign string) {
	var params = url.Values{}
	jsonbs, _ := json.Marshal(message)
	jsonmap := make(map[string]interface{})
	json.Unmarshal(jsonbs, &jsonmap)
	for k, v := range jsonmap {
		params.Add(k, fmt.Sprintf("%v", v))
	}

	params.Del(`sign`)
	params.Del(`status`)
	params.Del(`msg`)
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

	sign = MD5Sum(src)
	return sign
}
