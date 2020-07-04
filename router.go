package mygin

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
	path     string
	engine   *Engine
	handlers HandlersChain
}

func (r *routerGroup) Group(path string) RouterGroup {
	handlers := make(HandlersChain, len(r.handlers), len(r.handlers)+3)
	// 继承中间件
	copy(handlers, r.handlers)
	// TODO:针对是否有'/'进行处理，防止误写
	return &routerGroup{
		path:     r.path + path,
		engine:   r.engine,
		handlers: handlers,
	}
}

func (r *routerGroup) handle(method, path string) *router {
	handlers := make(HandlersChain, len(r.handlers), len(r.handlers)+3)
	copy(handlers, r.handlers)
	return &router{
		method:   method,
		path:     r.path + path,
		engine:   r.engine,
		handlers: handlers,
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

func (r *routerGroup) Use(handles ...HandlerFunc) RouterGroup {
	r.handlers = append(r.handlers, handles...)
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
	method   string
	path     string
	engine   *Engine
	handlers HandlersChain
}

func (r *router) Use(handles ...HandlerFunc) Router {
	r.handlers = append(r.handlers, handles...)
	return r
}

func (r *router) Do(handlerFunc HandlerFunc) {
	// 添加方法
	r.engine.addRouter(r.method, r.path, func(c *Context) {
		// 附加方法
		c.handlers = append(r.handlers, handlerFunc)
		// 启动处理
		c.handlers[0](c)
	})
}
