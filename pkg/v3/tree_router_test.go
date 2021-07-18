package webv3

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

func TestHandlerBasedOnTree_Route(t *testing.T) {
	handler := NewHandlerBasedOnTree().(*HandlerBasedOnTree)
	// 要确认已经为支持的方法创建了节点
	assert.Equal(t, len(supportMethods), len(handler.forest))

	postNode := handler.forest[http.MethodPost]

	err := handler.Route(http.MethodPost, "/user", func(c *Context) {})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(postNode.children))

	n := postNode.children[0]
	assert.NotNil(t, n)
	assert.Equal(t, "user", n.pattern)
	assert.NotNil(t, n.handler)
	assert.Empty(t, n.children)

	// 我们只有
	//  user -> profile
	err = handler.Route(http.MethodPost, "/user/profile", func(c *Context) {})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(n.children))
	profileNode := n.children[0]
	assert.NotNil(t, profileNode)
	assert.Equal(t, "profile", profileNode.pattern)
	assert.NotNil(t, profileNode.handler)
	assert.Empty(t, profileNode.children)

	// 试试重复
	err = handler.Route(http.MethodPost, "/user", func(c *Context) {})
	assert.Nil(t, err)
	n = postNode.children[0]
	assert.NotNil(t, n)
	assert.Equal(t, "user", n.pattern)
	assert.NotNil(t, n.handler)
	// 有profile节点
	assert.Equal(t, 1, len(n.children))

	// 给 user 再加一个节点
	err = handler.Route(http.MethodPost, "/user/home", func(c *Context) {})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(n.children))
	homeNode := n.children[1]
	assert.NotNil(t, homeNode)
	assert.Equal(t, "home", homeNode.pattern)
	assert.NotNil(t, homeNode.handler)
	assert.Empty(t, homeNode.children)

	// 添加 /order/detail
	err = handler.Route(http.MethodPost, "/order/detail", func(c *Context) {})
	assert.Equal(t, 2, len(postNode.children))
	orderNode := postNode.children[1]
	assert.NotNil(t, orderNode)
	assert.Equal(t, "order", orderNode.pattern)
	// 此刻我们只有/order/detail，但是没有/order
	assert.Nil(t, orderNode.handler)
	assert.Equal(t, 1, len(orderNode.children))

	orderDetailNode := orderNode.children[0]
	assert.NotNil(t, orderDetailNode)
	assert.Empty(t, orderDetailNode.children)
	assert.Equal(t, "detail", orderDetailNode.pattern)
	assert.NotNil(t, orderDetailNode.handler)

	// 加一个 /order
	err = handler.Route(http.MethodPost, "/order", func(c *Context) {})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(postNode.children))
	orderNode = postNode.children[1]
	assert.Equal(t, "order", orderNode.pattern)
	// 此时我们有了 /order
	assert.NotNil(t, orderNode.handler)

	err = handler.Route(http.MethodPost, "/order/*", func(c *Context) {})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(orderNode.children))
	orderWildcard := orderNode.children[1]
	assert.NotNil(t, orderWildcard)
	assert.NotNil(t, orderWildcard.handler)
	assert.Equal(t, "*", orderWildcard.pattern)

	err = handler.Route(http.MethodPost, "/order/*/checkout", func(c *Context) {})
	assert.Equal(t, ErrorInvalidRouterPattern, err)

	err = handler.Route(http.MethodConnect, "/order/checkout", func(c *Context) {})
	assert.Equal(t, ErrorInvalidMethod, err)

	err = handler.Route(http.MethodPost, "/order/:id", func(c *Context){})
	assert.Nil(t, err)
	// 这时候我们有/order/* 和 /order/:id
	// 因为我们并没有认为它们不兼容，而是/order/:id优先
	assert.Equal(t, 3, len(orderNode.children))
	orderParamNode := orderNode.children[2]
	assert.Equal(t, ":id", orderParamNode.pattern)

}

func TestHandlerBasedOnTree_findRouter(t *testing.T) {
	handler := NewHandlerBasedOnTree().(*HandlerBasedOnTree)
	_ = handler.Route(http.MethodPost, "/user", func(c *Context) {})
	ctx := NewContext(nil, nil)
	fn, found := handler.findRouter(http.MethodPost, "/user", ctx)
	assert.True(t, found)
	assert.NotNil(t, fn)
	_, found = handler.findRouter(http.MethodPost,"/user/profile", ctx)
	assert.False(t, found)

	_ = handler.Route(http.MethodPost, "/user/profile", func(c *Context) {})
	_, found = handler.findRouter(http.MethodPost, "/user/profile", ctx)
	assert.True(t, found)

	_, found = handler.findRouter(http.MethodPost, "/user", ctx)
	assert.True(t, found)

	var detailHandler handlerFunc = func(c *Context) {}
	_ = handler.Route(http.MethodPost, "/order/detail", detailHandler)
	_, found = handler.findRouter(http.MethodPost,"/order", ctx)
	assert.False(t, found)

	fn, found = handler.findRouter(http.MethodPost,"/order/detail", ctx)
	assert.True(t, found)
	assert.True(t, handlerFuncEquals(detailHandler, fn))

	var wildcardHandler handlerFunc = func(c *Context) {}
	_ = handler.Route(http.MethodPost, "/order/*", wildcardHandler)
	_, found = handler.findRouter(http.MethodPost,"/order", ctx)
	assert.False(t, found)

	fn, found = handler.findRouter(http.MethodPost,"/order/detail", ctx)
	assert.True(t, found)
	assert.True(t, handlerFuncEquals(detailHandler, fn))

	fn, found = handler.findRouter(http.MethodPost,"/order/checkout", ctx)
	assert.True(t, found)
	assert.True(t, handlerFuncEquals(wildcardHandler, fn))

	_, found = handler.findRouter(http.MethodGet,"/order/checkout", ctx)
	assert.False(t, found)

	// 参数路由
	handler.Route(http.MethodPost, "/order/*", wildcardHandler)
}

func handlerFuncEquals(hf1 handlerFunc, hf2 handlerFunc) bool {
	return reflect.ValueOf(hf1).Pointer() == reflect.ValueOf(hf2).Pointer()
}