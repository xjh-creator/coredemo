package framework

import (
	"net/http"
	"strings"
)

//Core 框架核心结构
type Core struct {
	router map[string]map[string]ControllerHandler
}

func NewCore() *Core {
	//定义二级map
	getRouter := map[string]ControllerHandler{}
	postRouter := map[string]ControllerHandler{}
	putRouter := map[string]ControllerHandler{}
	deleteRouter := map[string]ControllerHandler{}

	//将二级map写入一级map
	router := map[string]map[string]ControllerHandler{}
	router["GET"] = getRouter
	router["POST"] = postRouter
	router["PUT"] = putRouter
	router["DELETE"] = deleteRouter

	return &Core{router: router}
}

// Get 对应 Method = GET
func (c *Core)Get(url string,handler ControllerHandler)  {
	upperUrl := strings.ToUpper(url)
	c.router["GET"][upperUrl] = handler
}

// Post 对应 Method = POST
func (c *Core)Post(url string,handler ControllerHandler)  {
	upperUrl := strings.ToUpper(url)
	c.router["POST"][upperUrl] = handler
}

// Put 对应 Method = PUT
func (c *Core)Put(url string,handler ControllerHandler)  {
	upperUrl := strings.ToUpper(url)
	c.router["PUT"][upperUrl] = handler
}

// Delete 对应 Method = DELETE
func (c *Core)Delete(url string,handler ControllerHandler)  {
	upperUrl := strings.ToUpper(url)
	c.router["DELETE"][upperUrl] = handler
}

// FindRouteByRequest 匹配路由，如果匹配不到，则返回nil
func (c *Core)FindRouteByRequest(request *http.Request) ControllerHandler {
	// url 和 method 全部转为大写，保证大小写不敏感
	url := request.URL.Path
	method := request.Method
	upperMethod := strings.ToUpper(method)
	upperUrl := strings.ToUpper(url)

	//查找第一层map
	if methodHandlers,ok := c.router[upperMethod];ok{
		//查找第二层
		if handler,ok := methodHandlers[upperUrl];ok{
			return handler
		}
	}
	return nil
}

func (c *Core)ServeHTTP(response http.ResponseWriter,request *http.Request)  {
	// 封装自定义context
	ctx := NewContext(request, response)

	// 寻找路由
	router := c.FindRouteByRequest(request)
	if router == nil{
		// 没有找到，这里打印日志
		ctx.Json(404,"not found")
		return
	}

	// 调用路由函数，如果返回 err 代表内部错误，返回500
	if err := router(ctx);err != nil{
		ctx.Json(500,"inner error")
		return
	}
}

func (c *Core)Group(prefix string) IGroup {
	return NewGroup(c,prefix)
}
