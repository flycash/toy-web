package main

import "fmt"

func main() {
	// 直接初始化一个三个元素的数组。大括号里面多一个或者少一个都编译不通过
	a1 := [3]int{9, 8, 7}
	fmt.Printf("a1: %v, len: %d, cap: %d", a1, len(a1), cap(a1))

	// 初始化一个三个元素的数组，所有元素都是0
	var a2 [3]int
	fmt.Printf("a2: %v, len: %d, cap: %d", a2, len(a2), cap(a2))

	//a1 = append(a1, 12) 数组不支持 append 操作

	// 按下标索引
	fmt.Printf("a1[1]: %d", a1[1])
	// 超出下标范围，直接崩溃，编译不通过
	//fmt.Printf("a1[99]: %d", a1[99])
}