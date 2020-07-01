package middleware

import (
	"fmt"
	"github.com/shlande/sn"
	"log"
	"testing"
)

type DTO struct {
	Hi  string
	Ok  int
	Ahh float64 `json:"mmm"`
}

func TestArgs(t *testing.T) {
	mygin := sn.New()
	root := mygin.Router()
	root.Post("/").Use(Args()).Do(func(context *sn.Context) {
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
