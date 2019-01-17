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
		//switch t := v.(type) {
		//default:
		//	params.Add(k, fmt.Sprintf("%v", v))
		//case map[string]interface{}:
		//	params.Add(k, fmt.Sprintf("%v", v))
		//}
	}
	params.Del(`sign`)
	//fmt.Println("=====params", params)

	var keys = make([]string, 0, 0)
	for key := range params {
		if params.Get(key) != `` {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)
	//fmt.Println("=====keys", keys)

	var pList = make([]string, 0, 0)
	for _, key := range keys {
		var value = strings.TrimSpace(params.Get(key))
		if len(value) > 0 {
			pList = append(pList, key+"="+value)
		}
	}
	var src = strings.Join(pList, "&")
	src += "&key=" + privKey
	//fmt.Println("=====src", src)

	sign = MD5Sum(src)

	return sign
}

//// Signature get the sign info
//func Signature(srcdata interface{}, bizkey string) string {
//	md5ctx := md5.New()
//
//	switch v := reflect.ValueOf(srcdata); v.Kind() {
//	case reflect.String:
//		md5ctx.Write([]byte(v.String() + bizkey))
//		return hex.EncodeToString(md5ctx.Sum(nil))
//	case reflect.Map:
//		orderStr := orderParam(v.Interface(), bizkey)
//		md5ctx.Write([]byte(orderStr))
//		return hex.EncodeToString(md5ctx.Sum(nil))
//	case reflect.Struct:
//		orderStr := Struct2map(v.Interface(), bizkey)
//		md5ctx.Write([]byte(orderStr))
//		return hex.EncodeToString(md5ctx.Sum(nil))
//	default:
//		return ""
//	}
//}
//
//func orderParam(source interface{}, bizKey string) (returnStr string) {
//	switch v := source.(type) {
//	case map[string]string:
//		keys := make([]string, 0, len(v))
//
//		for k := range v {
//			if k == "sign" {
//				continue
//			}
//			keys = append(keys, k)
//		}
//		sort.Strings(keys)
//		var buf bytes.Buffer
//		for _, k := range keys {
//			if v[k] == "" {
//				continue
//			}
//			if buf.Len() > 0 {
//				buf.WriteByte('&')
//			}
//
//			buf.WriteString(k)
//			buf.WriteByte('=')
//			buf.WriteString(v[k])
//		}
//		buf.WriteString(bizKey)
//		returnStr = buf.String()
//	case map[string]interface{}:
//		keys := make([]string, 0, len(v))
//
//		for k := range v {
//			if k == "sign" {
//				continue
//			}
//			keys = append(keys, k)
//		}
//		sort.Strings(keys)
//		var buf bytes.Buffer
//		for _, k := range keys {
//			if v[k] == "" {
//				continue
//			}
//			if buf.Len() > 0 {
//				buf.WriteByte('&')
//			}
//
//			buf.WriteString(k)
//			buf.WriteByte('=')
//			// buf.WriteString(srcmap[k])
//			switch vv := v[k].(type) {
//			case string:
//				buf.WriteString(vv)
//			case int:
//				buf.WriteString(strconv.FormatInt(int64(vv), 10))
//			default:
//				panic("params type not supported")
//			}
//		}
//		buf.WriteString(bizKey)
//		returnStr = buf.String()
//	}
//	// fmt.Println(returnStr)
//	return
//}
//
//func Struct2map(content interface{}, bizKey string) string {
//	var tempArr []string
//	temString := ""
//	var val map[string]string
//	if marshalContent, err := json.Marshal(content); err != nil {
//		fmt.Println(err)
//	} else {
//		d := json.NewDecoder(bytes.NewBuffer(marshalContent))
//		d.UseNumber()
//		if err := d.Decode(&val); err != nil {
//			fmt.Println(err)
//		} else {
//			for k, v := range val {
//				val[k] = v
//			}
//		}
//	}
//	i := 0
//	for k, v := range val {
//		// 去除冗余未赋值struct
//		if v == "" {
//			continue
//		}
//		i++
//		tempArr = append(tempArr, k+"="+v)
//	}
//	sort.Strings(tempArr)
//	for n, v := range tempArr {
//		if n+1 < len(tempArr) {
//			temString = temString + v + "&"
//		} else {
//			temString = temString + v + bizKey
//		}
//	}
//	return temString
//}
