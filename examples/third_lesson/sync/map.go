package main

import (
	"fmt"
	"sync"
)

func main() {
	m := sync.Map{}
	m.Store("cat", "Tom")
	m.Store("mouse", "Jerry")

	// 这里重新读取出来的，就是
	val, ok := m.Load("cat")
	if ok {
		fmt.Println(len(val.(string)))
	}
}