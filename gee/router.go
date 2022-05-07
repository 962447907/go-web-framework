package gee

import (
	"errors"
	"fmt"
	"gee/gee/context"
	"log"
	"net/http"
	"sort"
	"strings"
)

type node struct {
	path     string           // 路由路径
	part     string           // 路由中由'/'分隔的部分
	children map[string]*node // 子节点
	isWild   bool             // 是否是通配符节点
}

var (
	methods = []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}
)

func (n *node) String() string {
	return fmt.Sprintf("node{part=%s, part=%s, isWild=%t}", n.part, n.part, n.isWild)
}

func (n *node) travel(nodes *[]*node) {
	if n.path != "" {
		*nodes = append(*nodes, n)
	}
	keys := n.sort()
	for _, key := range keys {
		n.children[key].travel(nodes)
	}
}

func (n *node) sort() []string {
	// 得到各个key
	var keys []string
	for key := range n.children {
		keys = append(keys, key)
	}
	// 给key排序，从小到大
	sort.Sort(sort.StringSlice(keys))
	return keys
}

type router struct {
	roots    map[string]*node
	handlers map[string]context.HandlerFunc
}

func newRouter() *router {
	r := &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]context.HandlerFunc),
	}
	for _, method := range methods {
		r.roots[method] = &node{children: make(map[string]*node)}
	}
	return r
}

// parsePath Only one * is allowed
func parsePath(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

// addRoute 绑定路由到handler
func (r *router) addRoute(method, path string, handler context.HandlerFunc) {
	parts := parsePath(path)
	root := r.roots[method]
	if root.path != "" && len(parts) == 0 {
		panic(errors.New(fmt.Sprintf("duplicate route declaration:%4s - %s", method, root.path)))
	}
	key := r.getRouteKey(method, path)
	// 将parts插入到路由树
	for _, part := range parts {
		if root.children[part] == nil {
			root.children[part] = &node{
				part:     part,
				children: make(map[string]*node, 0),
				isWild:   part[0] == ':' || part[0] == '*'}
		}

		root = root.children[part]
	}
	if root.path != "" {
		panic(errors.New(fmt.Sprintf("duplicate route declaration:%4s - %s", method, root.path)))
	}
	// 相当于前缀树的Stop,标识路由结束
	root.path = path
	// 绑定路由和handler
	r.handlers[key] = handler
}

// getRoute 获取路由树节点以及路由变量
// method用来判断属于哪一个方法路由树，path用来获取路由树节点和参数
func (r *router) getRoute(method, path string) (node *node, params map[string]string) {
	params = map[string]string{}
	searchParts := parsePath(path)

	node = r.roots[method]
	// 如果根节点没有子节点并且根节点没有对应路由
	if len(node.children) == 0 && node.path != "/" {
		return nil, nil
	}
	// 在该方法的路由树上查找该路径
	for _, part := range searchParts {
		node = node.children[part]
		if node == nil {
			return nil, nil
		}
	}

	return node, params

}

// handle 用来绑定路由和handlerFunc
func (r *router) handle(c *context.Context) {
	// 获取路由树节点和动态路由中的参数
	node, params := r.getRoute(c.Method, c.Path)
	if node != nil {
		c.Params = params
		key := r.getRouteKey(c.Method, node.path)
		c.Handlers = append(c.Handlers, r.handlers[key])
	} else {
		c.Handlers = []context.HandlerFunc{func(ctx *context.Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND %s \n", c.Path)
		}}
	}
	c.Next()
}

func (r *router) getRouteKey(method string, path string) string {
	key := method + "-" + path
	return key
}

func (r *router) getRoutes(method string) []*node {
	root := r.roots[method]
	nodes := make([]*node, 0)
	root.travel(&nodes)
	return nodes
}

func (r *router) printRoutes() {
	for _, method := range methods {

		nodes := r.getRoutes(method)
		for _, node := range nodes {
			log.Printf("Route %4s - %s", method, node.path)
		}
	}
}
