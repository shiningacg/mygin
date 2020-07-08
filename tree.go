package mygin

import "fmt"

func NewTree() *tree {
	return &tree{
		root: &routerGroup{},
	}
}

type tree struct {
	root *routerGroup
}

func (t *tree) match(path, method string) HandlerFunc {
	nodes := getNodeFromPath(path)
	fmt.Println(nodes)
	g := t.root
	var handler HandlerFunc
	for i, node := range nodes {
		temp := g.getRouterGroup(node)
		if temp == nil {
			handler = g.handler[method]
			break
		}
		g = temp
		if i == len(nodes)-1 {
			handler = g.handler[method]
		}
	}
	if handler == nil {
		return bodyWriter(default404Body)
	}
	return handler
}
