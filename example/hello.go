package main

import (
	"fmt"
	"github.com/shiningacg/mygin"
	"github.com/shiningacg/mygin/middleware"
	"log"
)

func main() {
	server := mygin.New()
	r := server.Router()
	r.Group("/group").Any("/hi").Do(func(context *mygin.Context) {
		context.Body([]byte("nihao"))
	})
	r.Any("/any").Do(func(context *mygin.Context) {
		context.Body([]byte("lala"))
	})
	r.Get("/hello").Use(TestMiddleware()).Do(func(context *mygin.Context) {
		context.Body([]byte("hello world"))
		fmt.Println("ccc")
	})
	r.Post("/user/:id").Do(func(context *mygin.Context) {
		fmt.Println(context.RouterValue("id"))
		fmt.Println(context.RouterValue("name"))
	})
	r.Any("/*").Do(func(context *mygin.Context) {
		context.Body([]byte("没有找到页面"))
	})
	r.Post("/json").Use(middleware.JsonAPI(), TestJsonApi()).Do(func(context *mygin.Context) {
		middleware.SucResponse(context, "测试json数据成功！")
	})
	r.Head("/json").Use(middleware.JsonAPI(), TestJsonApi()).Do(func(context *mygin.Context) {
		middleware.SucResponse(context, "测试head")
	})
	err := server.Run(":3112")
	if err != nil {
		log.Fatal(err)
	}
}

/*func _main()  {
	New()
}

type MyHttpServer struct {}

func (m MyHttpServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("hello"))
}

func New() {
	server := MyHttpServer{}
	handler := mygin.WrapHttpHandle(MyHttpServer{})
	router.aa().Do(handler)
}*/
