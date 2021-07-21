package main

func init() {
	// 因为我们不能确定 init 方法的执行顺序，
	// 只能曲线救国
	initBeforeSomething()
	initSomething()
	initAfterSomething()
}

func initBeforeSomething()  {
	
}

func initSomething()  {
	
}

func initAfterSomething()  {
	
}