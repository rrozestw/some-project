package main

import (
	"github.com/gocraft/web"
	log "logex-gls"
	"net/http/httputil"
)

func (c *Context) IncomingRequestLoggingWithBody(rw web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	dump, err := httputil.DumpRequest(r.Request, true)
	if err != nil {
		log.Warn(err)
	}
	log.Debug(string(dump))
	next(rw, r)
}

func (c *Context) IncomingRequestLoggingWithoutBody(rw web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	dump, err := httputil.DumpRequest(r.Request, false)
	if err != nil {
		log.Warn(err)
	}
	log.Debug(string(dump))
	next(rw, r)
}
