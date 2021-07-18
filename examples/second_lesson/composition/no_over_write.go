package main

import "fmt"

func main() {
	son := Son{
		Parent{},
	}

	son.SayHello()
}

type Parent struct {

}

func (p Parent) SayHello() {
	fmt.Println("I am " + p.Name())
}

func (p Parent) Name() string {
	return "Parent"
}

type Son struct {
	Parent
}

// 定义了自己的 Name() 方法
func (s Son) Name() string {
	return "Son"
}

