package webv3

import (
	"errors"
	"net/http"
	"sort"
	"strings"
)

var ErrorInvalidRouterPattern = errors.New("invalid router pattern")
var ErrorInvalidMethod = errors.New("invalid method")

type HandlerBasedOnTree struct {
	forest map[string]*node
}

var supportMethods = [4]string {
	http.MethodGet,
	http.MethodPost,
	http.MethodPut,
	http.MethodDelete,
}

func NewHandlerBasedOnTree() Handler {
	forest := make(map[string]*node, len(supportMethods))
	for _, m :=range supportMethods {
		forest[m] = newRootNode(m)
	}
	return &HandlerBasedOnTree{
		forest: forest,
	}
}

// ServeHTTP 就是从树里面找节点
// 找到了就执行
func (h *HandlerBasedOnTree) ServeHTTP(c *Context) {
	handler, found := h.findRouter(c.R.Method, c.R.URL.Path, c)
	if !found {
		c.W.WriteHeader(http.StatusNotFound)
		_, _ = c.W.Write([]byte("Not Found"))
		return
	}
	handler(c)
}

func (h *HandlerBasedOnTree) findRouter(method, path string, c *Context) (handlerFunc, bool) {
	// 去除头尾可能有的/，然后按照/切割成段
	paths := strings.Split(strings.Trim(path, "/"), "/")
	cur, ok := h.forest[method]
	if !ok {
		return nil, false
	}
	for _, p := range paths {
		// 从子节点里边找一个匹配到了当前 p 的节点
		matchChild, found := h.findMatchChild(cur, p, c)
		if !found {
			return nil, false
		}
		cur = matchChild
	}
	// 到这里，应该是找完了
	if cur.handler == nil {
		// 到达这里是因为这种场景
		// 比如说你注册了 /user/profile
		// 然后你访问 /user
		return nil, false
	}
	return cur.handler, true
}

// Route 就相当于往树里面插入节点
func (h *HandlerBasedOnTree) Route(method string, pattern string,
	handlerFunc handlerFunc) error {

	err := h.validatePattern(pattern)
	if err != nil {
		return err
	}

	// 将pattern按照URL的分隔符切割
	// 例如，/user/friends 将变成 [user, friends]
	// 将前后的/去掉，统一格式
	pattern = strings.Trim(pattern, "/")
	paths := strings.Split(pattern, "/")

	// 当前指向根节点
	cur, ok := h.forest[method]
	if !ok {
		return ErrorInvalidMethod
	}
	for index, path := range paths {

		// 从子节点里边找一个匹配到了当前 path 的节点
		matchChild, found := h.findMatchChild(cur, path, nil)
		// != nodeTypeAny 是考虑到 /order/* 和 /order/:id 这种注册顺序
		if found && matchChild.nodeType != nodeTypeAny {
			cur = matchChild
		} else {
			// 为当前节点根据
			h.createSubTree(cur, paths[index:], handlerFunc)
			return nil
		}
	}
	// 离开了循环，说明我们加入的是短路径，
	// 比如说我们先加入了 /order/detail
	// 再加入/order，那么会走到这里
	cur.handler = handlerFunc
	return nil
}

func (h *HandlerBasedOnTree) validatePattern(pattern string) error {
	// 校验 *，如果存在，必须在最后一个，并且它前面必须是/
	// 即我们只接受 /* 的存在，abc*这种是非法

	pos := strings.Index(pattern, "*")
	// 找到了 *
	if pos > 0 {
		// 必须是最后一个
		if pos != len(pattern) - 1 {
			return ErrorInvalidRouterPattern
		}
		if pattern[pos-1] != '/' {
			return ErrorInvalidRouterPattern
		}
	}
	return nil
}

func (h *HandlerBasedOnTree) findMatchChild(root *node,
	path string, c *Context) (*node, bool) {
	candidates := make([]*node, 0, 2)
	for _, child := range root.children {
		if child.matchFunc(path, c) {
			candidates = append(candidates, child)
		}
	}

	if len(candidates) == 0 {
		return nil, false
	}

	// type 也决定了它们的优先级
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].nodeType < candidates[j].nodeType
	})
	return candidates[len(candidates) - 1], true
}

func (h *HandlerBasedOnTree) createSubTree(root *node, paths []string, handlerFn handlerFunc) {
	cur := root
	for _, path := range paths {
		nn := newNode(path)
		cur.children = append(cur.children, nn)
		cur = nn
	}
	cur.handler = handlerFn
}

