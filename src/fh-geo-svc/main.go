package main

import (
	util "fh-common"
	"fmt"
	"net/http"
	"os"
	//"strconv"

	_ "github.com/davecgh/go-spew/spew"
	"github.com/gocraft/web"

	_ "github.com/surma/stacksignal"
	log "logex-gls"
	"math/rand"
	_ "net/http/pprof"
	"time"
)

type Context struct {
	UserId int64
}

func (c *Context) Favicon(rw web.ResponseWriter, req *web.Request) {
	rw.WriteHeader(http.StatusNotFound)
}

func (c *Context) Index(rw web.ResponseWriter, req *web.Request) {
	log.Info("check")
	fmt.Fprintln(rw, `<!DOCTYPE html><html lang="en"><head><meta charset="utf-8"></head><body><a href="/api/v1/geolocate/:ip">/api/v1/geolocate/:ip</a>`)
}

func (c *Context) Healthz(rw web.ResponseWriter, req *web.Request) {
	log.Info("healthz - ok")
	fmt.Fprintln(rw, "Health OK")
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	configure_from_env()
	log.Info("Configuring web router...")
	go func() {
		if len(os.Getenv("PPROF_LISTEN")) > 0 {
			log.Warn("PPROF enabled, insecure mode, listening on", os.Getenv("PPROF_LISTEN"))
			log.Println(http.ListenAndServe(os.Getenv("PPROF_LISTEN"), nil))
		}
	}()
	
	if len(os.Args) >= 2 && os.Args[1] == "parse" {
		csv_loader()
		return
	}

	r := web.New(Context{})
	r.Middleware(web.LoggerMiddleware)
	r.Middleware((*Context).GlsRequestIdMiddleware)
	r.Middleware((*Context).CircuitBreakerMiddleware)

	if os.Getenv("DEBUG_REQUEST") == "1" {
		log.Warn("DEBUG_REQUEST enabled, insecure mode.")
		r.Middleware((*Context).IncomingRequestLoggingWithoutBody)
	} else if os.Getenv("DEBUG_REQUEST_BODY") == "1" {
		log.Warn("DEBUG_REQUEST_BODY enabled, insecure mode.")
		r.Middleware((*Context).IncomingRequestLoggingWithBody)
	}
	r.Middleware((*Context).UserRequired)
	r.Get("/favicon.ico", (*Context).Favicon)
	r.Get("/", (*Context).Index)
	r.Get("/healthz", (*Context).Healthz)
	r.Get("/api/v1/geolocate/:ip", (*Context).GeolocateEndpoint)

	log.Info("Starting http server...")
	err := http.ListenAndServe(os.Getenv("LISTEN"), r)
	log.Error("HTTP Server has finished with error: ", err)
}

func configure_from_env() {
	log.Info("Configuring...")
	//util.CheckOrSetEnv("PPROF_LISTEN", "127.0.0.1:6061")
	util.CheckOrSetEnv("LISTEN", "0.0.0.0:8080")
	util.CheckOrSetEnv("POSTGRES_DB", "fhgeo")
	util.CheckOrSetEnv("POSTGRES_USER", "fhgeo")
	util.CheckOrSetEnvPassword("POSTGRES_PASS", "1234")
	util.CheckOrSetEnv("POSTGRES_HOST", "localhost")
	util.CheckOrSetEnv("POSTGRES_SSLMODE", "disable")
	util.CheckOrSetEnv("POSTGRES_PORT", "5432")
}
