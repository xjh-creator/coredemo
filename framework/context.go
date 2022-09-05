package framework

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type Context struct {
	request *http.Request
	responseWriter http.ResponseWriter
	ctx            context.Context

	// 当前请求的handler链条
	handlers        []ControllerHandler
	index int // 当前请求调用到调用链哪个节点

	// 是否超时标记位
	hasTimeout bool
	// 写保护机制
	writerMux *sync.Mutex

	params map[string]string // url路由匹配的参数
}

func NewContext(r *http.Request,w http.ResponseWriter) *Context {
	return &Context{
		request: r,
		responseWriter: w,
		ctx:r.Context(),
		writerMux: &sync.Mutex{},
		index: -1,
	}
}

// region base function

func (ctx *Context)WriterMux() *sync.Mutex {
	return ctx.writerMux
}

func (ctx *Context)GetRequest() *http.Request {
	return ctx.request
}

func (ctx *Context)GetResponse() http.ResponseWriter {
	return ctx.responseWriter
}

func (ctx *Context)SetHasTimeout()  {
	ctx.hasTimeout = true
}

func (ctx *Context)HasTimeout() bool {
	return ctx.hasTimeout
}

// endregion

// region context function

func (ctx *Context)BaseContext() context.Context {
	return ctx.request.Context()
}

func (ctx *Context)Deadline() (deadline time.Time, ok bool) {
	return ctx.BaseContext().Deadline()
}

func (ctx *Context)Done() <-chan struct{} {
	return ctx.BaseContext().Done()
}

func (ctx *Context)Err() error {
	return ctx.BaseContext().Err()
}

func (ctx *Context)Value(key interface{}) interface{} {
	return ctx.BaseContext().Value(key)
}

// endregion

// region response function

func (ctx *Context)Json(status int,obj interface{}) error {
	if ctx.HasTimeout(){
		return nil
	}
	ctx.responseWriter.Header().Set("Content-Type","application/json")
	ctx.responseWriter.WriteHeader(status)
	byt,err := json.Marshal(obj)
	if err != nil{
		ctx.responseWriter.WriteHeader(500)
		return err
	}
	ctx.responseWriter.Write(byt)
	return nil
}

func (ctx *Context)HTML(status int,obj interface{},template string) error {
	return nil
}

func (ctx *Context)Text(status int,obj string) error {
	return nil
}

// endregion

// SetHandlers 为 context 设置 handlers
func (ctx *Context) SetHandlers(handlers []ControllerHandler) {
	ctx.handlers = handlers
}

// Next 核心函数，调用context的下一个函数
func (ctx *Context) Next() error {
	ctx.index++
	if ctx.index < len(ctx.handlers) {
		if err := ctx.handlers[ctx.index](ctx); err != nil {
			return err
		}
	}
	return nil
}

