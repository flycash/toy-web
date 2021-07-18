package webv2

import (
	"errors"
	"net/http"
	"strings"
)

var ErrorInvalidRouterPattern = errors.New("invalid router pattern")

type HandlerBasedOnTree struct {
	root *node

}
var supportMethods = [4]string {http.MethodPost, http.MethodGet,
	http.MethodDelete, http.MethodPut}
func NewHandlerBasedOnTree() Handler {
	root := &node{
	}
	return &HandlerBasedOnTree{
		root: root,
	}
}

// ServeHTTP 就是从树里面找节点
// 找到了就执行
func (h *HandlerBasedOnTree) ServeHTTP(c *Context) {
	handler, found := h.findRouter(c.R.URL.Path)
	if !found {
		c.W.WriteHeader(http.StatusNotFound)
		_, _ = c.W.Write([]byte("Not Found"))
		return
	}
	handler(c)
}

// 这个是不好测试的版本，可以尝试为这个写单元测试
// 会发现很难构造 request，也很难对 ResponseWriter做断言
//func (h *HandlerBasedOnTree) ServeHTTP(c *Context) {
//	url := strings.Trim(c.R.URL.Path, "/")
//	paths := strings.Split(url, "/")
//	cur := h.root
//	for _, path := range paths {
//		// 从子节点里边找一个匹配到了当前 path 的节点
//		matchChild, found := h.findMatchChild(cur, path)
//		if !found {
//			// 找不到匹配的路径，直接返回
//			c.W.WriteHeader(404)
//			_, _ = c.W.Write([]byte("Not Found"))
//			return
//		}
//		cur = matchChild
//	}
//	// 到这里，应该是找完了
//	if cur.handler == nil {
//		// 到达这里是因为这种场景
//		// 比如说你注册了 /user/profile
//		// 然后你访问 /user
//		c.W.WriteHeader(404)
//		_, _ = c.W.Write([]byte("Not Found"))
//		return
//	}
//	cur.handler(c)
//}

func (h *HandlerBasedOnTree) findRouter(path string) (handlerFunc, bool) {
	// 去除头尾可能有的/，然后按照/切割成段
	paths := strings.Split(strings.Trim(path, "/"), "/")
	cur := h.root
	for _, p := range paths {
		// 从子节点里边找一个匹配到了当前 path 的节点
		matchChild, found := h.findMatchChild(cur, p)
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
	cur := h.root
	for index, path := range paths {

		// 从子节点里边找一个匹配到了当前 path 的节点
		matchChild, found := h.findMatchChild(cur, path)
		if found {
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

func (h *HandlerBasedOnTree) findMatchChild(root *node, path string) (*node, bool) {
	var wildcardNode *node
	for _, child := range root.children {
		// 并不是 * 的节点命中了，直接返回
		// != * 是为了防止用户乱输入
		if child.path == path &&
			child.path != "*"{
			return child, true
		}
		// 命中了通配符的，我们看看后面还有没有更加详细的
		if child.path == "*" {
			wildcardNode = child
		}
	}
	return wildcardNode, wildcardNode != nil
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

type node struct {
	path string
	children []*node

	// 如果这是叶子节点，
	// 那么匹配上之后就可以调用该方法
	handler handlerFunc
}

func newNode(path string) *node {
	return &node{
		path: path,
		children: make([]*node, 0, 2),
	}
}
