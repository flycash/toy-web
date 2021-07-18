package webv2

import (
	"net/http"
)

// Routable 可路由的
type Routable interface {
	// Route 设定一个路由，命中该路由的会执行handlerFunc的代码
	Route(method string, pattern string, handlerFunc handlerFunc) error
}

// Server 是http server 的顶级抽象
type Server interface {
	Routable
	// Start 启动我们的服务器
	Start(address string) error
}

// sdkHttpServer 这个是基于 net/http 这个包实现的 http server
type sdkHttpServer struct {
	// Name server 的名字，给个标记，日志输出的时候用得上
	Name    string
	handler Handler
	root    Filter
}

func (s *sdkHttpServer) Route(method string, pattern string,
	handlerFunc handlerFunc) error {
	return s.handler.Route(method, pattern, handlerFunc)
}

func (s *sdkHttpServer) Start(address string) error {
	return http.ListenAndServe(address, s)
}

func (s *sdkHttpServer) ServeHTTP(writer http.ResponseWriter, request *http.Request)  {
	c := NewContext(writer, request)
	s.root(c)
}

func NewSdkHttpServer(name string, builders ...FilterBuilder) Server {

	// 改用我们的树
	handler := NewHandlerBasedOnTree()
	//handler := NewHandlerBasedOnMap()
	// 因为我们是一个链，所以我们把最后的业务逻辑处理，也作为一环
	var root Filter = handler.ServeHTTP
	// 从后往前把filter串起来
	for i := len(builders) - 1; i >= 0; i-- {
		b := builders[i]
		root = b(root)
	}
	res := &sdkHttpServer{
		Name: name,
		handler: handler,
		root: root,
	}
	return res
}

