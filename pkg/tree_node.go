package web

import (
	"strings"
)

const (

	// 根节点，只有根用这个
	nodeTypeRoot = iota

	// *
	nodeTypeAny

	// 路径参数
	nodeTypeParam

	// 正则
	nodeTypeReg

	// 静态，即完全匹配
	nodeTypeStatic
)

const any = "*"

// matchFunc 承担两个职责，一个是判断是否匹配，一个是在匹配之后
// 将必要的数据写入到 Context
// 所谓必要的数据，这里基本上是指路径参数
type matchFunc func(path string, c *Context) bool

type node struct {
	children []*node

	// 如果这是叶子节点，
	// 那么匹配上之后就可以调用该方法
	handler   handlerFunc
	matchFunc matchFunc

	// 原始的 pattern。注意，它不是完整的pattern，
	// 而是匹配到这个节点的pattern
	pattern string
	nodeType int
}

// 静态节点
func newStaticNode(path string) *node {
	return &node{
		children: make([]*node, 0, 2),
		matchFunc: func(p string, c *Context) bool {
			return path == p && p != "*"
		},
		nodeType: nodeTypeStatic,
		pattern:  path,
	}
}


func newRootNode(method string) *node {
	return &node{
		children: make([]*node, 0, 2),
		matchFunc: func( p string, c *Context) bool {
			panic("never call me")
		},
		nodeType: nodeTypeRoot,
		pattern:  method,
	}
}

func newNode(path string) *node {
	if path == "*"{
		return newAnyNode()
	}
	if strings.HasPrefix(path, ":") {
		return newParamNode(path)
	}
	return newStaticNode(path)
}

// 通配符 * 节点
func newAnyNode() *node {
	return &node{
		// 因为我们不允许 * 后面还有节点，所以这里可以不用初始化
		//children: make([]*node, 0, 2),
		matchFunc: func(p string, c *Context) bool {
			return true
		},
		nodeType: nodeTypeAny,
		pattern:  any,
	}
}

// 路径参数节点
func newParamNode(path string) *node {
	paramName := path[1:]
	return &node{
		children: make([]*node, 0, 2),
		matchFunc: func(p string, c *Context) bool {
			if c != nil {
				c.PathParams[paramName] = p
			}
			// 如果自身是一个参数路由，
			// 然后又来一个通配符，我们认为是不匹配的
			return p != any
		},
		nodeType: nodeTypeParam,
		pattern:  path,
	}
}

// 正则节点
//func newRegNode(path string) *node {
//	// 依据你的规则拿到正则表达式
//	return &node{
//		children: make([]*node, 0, 2),
//		matchFunc: func(p string, c *Context) bool {
//			// 怎么写？
//		},
//		nodeType: nodeTypeParam,
//		pattern: path,
//	}
//}

