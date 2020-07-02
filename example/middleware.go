package main

import (
	"fmt"
	"github.com/shlande/mygin"
	"github.com/shlande/mygin/middleware"
)

func TestMiddleware() mygin.HandlerFunc {
	md := MiddleWare{}
	return md.Handle
}

type MiddleWare struct{}

func (m *MiddleWare) Handle(c *mygin.Context) {
	fmt.Println("hihihih")
	c.Next()
	fmt.Println("aaaa")
}

func TestJsonApi() mygin.HandlerFunc {
	return func(context *mygin.Context) {
		context.Next()
		middleware.JsonRawResponse(context).Set("traceId", "111")
	}
}
