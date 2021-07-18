package web

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Hook 是一个钩子函数。注意，
// ctx 是一个有超时机制的 context.Context
// 所以你必须处理超时的问题
type Hook func(ctx context.Context) error

// BuildCloseServerHook 这里其实可以考虑使用 errgroup，
// 但是我们这里不用是希望每个 server 单独关闭
// 互相之间不影响
func BuildCloseServerHook(servers ...Server) Hook {
	return func(ctx context.Context) error {
		wg := sync.WaitGroup{}
		doneCh := make(chan struct{})
		wg.Add(len(servers))

		for _, s := range servers {
			go func(svr Server) {
				err := svr.Shutdown(ctx)
				if err != nil {
					fmt.Printf("server shutdown error: %v \n", err)
				}
				time.Sleep(time.Second)
				wg.Done()
			}(s)
		}
		go func() {
			wg.Wait()
			doneCh <- struct{}{}
		}()
		select {
		case <- ctx.Done():
			fmt.Printf("closing servers timeout \n")
			return ErrorHookTimeout
		case <- doneCh:
			fmt.Printf("close all servers \n")
			return nil
		}
	}
}
