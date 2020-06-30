package main

import (
	"github.com/shlande/sn"
	"log"
)

func main() {
	server := sn.New()
	r := server.Router()
	r.Get("/hello").Do(func(context *sn.Context) {
		context.Body([]byte("hello world"))
	})
	err := server.Run(":3112")
	if err != nil {
		log.Fatal(err)
	}
}
