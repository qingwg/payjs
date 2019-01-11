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
	fmt.Println("=====message", message)
	var params = url.Values{"a": []string{}}
	jsonbs, _ := json.Marshal(message)
	jsonmap := make(map[string]interface{})
	json.Unmarshal(jsonbs, &jsonmap)
	for k, v := range jsonmap {
		switch t := v.(type) {
		default:
			params.Add(k, fmt.Sprintf("%v", v))
		case map[string]interface{}:
			for kk, vv := range t {
				params.Add(k+"["+kk+"]", fmt.Sprintf("%v", vv))
			}
		}
	}

	params.Del(`sign`)

	var keys = make([]string, 0, 0)
	for key := range params {
		if params.Get(key) != `` {
			keys = append(keys, key)
		}
		//keys = append(keys, key)
	}
	sort.Strings(keys)
	fmt.Println("=====keys", keys)

	var pList = make([]string, 0, 0)
	for _, key := range keys {
		var value = strings.TrimSpace(params.Get(key))
		if len(value) > 0 {
			pList = append(pList, key+"="+value)
		}
	}
	var src = strings.Join(pList, "&")
	src += "&key=" + privKey
	fmt.Println("=====src", src)

	sign = MD5Sum(src)

	return sign
}
