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

func (r *JsonResponse) Encode() ([]byte, error) {
	b, err := json.Marshal(r.m)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (r *JsonResponse) IsEmpty() bool {
	return len(r.m) == 0
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
		if err := context.GetError(); err != nil {
			errResponse(rsp, -1, err)
		}
		if rsp.IsEmpty() {
			sucResponse(rsp, nil)
		}
		b, err := rsp.Encode()
		if err != nil {
			rsp = NewJsonResponse()
			errResponse(rsp, 500, errors.New("服务内部错误：数据无法转换"))
			b, _ = rsp.Encode()
		}
		context.Body(b)
	}
}

func IsJsonAPILoad(ctx *mygin.Context) bool {
	return ctx.Value(JsonApiKey) == nil
}

func JsonRawResponse(ctx *mygin.Context) *JsonResponse {
	if rsp := ctx.Value(JsonApiKey); rsp != nil {
		return rsp.(*JsonResponse)
	}
	panic("JsonApi中间件没有被正确加载")
}
