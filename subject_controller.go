package main

import (
	"coredemo/framework/gin"
	"coredemo/provider/demo"
	"fmt"
)

func SubjectAddController(c *gin.Context){
	c.ISetOkStatus().IJson("ok, SubjectAddController")
}

func SubjectListController(c *gin.Context){
	// 获取demo服务实例
	demoService := c.MustMake(demo.Key).(demo.Service)

	// 调用服务实例的方法
	foo := demoService.GetFoo()

	// 输出结果
	c.ISetOkStatus().IJson(foo)
}

func SubjectDelController(c *gin.Context){
	c.ISetOkStatus().IJson("ok, SubjectDelController")
}

func SubjectUpdateController(c *gin.Context){
	c.ISetOkStatus().IJson("ok, SubjectUpdateController")
}

func SubjectGetController(c *gin.Context){
	subjectId, _ := c.DefaultParamInt("id", 0)
	c.ISetOkStatus().IJson("ok, SubjectGetController:" + fmt.Sprint(subjectId))
}

func SubjectNameController(c *gin.Context){
	c.ISetOkStatus().IJson("ok, SubjectNameController")
}