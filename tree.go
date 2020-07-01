package mygin

import (
	"strings"
)

type node struct {
	path    string
	handler map[string]HandlerFunc
}

func (n *node) add(httpMethod string, handlerFunc HandlerFunc) {
	n.handler[httpMethod] = handlerFunc
}

func (n *node) match(method string) HandlerFunc {
	var (
		has    bool
		handle HandlerFunc
	)
	if handle, has = n.handler[method]; !has {
		if handle, has = n.handler[ANY]; has {
			return handle
		}
		return bodyWriter(default404Body)
	}
	return handle
}

func NewTree() *tree {
	return &tree{
		nodes: make(map[string]*node),
	}
}

type tree struct {
	nodes map[string]*node
}

func (t *tree) addRouter(httpMethod, path string, handler HandlerFunc) {
	n, has := t.nodes[path]
	if !has {
		handler := make(map[string]HandlerFunc)
		n = &node{
			path:    path,
			handler: handler,
		}
	}
	n.add(httpMethod, handler)
	t.nodes[path] = n
}

func (t *tree) match(path, method string) HandlerFunc {
	var (
		n   *node
		has bool
	)
	path = parsePath(path)
	if n, has = t.nodes[path]; !has {
		return bodyWriter(default404Body)
	}
	return n.match(method)
}

// 处理一下请求，但是仅仅是get类
func parsePath(path string) string {
	raw := strings.Split(path, "?")
	if len(raw) == 1 {
		return raw[0]
	}
	return strings.Join(raw[:len(raw)-1], "")
}
