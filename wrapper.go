package mygin

import "net/http"

func WrapHttpHandle(handler http.Handler) func(ctx *Context) {
	return func(ctx *Context) {
		ctx.Proto()
		handler.ServeHTTP(ctx.Write, ctx.Request)
	}
}
