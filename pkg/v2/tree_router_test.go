package webv2

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

func TestHandlerBasedOnTree_Route(t *testing.T) {
	handler := NewHandlerBasedOnTree().(*HandlerBasedOnTree)
	assert.NotNil(t, handler.root)

	err := handler.Route(http.MethodPost, "/user", func(c *Context) {})
	assert.Nil(t, err)
	// 开始做断言，这个时候我们应该确认，在根节点之下只有一个user节点
	assert.Equal(t, 1, len(handler.root.children))

	n := handler.root.children[0]
	assert.NotNil(t, n)
	assert.Equal(t, "user", n.path)
	assert.NotNil(t, n.handler)
	assert.Empty(t, n.children)

	// 我们只有
	//  user -> profile
	err = handler.Route(http.MethodPost, "/user/profile", func(c *Context) {})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(n.children))
	profileNode := n.children[0]
	assert.NotNil(t, profileNode)
	assert.Equal(t, "profile", profileNode.path)
	assert.NotNil(t, profileNode.handler)
	assert.Empty(t, profileNode.children)

	// 试试重复
	err = handler.Route(http.MethodPost, "/user", func(c *Context) {})
	assert.Nil(t, err)
	n = handler.root.children[0]
	assert.NotNil(t, n)
	assert.Equal(t, "user", n.path)
	assert.NotNil(t, n.handler)
	// 有profile节点
	assert.Equal(t, 1, len(n.children))

	// 给 user 再加一个节点
	err = handler.Route(http.MethodPost, "/user/home", func(c *Context) {})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(n.children))
	homeNode := n.children[1]
	assert.NotNil(t, homeNode)
	assert.Equal(t, "home", homeNode.path)
	assert.NotNil(t, homeNode.handler)
	assert.Empty(t, homeNode.children)

	// 添加 /order/detail
	err = handler.Route(http.MethodPost, "/order/detail", func(c *Context) {})
	assert.Equal(t, 2, len(handler.root.children))
	orderNode := handler.root.children[1]
	assert.NotNil(t, orderNode)
	assert.Equal(t, "order", orderNode.path)
	// 此刻我们只有/order/detail，但是没有/order
	assert.Nil(t, orderNode.handler)
	assert.Equal(t, 1, len(orderNode.children))

	orderDetailNode := orderNode.children[0]
	assert.NotNil(t, orderDetailNode)
	assert.Empty(t, orderDetailNode.children)
	assert.Equal(t, "detail", orderDetailNode.path)
	assert.NotNil(t, orderDetailNode.handler)

	// 加一个 /order
	err = handler.Route(http.MethodPost, "/order", func(c *Context) {})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(handler.root.children))
	orderNode = handler.root.children[1]
	assert.Equal(t, "order", orderNode.path)
	// 此时我们有了 /order
	assert.NotNil(t, orderNode.handler)

	err = handler.Route(http.MethodPost, "/order/*", func(c *Context) {})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(orderNode.children))
	orderWildcard := orderNode.children[1]
	assert.NotNil(t, orderWildcard)
	assert.NotNil(t, orderWildcard.handler)
	assert.Equal(t, "*", orderWildcard.path)

	err = handler.Route(http.MethodPost, "/order/*/checkout", func(c *Context) {})
	assert.NotNil(t, err)
}

func TestHandlerBasedOnTree_findRouter(t *testing.T) {
	handler := NewHandlerBasedOnTree().(*HandlerBasedOnTree)
	_ = handler.Route(http.MethodPost, "/user", func(c *Context) {})
	fn, found := handler.findRouter("/user")
	assert.True(t, found)
	assert.NotNil(t, fn)
	_, found = handler.findRouter("/user/profile")
	assert.False(t, found)

	_ = handler.Route(http.MethodPost, "/user/profile", func(c *Context) {})
	_, found = handler.findRouter("/user/profile")
	assert.True(t, found)

	_, found = handler.findRouter("/user")
	assert.True(t, found)

	var detailHandler handlerFunc = func(c *Context) {}
	_ = handler.Route(http.MethodPost, "/order/detail", detailHandler)
	_, found = handler.findRouter("/order")
	assert.False(t, found)

	fn, found = handler.findRouter("/order/detail")
	assert.True(t, found)
	assert.True(t, handlerFuncEquals(detailHandler, fn))

	var wildcardHandler handlerFunc = func(c *Context) {}
	_ = handler.Route(http.MethodPost, "/order/*", wildcardHandler)
	_, found = handler.findRouter("/order")
	assert.False(t, found)

	fn, found = handler.findRouter("/order/detail")
	assert.True(t, found)
	assert.True(t, handlerFuncEquals(detailHandler, fn))

	fn, found = handler.findRouter("/order/checkout")
	assert.True(t, found)
	assert.True(t, handlerFuncEquals(wildcardHandler, fn))
}

func handlerFuncEquals(hf1 handlerFunc, hf2 handlerFunc) bool {
	return reflect.ValueOf(hf1).Pointer() == reflect.ValueOf(hf2).Pointer()
}