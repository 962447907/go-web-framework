package gee

import (
	"testing"
)

func TestRouter(t *testing.T) {
	r := newRouter()
	// r.addRoute("GET", "/", nil)
	r.addRoute("GET", "/a/b/", nil)
	r.addRoute("GET", "/", nil)
	r.addRoute("GET", "/a/b/c", nil)
	r.addRoute("GET", "/hello/:name", nil)
	r.addRoute("GET", "/hi/:name", nil)
	r.addRoute("GET", "/assets/*filepath", nil)
}
