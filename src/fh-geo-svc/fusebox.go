package main

import (
	"github.com/gocraft/web"
	//log "logex-gls"
)

// TODO
func (c *Context) CircuitBreakerMiddleware(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	next(rw, req)
}
