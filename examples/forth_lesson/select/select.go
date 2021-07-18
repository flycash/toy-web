package main

import (
	"fmt"
	"time"
)

func main() {
	// 这个不能在 main 函数运行，是因为运行起来，
	// 所有的goroutine都被我们搞sleep了，直接就崩了
	//Select()
}

func Select() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(time.Second)
		ch1 <- "msg from channel1"
	}()

	go func() {
		time.Sleep(time.Second)
		ch2 <- "msg from channel2"
	}()

	for {
		select {
		case msg := <- ch1:
			fmt.Println(msg)
		case msg := <- ch2:
			fmt.Println(msg)
		}
	}
}
