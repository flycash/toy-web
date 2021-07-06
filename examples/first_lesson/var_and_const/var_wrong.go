package main

var aa = "hello"
// var aa = "bbb" 这个包已经有一个 a 了，所以再次声明会导致编译
func main() {
	aa := 13 // 虽然包外面已经有一个 aa 了，但是这里从包变成了局部变量
	println(aa)

	var bb = 15
	//var bb = 16 // 重复声明，也会导致编译不通过
	println(bb)

	bb = 17 // OK,没有重复声明，只是赋值了新的值
	// bb := 18 // 不行，因为 := 就是声明并且赋值的简写，相当于重复声明了 bb
}
