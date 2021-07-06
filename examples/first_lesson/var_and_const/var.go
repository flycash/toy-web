package main

// Global 首字母大写，全局可以访问
var Global = "全局变量"

// 首字母小写，只能在这个包里面使用
// 其子包也不能用
var local = "包变量"

var (
	First string = "abc"
	second int32 = 16
)

func main() {
	// int 是灰色的，是因为 golang 自己可以做类型推断，它觉得你可以省略
	var a int = 13
	println(a)

	// 这里我们省略了类型
	var b = 14
	println(b)

	// 这里 uint 不可省略，因为生路之后，因为不加 uint 类型，15会被解释为 int 类型
	var c uint = 15
	println(c)

	// 这一句无法通过编译，因为 golang 是强类型语言，并且不会帮你做任何的转换
	// println(a == c)

	// 只声明不赋值，d 是默认值 0，类型不可以省略
	var d int
	println(d)
}



