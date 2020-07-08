package mygin

import (
	"strings"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
	PATCH  = "PATCH"
	ANY    = "ANY"
)

type RouterGroup interface {
	Group(path string) RouterGroup // 向下生成子路由

	// http方法
	Get(path string) Router
	Post(path string) Router
	Put(path string) Router
	Delete(path string) Router
	Patch(path string) Router
	Any(path string) Router
	Use(handles ...HandlerFunc) RouterGroup
}

type routerGroup struct {
	path           string
	name           string
	subRouterGroup map[string]*routerGroup
	// 管理处理方法用
	handler map[string]HandlerFunc
	// 存放中间件
	middleware HandlersChain
}

//
func (r *routerGroup) Group(name string) RouterGroup {
	// 判断path，检查是否以'/'开始
	if name[0] != '/' {
		panic("路径必须以'/'开始")
	}
	// 判断是否有多个节点
	nodes := getNodeFromPath(name[1:])
	if len(nodes) != 1 {
		panic("不能添加多个节点或空节点")
	}
	r.addRouterGroup(name)
	return r.getRouterGroup(name)
}

func (r *routerGroup) GetHandler(method string) HandlerFunc {
	handler, has := r.handler[method]
	if !has {
		handler = r.handler[ANY]
	}
	if handler == nil {
		return bodyWriter(default404Body)
	}
	return handler
}

// 处理一下请求，但是仅仅是get类
func parsePath(path string) string {
	raw := strings.Split(path, "?")
	if len(raw) == 1 {
		return raw[0]
	}
	return strings.Join(raw[:len(raw)-1], "")
}

// 获取路由的节点数
func getNodeFromPath(nodeString string) []string {
	nodeString = parsePath(nodeString)
	if nodeString == "" {
		return nil
	}
	return strings.Split(nodeString, "/")[1:]
}

func getNodeFromTemplate(template string) (nodeName []string, args []string) {
	var node string
	var i int
	arr := strings.Split(template, "/")[1:]
	for i, node = range arr {
		if prefix := node[0]; prefix == ':' || prefix == '*' {
			nodeName = arr[:i]
			args = arr[i:]
			goto END
		}
	}
	// 没有参数
	nodeName = arr
	return
END:
	for i, arg := range args {
		if arg == "*" {
			continue
		}
		args[i] = arg[1:]
	}
	return
}

func (r *routerGroup) getRouterGroup(nodeName string) *routerGroup {
	return r.subRouterGroup[nodeName]
}

// 不允许越级添加，这是底层实现
func (r *routerGroup) addRouterGroup(nodeName string) {
	if _, has := r.subRouterGroup[nodeName]; has {
		return
	}
	md := make(HandlersChain, len(r.middleware), len(r.middleware)+3)
	// 继承中间件
	copy(md, r.middleware)
	// 初始化sub数组
	if r.subRouterGroup == nil {
		r.subRouterGroup = make(map[string]*routerGroup)
	}
	r.subRouterGroup[nodeName] = &routerGroup{
		path:       r.path + "/" + nodeName,
		name:       nodeName,
		middleware: md,
	}
}

func (r *routerGroup) addHandler(method string, handlerFunc HandlerFunc) {
	if r.handler == nil {
		r.handler = make(map[string]HandlerFunc)
	}
	r.handler[method] = handlerFunc
}

func (r *routerGroup) addRouter(rt *router) {
	var node *routerGroup
	names, args := getNodeFromTemplate(rt.matchTemplate)
	for _, name := range names {
		r.addRouterGroup(name)
		node = r.getRouterGroup(name)
	}
	if node == nil {
		node = rt.group
	}
	node.addHandler(rt.method, func(c *Context) {
		argstr := c.Request.RequestURI[len(node.path):]
		inputs := getNodeFromPath(argstr)
		// 判断是否是全局匹配
		if len(args) == 1 && args[0] == "*" {
			goto HANDLE
		}
		// 参数个数不匹配
		if len(args) != len(inputs) {
			bodyWriter(default404Body)(c)
			return
		}
		// 路由中添加方法
		for i, key := range args {
			c.setRouterValue(key, inputs[i])
		}
	HANDLE:
		// 附加方法
		c.handlers = rt.handlers
		// 启动处理
		c.handlers[0](c)
	})
}

// 用来绑定处理方法
func (r *routerGroup) handle(method, path string) *router {
	handlers := make(HandlersChain, len(r.middleware), len(r.middleware)+3)
	copy(handlers, r.middleware)
	return &router{
		group:         r,
		method:        method,
		matchTemplate: path,
		handlers:      handlers,
	}
}

func (r *routerGroup) Get(path string) Router {
	return r.handle(GET, path)
}

func (r *routerGroup) Post(path string) Router {
	return r.handle(POST, path)
}

func (r *routerGroup) Put(path string) Router {
	return r.handle(PUT, path)
}

func (r *routerGroup) Delete(path string) Router {
	return r.handle(DELETE, path)
}

func (r *routerGroup) Patch(path string) Router {
	return r.handle(PATCH, path)
}

func (r *routerGroup) Any(path string) Router {
	return r.handle(ANY, path)
}

func (r *routerGroup) Use(mds ...HandlerFunc) RouterGroup {
	r.middleware = append(r.middleware, mds...)
	return r
}

// 路由对象，用来设置需要添加的处理方法和中间件
type Router interface {
	// 添加中间件
	Use(handles ...HandlerFunc) Router
	// 设置处理方法
	Do(handlerFunc HandlerFunc)
}

type router struct {
	method        string
	matchTemplate string
	group         *routerGroup
	handlers      HandlersChain
	end           bool
}

func (r *router) Use(handles ...HandlerFunc) Router {
	if r.end {
		panic("无法在添加了处理方法后再添加中间件！")
	}
	r.handlers = append(r.handlers, handles...)
	return r
}

func (r *router) Do(handlerFunc HandlerFunc) {
	// 添加方法
	r.Use(handlerFunc)
	r.end = true
	r.group.addRouter(r)
}
