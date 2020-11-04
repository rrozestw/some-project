package main

import (
	_ "github.com/davecgh/go-spew/spew"
	"github.com/gocraft/web"
	"math/rand"
	"strconv"

	"github.com/tylerb/gls"
	log "logex-gls"
)

// ##################################################################
// #  Tracking requests IDs in logger with goroutine-local storage. #
// ##################################################################
// In a typical scenario it is very useful to a have request tracking functionality across the whole stack.
// It calls for a logger that supports dumping request id/details on every line we log.
// This is tricky to implement, as requires passing some context-like struct to functions we call.
// It works, but is changing clean otherwise signatures.
// The code below is an alternative, it uses *no warranty* goroutine-local storage with a modified logex library
// to display those.
// WARNING: test with every golang version used in production, as it relies on some internal impl. details.
// WARNING: when this http stack handler creates a new goroutine, it is loosing that tracking data. This is in
// practice way less issue and it seems, and can be mitigated by putting code below into the new implementation
// of a new goroutine.
func (c *Context) GlsRequestIdMiddleware(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	gls.Set(log.RAND_ID, strconv.Itoa(rand.Intn(9999)))
	gls.Set(log.REQUEST_ID, req.Header.Get("X-Request-Id"))
	gls.Set(log.FORWARDED_FOR, req.Header.Get("X-Forwarded-For"))
	gls.Set(log.USER_ID, "?")
	gls.Set(log.REQUEST_PATH, req.URL.Path)
	gls.Set(log.REQUEST_IP_ADDR, req.RemoteAddr)
	defer gls.Cleanup()
	next(rw, req)
}
