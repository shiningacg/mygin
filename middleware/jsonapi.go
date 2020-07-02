package middleware

import (
	"encoding/json"
	"errors"
	"github.com/shiningacg/mygin"
)

func NewJsonResponse() *JsonResponse {
	return &JsonResponse{
		m: make(map[string]interface{}),
	}
}

type JsonResponse struct {
	m map[string]interface{}
}

func (r *JsonResponse) Set(key string, value interface{}) {
	r.m[key] = value
}

func (r *JsonResponse) Value(key string) interface{} {
	return r.m[key]
}

func (r *JsonResponse) Encode() []byte {
	b, err := json.Marshal(r.m)
	if err != nil {
		return nil
	}
	return b
}

const JsonApiKey = "JsonApi"

func ErrResponse(ctx *mygin.Context, code int, err error) {
	errResponse(JsonRawResponse(ctx), code, err)
}

func errResponse(rr *JsonResponse, code int, err error) {
	rr.Set("code", code)
	rr.Set("message", err.Error())
}

func SucResponse(ctx *mygin.Context, data interface{}) {
	sucResponse(JsonRawResponse(ctx), data)
}

func sucResponse(rr *JsonResponse, data interface{}) {
	rr.Set("code", 200)
	rr.Set("message", "ok")
	rr.Set("data", data)
}

func JsonAPI() mygin.HandlerFunc {
	return func(context *mygin.Context) {
		rsp := NewJsonResponse()
		context.Set(JsonApiKey, rsp)
		context.Next()
		b := rsp.Encode()
		if b == nil {
			rsp = NewJsonResponse()
			errResponse(rsp, 500, errors.New("服务内部错误：数据无法转换"))
			b = rsp.Encode()
		}
		context.Body(b)
	}
}

func JsonRawResponse(ctx *mygin.Context) *JsonResponse {
	if rsp := ctx.Value(JsonApiKey); rsp != nil {
		return rsp.(*JsonResponse)
	}
	panic("JsonApi中间件没有被正确加载")
}
