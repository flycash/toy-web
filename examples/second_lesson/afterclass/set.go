package afterclass

type Set interface {
	Put(key string)
	Keys() []string
	Contains(key string) bool
	Remove(key string)
	// 如果之前已经有了，就返回旧的值，absent =false
	// 如果之前没有，就塞下去，返回 absent = true
	PutIfAbsent(key string) (old string, absent bool)
}
