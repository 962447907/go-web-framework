package main

import (
	"gee/gee"
	"gee/gee/context"
	"gee/gee/middlewares"
	"log"
	"net/http"
)

func main() {
	// r := gee.New()
	// r.GET("/", func(c *gee.Context) {
	// 	c.HTML(http.StatusOK, "/")
	// })
	// r.GET("/a/c", func(c *gee.Context) {
	// 	// expect /hello?name=geektutu
	// 	c.String(http.StatusOK, "/a/c")
	// })
	//
	// r.GET("/a/c", func(c *gee.Context) {
	// 	// expect /hello?name=geektutu
	// 	c.String(http.StatusOK, "/a/d")
	// })
	//
	// r.GET("/a/b", func(c *gee.Context) {
	// 	// expect /hello?name=geektutu
	// 	c.String(http.StatusOK, "/a/b")
	// })
	// r.GET("/a/a", func(c *gee.Context) {
	// 	// expect /hello?name=geektutu
	// 	c.String(http.StatusOK, "/a/a")
	// })
	// r.GET("/a/e", func(c *gee.Context) {
	// 	// expect /hello?name=geektutu
	// 	c.String(http.StatusOK, "/a/e")
	// })
	//
	// r.GET("/a/b/c", func(c *gee.Context) {
	// 	// expect /hello/geektutu
	// 	c.String(http.StatusOK, "/a/b/c")
	// })
	r := gee.New()
	r.GET("/index", func(c *context.Context) {
		c.HTML(http.StatusOK, "<h1>Index Page</h1>")
	})
	r.Use(middlewares.Logger(), middlewares.Recovery())
	// index out of range for testing Recovery()
	r.GET("/panic", func(c *context.Context) {
		names := []string{"panic"}
		c.String(http.StatusOK, names[100])
	})
	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *context.Context) {
			c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
		})

		v1.GET("/hello", func(c *context.Context) {
			// expect /hello?name=geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}
	v2 := r.Group("/v2")
	{
		v2.POST("/login", func(c *context.Context) {
			c.JSON(http.StatusOK, context.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})

	}
	log.Panicln(r.Run(":9998"))
}
