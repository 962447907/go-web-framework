package gee

import (
	"net/http"
)

// HandlerFunc defines the request handler used by gee
type HandlerFunc func(ctx *Context)

// Engine implement the interface of ServeHTTP
type Engine struct {
	*RouterGroup
	router *router
}

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc // support middleware
	parent      *RouterGroup  // support nesting
	engine      *Engine       // all groups share an Engine instance
}

// New is the constructor of gee.Engine
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	return engine
}

// Group is defined to create a new RouterGroup
// remember all groups share the same Engine instance
func (routerGroup *RouterGroup) Group(prefix string) *RouterGroup {
	newGroup := &RouterGroup{
		prefix: routerGroup.prefix + prefix,
		engine: routerGroup.engine,
	}
	return newGroup
}

func (routerGroup *RouterGroup) addRoute(method string, pattern string, handler HandlerFunc) {
	routerGroup.engine.router.addRoute(method, routerGroup.prefix+pattern, handler)
}

// GET defines the method to add GET request
func (routerGroup *RouterGroup) GET(pattern string, handler HandlerFunc) {
	routerGroup.addRoute(http.MethodGet, pattern, handler)
}

// POST defines the method to add POST request
func (routerGroup *RouterGroup) POST(pattern string, handler HandlerFunc) {
	routerGroup.addRoute(http.MethodPost, pattern, handler)
}

// PUT defines the method to add PUT request
func (routerGroup *RouterGroup) PUT(pattern string, handler HandlerFunc) {
	routerGroup.addRoute(http.MethodPut, pattern, handler)
}

// DELETE defines the method to add DELETE request
func (routerGroup *RouterGroup) DELETE(pattern string, handler HandlerFunc) {
	routerGroup.addRoute(http.MethodDelete, pattern, handler)
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	engine.router.printRoutes()
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}
