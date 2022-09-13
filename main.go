package main

import (
	"context"
	"coredemo/framework/gin"
	"coredemo/framework/middleware"
	"coredemo/app/provider/demo"
	"coredemo/framework/provider/app"
	hadeHttp "coredemo/app/http"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main()  {
	core := gin.New()

	// 绑定具体的服务
	core.Bind(&app.HadeAppProvider{})
	core.Bind(&demo.DemoProvider{})

	core.Use(gin.Recovery())
	core.Use(middleware.Cost())

	hadeHttp.Routes(core)

	server := &http.Server{
		Handler: core,
		Addr: ":8080",
	}
	// 这个goroutine是启动服务的goroutine
	go func() {
		server.ListenAndServe()
	}()

	// 当前的goroutine等待信号量
	quit := make(chan os.Signal)
	// 监控信号：SIGINT, SIGTERM, SIGQUIT
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// 这里会阻塞当前goroutine等待信号
	<-quit

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 调用Server.Shutdown graceful结束
	if err := server.Shutdown(timeoutCtx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
}
