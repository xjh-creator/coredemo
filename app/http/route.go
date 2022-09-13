package http

import (
	"coredemo/framework/gin"
	"coredemo/provider/demo"
)

func Routes(r *gin.Engine) {

	r.Static("/dist/", "./dist/")

	demo.Register(r)
}
