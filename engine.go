package mygin

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

var (
	default404Body = []byte("404 page not found")
	default405Body = []byte("405 method not allowed")
	default500Body = []byte("500 internal server error")
	default400Body = []byte("400 bad request")

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

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := NewContext(req, w)
	c.Request = req
	c.Write = w
	e.handleHTTPRequest(c)
}

func (e *Engine) handleHTTPRequest(ctx *Context) {
	// 路由匹配
	handle := e.tree.match(ctx.Request.RequestURI, ctx.Request.Method)
	// 开始进行处理
	handle(ctx)
	// 写body和status
	if ctx.IsProto() {
		writeHeader(ctx)
		writeBody(ctx)
	}
}

func writeHeader(ctx *Context) {
	// 自动设置content-length，因为head会刷新掉
	if body := ctx.GetBody(); ctx.GetHeader("Content-Length") == "" && body != nil {
		ctx.Header("Content-Length", fmt.Sprint(len(body)))
	}
	headers := ctx.GetHeaders()
	for key, value := range headers {
		ctx.Write.Header().Set(key, value)
	}
	ctx.Write.WriteHeader(ctx.GetStatus())
}

func writeBody(ctx *Context) {
	// 如果是head请求，不允许返回body内容
	if ctx.Request.Method == HEAD {
		ctx.Body(nil)
	}
	body := ctx.GetBody()
	if code := ctx.GetStatus(); body == nil && code != 200 {
		if err := ctx.GetError(); err != nil {
			body = []byte(err.Error())
		} else {
			switch code {
			case 400:
				body = default400Body
			case 500:
				body = default500Body
			case 404:
				body = default404Body
			default:
				body = default500Body
			}
		}
	}
	for {
		n, err := ctx.Write.Write(body)
		if err != nil {
			if err != io.EOF {
				log.Print(err)
			}
			break
		}
		if n == len(body) {
			break
		}
	}
	err := ctx.Request.Body.Close()
	if err != nil {
		log.Print(err)
	}
}

func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

func (e *Engine) Router() RouterGroup {
	return e.tree.root
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
