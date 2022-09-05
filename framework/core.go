package framework

import (
	"log"
	"net/http"
	"strings"
)

//Core 框架核心结构
type Core struct {
	router map[string]*Tree
	middlewares []ControllerHandler //在core设置中间件
}

func NewCore() *Core {
	// 初始化路由
	router := map[string]*Tree{}
	router["GET"] = NewTree()
	router["POST"] = NewTree()
	router["PUT"] = NewTree()
	router["DELETE"] = NewTree()

	core := &Core{router: router}
	return core
}

// Use 注册中间件
func (c *Core) Use(middlewares ...ControllerHandler) {
	c.middlewares = append(c.middlewares,middlewares...)
}

// === http method wrap

// Get 匹配 GET 方法, 增加路由规则
func (c *Core) Get(url string, handlers ...ControllerHandler) {
	// 将core的middleware 和 handlers结合起来
	allHandlers := append(c.middlewares, handlers...)
	if err := c.router["GET"].AddRouter(url, allHandlers); err != nil {
		log.Fatal("add router error: ", err)
	}
}

// Post 匹配 POST 方法, 增加路由规则
func (c *Core) Post(url string, handlers ...ControllerHandler) {
	allHandlers := append(c.middlewares, handlers...)
	if err := c.router["POST"].AddRouter(url, allHandlers); err != nil {
		log.Fatal("add router error: ", err)
	}
}

// Put 匹配 PUT 方法, 增加路由规则
func (c *Core) Put(url string, handlers ...ControllerHandler) {
	allHandlers := append(c.middlewares, handlers...)
	if err := c.router["PUT"].AddRouter(url, allHandlers); err != nil {
		log.Fatal("add router error: ", err)
	}
}

// Delete 匹配 DELETE 方法, 增加路由规则
func (c *Core) Delete(url string, handlers ...ControllerHandler) {
	allHandlers := append(c.middlewares, handlers...)
	if err := c.router["DELETE"].AddRouter(url, allHandlers); err != nil {
		log.Fatal("add router error: ", err)
	}
}

// ==== http method wrap end

// FindRouteNodeByRequest 匹配路由，如果没有匹配到，返回nil
func (c *Core) FindRouteNodeByRequest(request *http.Request) *node {
	// uri 和 method 全部转换为大写，保证大小写不敏感
	uri := request.URL.Path
	method := request.Method
	upperMethod := strings.ToUpper(method)

	// 查找第一层map
	if methodHandlers, ok := c.router[upperMethod]; ok {
		return methodHandlers.root.matchNode(uri)
	}
	return nil
}

func (c *Core)ServeHTTP(response http.ResponseWriter,request *http.Request)  {
	// 封装自定义context
	ctx := NewContext(request, response)

	// 寻找路由
	node := c.FindRouteNodeByRequest(request)
	if node == nil {
		// 如果没有找到，这里打印日志
		ctx.SetStatus(404).Json("not found")
		return
	}

	// 设置 context 中的 handlers 字段
	ctx.SetHandlers(node.handlers)

	// 设置参数
	params := node.parseParamsFromEndNode(request.URL.Path)
	ctx.SetParams(params)

	// 调用路由函数，如果返回 err 代表内部错误，返回500
	if err := ctx.Next();err != nil{
		ctx.SetStatus(500).Json("inner error")
		return
	}
}

func (c *Core)Group(prefix string) IGroup {
	return NewGroup(c,prefix)
}
