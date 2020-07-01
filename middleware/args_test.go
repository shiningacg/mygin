package middleware

import (
	"fmt"
	"github.com/shlande/mygin"
	"log"
	"testing"
)

type DTO struct {
	Hi  string
	Ok  int
	Ahh float64 `json:"mmm"`
}

func TestArgs(t *testing.T) {
	mygin := mygin.New()
	root := mygin.Router()
	root.Post("/").Use(Args()).Do(func(context *mygin.Context) {
		dto := &DTO{}
		err := Merge(context, dto)
		if err != nil {
			log.Print(err)
		}
		fmt.Println(dto)
	})
	err := mygin.Run(":3112")
	panic(err)
}
