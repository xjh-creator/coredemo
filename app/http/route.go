package http

import (
	"coredemo/framework/gin"
	"coredemo/app/http/module/demo"
)

func Routes(r *gin.Engine) {

	r.Static("/dist/", "./dist/")

	demo.Register(r)
}
