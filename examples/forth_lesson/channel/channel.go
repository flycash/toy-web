package main

import (
	"fmt"
	"time"
)

func main() {
	channelWithoutCache()
	channelWithCache()
}

func channelWithCache()  {
	ch := make(chan string, 1)
	go func() {

		ch <- "Hello, first msg from channel"
		time.Sleep(time.Second)
		ch <- "Hello, second msg from channel"
	}()

	time.Sleep(2 * time.Second)
	msg := <- ch
	fmt.Println(time.Now().String() + msg)
	msg = <- ch
	fmt.Println(time.Now().String() + msg)
	// 因为前面我们先睡了2秒，所以其实会有一个已经在缓冲了
	// 当我们尝试输出的时候，这个输出间隔就会明显小于1秒
	// 我电脑上的几次实验，差距都在1ms以内
}

func channelWithoutCache() {
	// 不带缓冲
	ch := make(chan string)
	go func() {
		time.Sleep(time.Second)
		ch <- "Hello, msg from channel"
	}()

	// 这里比较容易写成 msg <- ch，编译会报错
	msg := <- ch
	fmt.Println(msg)
}
