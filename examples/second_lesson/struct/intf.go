package main

// 首字母小写，所以是一个包私有的接口
type animal interface {
	// 这里可以有任意多个方法，不过我们一般建议是小接口，
	// 即接口里面不会有很多方法
	// 方法声明不需要 func 关键字

	Eat()
}

// 首字母大写，所以是一个包外可访问的接口
type Duck interface {
	Swim()
}
