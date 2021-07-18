package web

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"time"
)

var ErrorHookTimeout = errors.New("the hook timeout")

type GracefulShutdown struct {
	// 还在处理中的请求数
	reqCnt int64
	// 大于 1 就说明要关闭了
	closing int32
	
	// 用 channel 来通知已经处理完了所有请求
	zeroReqCnt chan struct{}
}

func NewGracefulShutdown() *GracefulShutdown {
	return &GracefulShutdown{
		zeroReqCnt: make(chan struct{}),
	}
}

// ShutdownFilterBuilder 这个东西怎么保持线程安全呢？
// 它的逻辑有点绕，核心就在于当我们准备关闭的时候，这个动作是单向的，就是说，我的closing一旦加1
// 就再也不会-1
// 所以我们不需要用一个锁把整个方法锁住
// 而实际上，基于这个理由，我们也不需要把 closing 声明为 int32
// 只需要声明 bool，然后在关闭的时候设置为 true。在这里直接检测 true or false就可以。
// 这种做法有一个很重要的点是，在设置值的时候，即便 bool 被高速缓存缓存了，
// 即便了 bool 在平台上，处理器并不能一条指令 设置好值，
// 但是也没什么关系。因为我们可以确认，最终 bool 会变为 true
// 这个做法更加难以理解，所以采用了使用 closing int32 的做法
func (g *GracefulShutdown) ShutdownFilterBuilder(next Filter) Filter {
	return func(c *Context) {
		// 开始拒绝所有的请求
		cl :=  atomic.LoadInt32(&g.closing)
		if cl > 0 {
			c.W.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		atomic.AddInt64(&g.reqCnt, 1)
		next(c)
		n := atomic.AddInt64(&g.reqCnt, -1)
		// 已经开始关闭了，而且请求数为0，
		if cl > 0 && n == 0 {
			g.zeroReqCnt <- struct{}{}
		}
	}
}

// RejectNewRequestAndWaiting 将会拒绝新的请求，并且等待处理中的请求
func (g *GracefulShutdown) RejectNewRequestAndWaiting(ctx context.Context) error {

	atomic.AddInt32(&g.closing, 1)

	// 特殊 case 关闭之前其实就已经处理完了请求。
	if atomic.LoadInt64(&g.reqCnt) == 0 {
		return nil
	}
	
	done := ctx.Done()
	// 因为是单向的，所以我们这里不用 for 在外面包
	// 所谓单向就是，我一触发就回不到原来正常处理请求的状态了
	// 这个 select 可以理解为，要么超时了
	// 要么我这里所有的请求都执行完了
	select {
	case <- done:
		fmt.Println("超时了，还没等到所有请求执行完毕")
		return ErrorHookTimeout
	case <- g.zeroReqCnt:
		fmt.Println("全部请求处理完了")
	}
	return nil
}

func WaitForShutdown(hooks...Hook) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, ShutdownSignals...)
	select {
	case sig := <-signals:
		fmt.Printf("get signal %s, application will shutdown \n", sig)
		// 十分钟都还不行，就直接强退了
		time.AfterFunc(time.Minute * 10, func() {
			fmt.Printf("Shutdown gracefully timeout, application will shutdown immediately. ")
			os.Exit(1)
		})
		for _, h := range hooks {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second * 30)
			err := h(ctx)
			if err != nil {
				fmt.Printf("failed to run hook, err: %v \n", err)
			}
			cancel()
		}
		os.Exit(0)
	}
}
