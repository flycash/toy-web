package main

import "fmt"

func main() {
	// 指针用 * 表示
	var p *ToyDuck = &ToyDuck{}
	// 解引用，得到结构体
	var duck ToyDuck = *p
	duck.Swim()

	// 只是声明了，但是没有使用
	var nilDuck *ToyDuck
	if nilDuck == nil {
		fmt.Println("nilDuck is nil")
	}
}