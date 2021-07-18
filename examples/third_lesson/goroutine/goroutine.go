package main

import (
	"fmt"
	"time"
)

func main() {
	GoRoutine()
}

func GoRoutine() {
	go func() {
		time.Sleep(10 * time.Second)
	}()
	// 这里直接输出，不会等待十秒
	fmt.Println("I am here")
}