package middleware

import (
	"context"
	"coredemo/framework"
	"fmt"
	"log"
	"time"
)

func TimeoutHandler(fun framework.ControllerHandler, d time.Duration) framework.ControllerHandler {
	// 使用函数回调
	return func(c *framework.Context) error {

		finish := make(chan struct{}, 1)
		panicChan := make(chan interface{}, 1)

		// 执行业务逻辑前预操作：初始化超时 context
		durationCtx, cancel := context.WithTimeout(c.BaseContext(), d)
		defer cancel()

		c.GetRequest().WithContext(durationCtx)

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()
			// 执行具体的业务逻辑
			fun(c)

			finish <- struct{}{}
		}()
		// 执行业务逻辑后操作
		select {
		case p := <-panicChan:
			log.Println(p)
			c.GetResponse().WriteHeader(500)
		case <-finish:
			fmt.Println("finish")
		case <-durationCtx.Done():
			c.SetHasTimeout()
			c.GetResponse().Write([]byte("time out"))
		}
		return nil
	}
}
