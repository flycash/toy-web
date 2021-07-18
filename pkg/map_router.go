package web

import (
	"fmt"
	"net/http"
	"sync"
)

// 一种常用的GO设计模式，
// 用于确保HandlerBasedOnMap肯定实现了这个接口
var _ Handler = &HandlerBasedOnMap{}


type HandlerBasedOnMap struct {
	handlers sync.Map
}

func (h *HandlerBasedOnMap) ServeHTTP(c *Context) {
	request := c.R
	key := h.key(request.Method, request.URL.Path)
	handler, ok := h.handlers.Load(key)
	if !ok {
		c.W.WriteHeader(http.StatusNotFound)
		_, _ = c.W.Write([]byte("not any router match"))
		return
	}

	handler.(handlerFunc)(c)
}

func (h *HandlerBasedOnMap) Route(method string, pattern string,
	handlerFunc handlerFunc) error {
	key := h.key(method, pattern)
	h.handlers.Store(key, handlerFunc)
	return nil
}

func (h *HandlerBasedOnMap) key(method string,
	path string) string {
	return fmt.Sprintf("%s#%s", method, path)
}

func NewHandlerBasedOnMap() *HandlerBasedOnMap {
	return &HandlerBasedOnMap{}
}
