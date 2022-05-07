package gee

import (
	"gee/gee/context"
	"gee/gee/middlewares"
	"net/http"
	"strings"
)

// Engine implement the interface of ServeHTTP
type Engine struct {
	*RouterGroup
	router       *router
	routerGroups []*RouterGroup
}

type RouterGroup struct {
	prefix      string                // like /a/b/c
	middlewares []context.HandlerFunc // support middleware
	parent      *RouterGroup          // support nesting
	engine      *Engine               // all groups share an Engine instance
}

// New is the constructor of gee.Engine
func New() *Engine {
	engine := &Engine{
		router: newRouter(),
	}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.routerGroups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// Group is defined to create a new RouterGroup
// remember all groups share the same Engine instance
func (routerGroup *RouterGroup) Group(prefix string) *RouterGroup {
	engine := routerGroup.engine
	newGroup := &RouterGroup{
		prefix: routerGroup.prefix + prefix,
		engine: engine,
		parent: routerGroup,
	}
	engine.routerGroups = append(engine.routerGroups, newGroup)
	return newGroup
}

// Use is defined to add middleware to the group
func (routerGroup *RouterGroup) Use(middlewares ...context.HandlerFunc) {
	routerGroup.middlewares = append(routerGroup.middlewares, middlewares...)
}

func (routerGroup *RouterGroup) addRoute(method string, pattern string, handler context.HandlerFunc) {
	routerGroup.engine.router.addRoute(method, routerGroup.prefix+pattern, handler)
}

// GET defines the method to add GET request
func (routerGroup *RouterGroup) GET(pattern string, handler context.HandlerFunc) {
	routerGroup.addRoute(http.MethodGet, pattern, handler)
}

// POST defines the method to add POST request
func (routerGroup *RouterGroup) POST(pattern string, handler context.HandlerFunc) {
	routerGroup.addRoute(http.MethodPost, pattern, handler)
}

// PUT defines the method to add PUT request
func (routerGroup *RouterGroup) PUT(pattern string, handler context.HandlerFunc) {
	routerGroup.addRoute(http.MethodPut, pattern, handler)
}

// DELETE defines the method to add DELETE request
func (routerGroup *RouterGroup) DELETE(pattern string, handler context.HandlerFunc) {
	routerGroup.addRoute(http.MethodDelete, pattern, handler)
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	engine.router.printRoutes()
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := context.NewContext(w, req)
	var handlerFuncs []context.HandlerFunc
	for _, routerGroup := range engine.routerGroups {
		if strings.HasPrefix(req.URL.Path, routerGroup.prefix) {
			handlerFuncs = append(handlerFuncs, routerGroup.middlewares...)
		}
	}
	c.Handlers = handlerFuncs
	engine.router.handle(c)
}

// Default use Logger() & Recovery middlewares
func Default() *Engine {
	engine := New()
	engine.Use(middlewares.Logger(), middlewares.Recovery())
	return engine
}
