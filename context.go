package sn

import (
	"math"
	"net/http"
)

const abortIndex int8 = math.MaxInt8 / 2

func NewContext(r *http.Request, p http.ResponseWriter) *Context {
	var caches = make(map[string]interface{})
	return &Context{
		handlers: nil,
		caches:   caches,
		Request:  r,
		Write:    p,
	}
}

// 自动生成
type Context struct {
	// 存放结果
	handlers HandlersChain
	index    int8
	caches   map[string]interface{}
	Request  *http.Request
	Write    http.ResponseWriter
}

func (c *Context) Value(key string) interface{} {
	return c.caches[key]
}

func (c *Context) Set(key string, value interface{}) {
	c.caches[key] = value
}

func (c *Context) Abort() {
	c.index = abortIndex
}

func (c *Context) IsAbort() bool {
	return c.index >= abortIndex
}

func (c *Context) Status(code int) {
	c.caches["SYS_STATUS"] = code
}

func (c *Context) Body(b []byte) {
	c.Set("SYS_BODY", b)
}

func (c *Context) GetBody() []byte {
	body := c.Value("SYS_BODY")
	if body == nil {
		return []byte{}
	}
	return body.([]byte)
}

func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}

// 设置状态码
func (c *Context) GetStatus() int {
	if code, has := c.caches["SYS_STATUS"]; has {
		return code.(int)
	}
	return 200
}

func (c *Context) Header(head string, content string) {
	if heads, has := c.caches["SYS_HEAD"]; has {
		heads.(map[string]string)[head] = content
		return
	}
	var heads = make(map[string]string)
	heads[head] = content
	c.Set("SYS_HEAD", heads)
}

func (c *Context) GetHeaders() map[string]string {
	if heads, has := c.caches["SYS_HEAD"]; has {
		return heads.(map[string]string)
	}
	return nil
}

func (c *Context) GetHeader(key string) string {
	if heads, has := c.caches["SYS_HEAD"]; has {
		if head, has := heads.(map[string]string)[key]; has {
			return head
		}
	}
	return ""
}
