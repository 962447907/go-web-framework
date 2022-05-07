package middlewares

import (
	"gee/gee/context"
	"log"
	"time"
)

func Logger() context.HandlerFunc {
	return func(c *context.Context) {
		// Start timer
		t := time.Now()
		c.Next()
		// Calculate resolution time
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
