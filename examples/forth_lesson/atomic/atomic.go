package main

import "sync/atomic"

var value int32 = 0
func main() {
	// 要传入 value 的指针
	// 把 value + 10
	atomic.AddInt32(&value, 10)
	nv := atomic.LoadInt32(&value)
	// 输出10
	println(nv)
	// 如果之前的值是10，那么就设置为新的值 20
	swapped := atomic.CompareAndSwapInt32(&value, 10, 20)
	// 输出 true
	println(swapped)

	// 如果之前的值是19，那么就设置为新的值 50
	// 显然现在 value 是 20
	swapped = atomic.CompareAndSwapInt32(&value, 19, 50)
	// 输出 false
	println(swapped)

	old := atomic.SwapInt32(&value, 40)
	// 应该是20，即原本的值
	println(old)
	// 输出新的值，也就是交换后的值，40
	println(value)
}