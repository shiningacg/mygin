package main

import (
	"fmt"
	"github.com/shlande/mygin"
	"github.com/shlande/mygin/middleware"
	"log"
)

func main() {
	server := mygin.New()
	r := server.Router()
	r.Get("/hello").Use(TestMiddleware()).Do(func(context *mygin.Context) {
		context.Body([]byte("hello world"))
		fmt.Println("ccc")
	})
	r.Post("/json").Use(middleware.JsonAPI(), TestJsonApi()).Do(func(context *mygin.Context) {
		middleware.SucResponse(context, "测试json数据成功！")
	})
	err := server.Run(":3112")
	if err != nil {
		log.Fatal(err)
	}
}
