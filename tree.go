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

func (n *node) match(method string, args []string) HandlerFunc {
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
	// 在这里包一层用于路由参数处理
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
	// 获取不包含参数的路径
	basePath := getBasePath(path)
	// 如果存在处理函数，那么进行多层包装
	n, has := t.nodes[basePath]
	if !has {
		handler := make(map[string]HandlerFunc)
		n = &node{
			path:    path,
			handler: handler,
		}
	}
	handler = n.wrapHandler(handler, n.handler[httpMethod], matcher)
	n.add(httpMethod, handler)
	t.nodes[path] = n
}

// 在添加路由的时候去除的参数
func getBasePath(path string) string {

}

// 在添加路由的时候获取参数
func getArgsPath(path string) string {
	basePath := getBasePath(path)
	if path == basePath {
		return ""
	}
	return path[len(basePath)+1:]
}

// 通过可变参数的形势获取match函数
func matchArgsTemplate(args []string, template string) bool {
	if template == "" {
		return args == nil
	}
	if template == "/*" {
		return true
	}
	temps := strings.Split(template[1:], "/")
	return len(temps) == len(args)
}

func (n *node) wrapHandler(wrapper, in HandlerFunc, argsTemplate string) HandlerFunc {
	return func(ctx *Context) {
		args := n.getArgs(ctx.Request.RequestURI)
		if matchArgsTemplate(args, argsTemplate) {
			// 设置路由参数

			wrapper(ctx)
		}
		if in == nil {
			in = bodyWriter(default404Body)
		}
		in(ctx)
	}
}

func (n *node) getArgs(s string) []string {
	s = s[len(n.path):]
	if s == "" {
		return nil
	}
	return strings.Split(s[1:], "/")
}

func (t *tree) match(path, method string) HandlerFunc {
	var (
		n   *node
		has bool
	)
	path = parsePath(path)
	// 便利，直到找到node
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
