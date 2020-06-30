package sn

import (
	"log"
	"net/http"
)

var (
	default404Body = []byte("404 page not found")
	default405Body = []byte("405 method not allowed")
)

// 中间件的实际调用的处理方法
type HandlerFunc func(*Context)

type HandlersChain []HandlerFunc

type Engine struct {
	// 路由表
	tree *tree
	//
}

func (e *Engine) addRouter(httpMethod, path string, handler HandlerFunc) {
	e.tree.addRouter(httpMethod, path, handler)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := NewContext(req, w)
	c.req = req
	c.p = w
	e.handleHTTPRequest(c)
}

func (e *Engine) handleHTTPRequest(ctx *Context) {
	// 解析路由
	req := ctx.req
	w := ctx.p
	e.tree.match(req.RequestURI, req.Method)(ctx)
	// 写body和status
	w.WriteHeader(ctx.GetStatus())
	writeHeader(ctx)
	_, err := w.Write(ctx.GetBody())
	if err != nil {
		log.Print(err)
	}
}

func writeHeader(ctx *Context) {
	headers := ctx.GetHeaders()
	for key, value := range headers {
		ctx.p.Header().Set(key, value)
	}
}

func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

func (e *Engine) Router() RouterGroup {
	return &routerGroup{
		path:     "/",
		engine:   e,
		handlers: make(HandlersChain, 0, 9),
	}
}

// 新建实例
func New() *Engine {
	engine := &Engine{
		tree: NewTree(),
	}
	return engine
}

func bodyWriter(body []byte) HandlerFunc {
	return func(context *Context) {
		context.Body(body)
	}
}
