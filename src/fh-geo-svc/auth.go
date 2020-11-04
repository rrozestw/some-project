package main

import (
	"strconv"
	//"net/http"
	"github.com/gocraft/web"
	_ "github.com/surma/stacksignal"
	"github.com/tylerb/gls"
	log "logex-gls"
)

// As this is designed to be high perforamance non-user facing service,
// using a horizonatlly scalable cryptography-based authentication (like RSA-signed JWT tokens)
// seems to be more reasonable to database lookups. Revocation is much less of an issue for service authentication,
// and we can make those reasonably short-lived.
//
// ****** FIXME TODO - no authentication is acutally implemented ********************
//
func (c *Context) UserRequired(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	log.Debug("User required middleware called", req.URL.Path)
	if req.URL.Path == "/healthz" ||
		req.URL.Path == "/favicon.ico" {
		next(rw, req)
	} else {
		/*
			auth := req.Header.Get("Authorization")
			if(len(auth) == 0) {
				log.Info("Authentication failed.")
				rw.WriteHeader(http.StatusUnauthorized)
				return
			}
		*/
		c.UserId = 1 // TODO extract from token, ensure token signature is verified, ensure that we fail on alg=none etc. (https://auth0.com/blog/critical-vulnerabilities-in-json-web-token-libraries/).
		gls.Set(log.USER_ID, strconv.FormatInt(c.UserId, 10))
		next(rw, req)
	}
}
