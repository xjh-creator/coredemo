package http

import "coredemo/framework/gin"

// NewHttpEngine is command
func NewHttpEngine() (*gin.Engine, error) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	Routes(r)
	return r, nil
}
