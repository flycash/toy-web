package main

import (
	"sync"
)

var mutex sync.Mutex
var rwMutex sync.RWMutex
func Mutex() {
	mutex.Lock()
	defer mutex.Unlock()
	// 你的代码
}

func RwMutex()  {
	// 加读锁
	rwMutex.RLock()
	defer rwMutex.RUnlock()

	// 也可以加写锁
	rwMutex.Lock()
	defer rwMutex.Unlock()
}

// 不可重入例子
func Failed1()  {
	mutex.Lock()
	defer mutex.Unlock()

	// 这一句会死锁
	// 但是如果你只有一个goroutine，那么这一个会导致程序崩溃
	mutex.Lock()
	defer mutex.Unlock()
}

// 不可升级
func Failed2()  {
	rwMutex.RLock()
	defer rwMutex.RUnlock()

	// 这一句会死锁
	// 但是如果你只有一个goroutine，那么这一个会导致程序崩溃
	mutex.Lock()
	defer mutex.Unlock()
}
