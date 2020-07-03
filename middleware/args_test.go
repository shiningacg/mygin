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
	Ahh float64 `json:"mmm" args:"required"`
}

func TestArgs(t *testing.T) {
	server := mygin.New()
	testArgsPost(server)
	err := server.Run(":3112")
	panic(err)
}

func testArgsGet(server *mygin.Engine) {
	root := server.Router()
	root.Get("/").Use(Args(&DTO{})).Do(func(context *mygin.Context) {
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
	root.Post("/").Use(JsonAPI(), Args(&DTO{})).Do(func(context *mygin.Context) {
		dto := &DTO{}
		err := Merge(context, dto)
		if err != nil {
			log.Print(err)
		}
		fmt.Println("dto:", dto)
		SucResponse(context, dto)
	})
}
