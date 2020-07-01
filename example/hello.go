package main

import (
	"fmt"
	"github.com/shlande/sn"
	"log"
)

func main() {
	server := sn.New()
	r := server.Router()
	r.Get("/hello").Use(TestMiddleware()).Do(func(context *sn.Context) {
		context.Body([]byte("hello world"))
		fmt.Println("ccc")
	})
	err := server.Run(":3112")
	if err != nil {
		log.Fatal(err)
	}
}
