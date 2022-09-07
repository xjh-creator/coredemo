package main

import (
	"context"
	"coredemo/framework"
	"fmt"
	"log"
	"time"
)

// FooControllerHandler 处理业务的控制器
func FooControllerHandler(c *framework.Context) error {
	finish := make(chan struct{}, 1)
	panicChan := make(chan interface{}, 1)

	// 一 、生成一个超时的Context
	durationCtx, cancel := context.WithTimeout(c.BaseContext(), time.Duration(1*time.Second))
	defer cancel()

	// mu := sync.Mutex{}
	//二 、开启新的 goroutine 处理业务逻辑
	go func() {
		defer func() {
			if p := recover(); p != nil {
				panicChan <- p
			}
		}()
		// Do real action
		time.Sleep(10 * time.Second)
		c.SetOkStatus().Json("ok")

		//新的 goroutine 结束时通过一个信号告知父 goroutine
		finish <- struct{}{}
	}()
	//三 、使用 select 关键字来监听三个事件：异常事件、结束实际、超时实际
	select {
	case p := <-panicChan:
		c.WriterMux().Lock()
		defer c.WriterMux().Unlock()
		log.Println(p)
		c.SetStatus(500).Json("panic")
	case <-finish:
		fmt.Println("finish")
	case <-durationCtx.Done():
		c.WriterMux().Lock()
		defer c.WriterMux().Unlock()
		c.SetStatus(500).Json("time out")
		c.SetHasTimeout()
	}
	return nil
}
