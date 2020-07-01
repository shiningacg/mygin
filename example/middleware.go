package main

import (
	"fmt"
	"github.com/shlande/sn"
)

func TestMiddleware() sn.HandlerFunc {
	md := MiddleWare{}
	return md.Handle
}

type MiddleWare struct{}

func (m *MiddleWare) Handle(c *sn.Context) {
	fmt.Println("hihihih")
	c.Next()
	fmt.Println("aaaa")
}
