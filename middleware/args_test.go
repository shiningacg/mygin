package middleware

import (
	"fmt"
	"github.com/shiningacg/mygin"
	"log"
	"testing"
)

type DTO struct {
	Hi  string
	Ok  int
	Ahh float64 `json:"mmm"`
}

func TestArgs(t *testing.T) {
	server := mygin.New()
	testArgsGet(server)
	err := server.Run(":3112")
	panic(err)
}

func testArgsGet(server *mygin.Engine) {
	root := server.Router()
	root.Get("/").Use(Args()).Do(func(context *mygin.Context) {
		dto := &DTO{}
		err := Merge(context, dto)
		if err != nil {
			log.Print(err)
		}
		fmt.Println(dto)
	})
}

func testArgsPost(server *mygin.Engine) {
	root := server.Router()
	root.Post("/").Use(Args()).Do(func(context *mygin.Context) {
		dto := &DTO{}
		err := Merge(context, dto)
		if err != nil {
			log.Print(err)
		}
		fmt.Println(dto)
	})
}
