package main

import (
	"fmt"
	"sync"
)

func main() {
	res := 0
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(val int) {
			res += val
			wg.Done()
		}(i)
	}
	// 把这个注释掉你会发现，什么结果你都可能拿到
	wg.Wait()
	fmt.Println(res)
}