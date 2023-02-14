package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/0x00b/gobbq/engine/entity"
	"google.golang.org/protobuf/proto"
)

type BasicAuth struct{ User, Passwd string }

type HttpParam struct {
	http.Header
	http.RoundTripper
	http.CookieJar
	*BasicAuth
}

type HttpHooker interface {
	Before(entity.Context, *http.Request) (After func(any, error) error)
}

var (
	Hooker HttpHooker
)

// HTTP http远程调用
func HTTP(ctx entity.Context, method, url string,
	request any, response any, params ...*HttpParam) (rspByte []byte, e error) {
	var requestBody []byte
	if request != nil {
		if p, ok := request.(proto.Message); ok {
			requestBody, _ = proto.Marshal(p)
		} else {
			typ := reflect.TypeOf(request)
			val := reflect.ValueOf(request)
			for typ.Kind() == reflect.Ptr {
				typ = typ.Elem()
				val = val.Elem()
			}
			switch typ.Kind() {
			case reflect.String:
				requestBody = []byte(val.Interface().(string))
			default:
				requestBody, e = json.Marshal(request)
				if e != nil {
					return nil, e
				}
			}
		}
	}
	contentType := http.DetectContentType(requestBody)

	var req *http.Request
	req, e = http.NewRequest(method,
		url,
		bytes.NewBuffer(requestBody))
	if e != nil {
		return nil, e
	}

	var param *HttpParam
	if len(params) != 0 {
		param = params[0]
	}
	if param.Header == nil {
		param.Header = make(http.Header)
	}
	client := &http.Client{}
	if param != nil {
		if param.RoundTripper != nil {
			client.Transport = param.RoundTripper
		}
		if param.Header != nil {
			req.Header = param.Header
		}
		if param.CookieJar != nil {
			client.Jar = param.CookieJar
		}
		if param.BasicAuth != nil {
			req.SetBasicAuth(param.BasicAuth.User, param.BasicAuth.Passwd)
		}
	}

	req.Header.Set("Content-Type", contentType)

	var responseBody []byte

	after := func(any, error) error { return nil }
	if Hooker != nil {
		// ctx := WithMeta(ctx, &ago.MetaInfo{
		// 	Action: url,
		// })
		after = Hooker.Before(ctx, req)
	}
	defer func() {
		e = after(response, e)
	}()

	resp, e := client.Do(req)
	if e != nil {
		return nil, e
	}
	defer resp.Body.Close()

	responseBody, e = ioutil.ReadAll(resp.Body)
	if e != nil {
		return nil, e
	}
	if resp.StatusCode/100 != 2 {
		return responseBody, fmt.Errorf("%d:%s", resp.StatusCode, resp.Status)
	}

	if response != nil {
		if r, ok := response.(proto.Message); ok {
			e = proto.Unmarshal(rspByte, r)
		} else {
			e = json.Unmarshal(rspByte, response)
		}
		if e != nil {
			return rspByte, e
		}
	}

	return rspByte, nil
}
