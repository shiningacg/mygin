package sn

import (
	"errors"
	"log"
	"net/http"
	"net/url"
)

var (
	default404Body = []byte("404 page not found")
	default405Body = []byte("405 method not allowed")

	ErrReachLimitSize = errors.New("发送的数据大小超过限制")
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
	c.Request = req
	c.Write = w
	e.handleHTTPRequest(c)
}

func (e *Engine) handleHTTPRequest(ctx *Context) {
	req := ctx.Request
	w := ctx.Write
	// 参数获取
	if req.Method == GET {
		// get参数获取
		e.parseGet(ctx)
	}
	// TODO：form数据获取
	// 路由匹配
	handle := e.tree.match(req.RequestURI, req.Method)
	// 开始进行处理
	handle(ctx)
	// 写body和status
	w.WriteHeader(ctx.GetStatus())
	writeHeader(ctx)
	_, err := w.Write(ctx.GetBody())
	if err != nil {
		log.Print(err)
	}
	err = req.Body.Close()
	if err != nil {
		log.Print(err)
	}
}

func (e *Engine) parseGet(ctx *Context) {
	u, err := url.Parse(ctx.Request.RequestURI)
	if err != nil {
		return
	}
	m := u.Query()
	for name, value := range m {
		ctx.Set(name, value[0])
	}
}

func writeHeader(ctx *Context) {
	headers := ctx.GetHeaders()
	for key, value := range headers {
		ctx.Write.Header().Set(key, value)
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
