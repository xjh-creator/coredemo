package framework

import "net/http"

//Core 框架核心结构
type Core struct {

}

func NewCore() *Core {
	return &Core{}
}

func (c *Core)ServeHTTP(response http.ResponseWriter,request *http.Request)  {
	//TODD
}
