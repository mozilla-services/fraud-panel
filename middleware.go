package main

import (
	"net/http"
	"time"

	"go.mozilla.org/fraud-panel/mozlog"
)

// Middleware wraps an HTTP handler with standard functionalities,
// like logging
type Middleware func(http.Handler) http.Handler

//  Run the request through all middlewares
func handleMiddlewares(h http.Handler, adapters ...Middleware) http.Handler {
	// To make the middleware run in the order in which they are specified,
	// we reverse through them in the Middleware function, rather than just
	// ranging over them
	for i := len(adapters) - 1; i >= 0; i-- {
		h = adapters[i](h)
	}
	return h
}

// logRequest is a middleware that writes details about each HTTP request processed
// but the various handlers.
func logRequest() Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
			mozlog.Event("request", mozlog.Fields{
				"method":             r.Method,
				"proto":              r.Proto,
				"remoteAddress":      r.RemoteAddr,
				"remoteAddressChain": "[" + r.Header.Get("X-Forwarded-For") + "]",
				"rid":                getRequestID(r),
				"t":                  time.Now().Sub(getRequestStartTime(r)) / time.Millisecond,
				"ua":                 r.UserAgent(),
				"url":                r.URL.String()})
		})
	}
}

func setResponseHeaders() Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Security-Policy", "default-src 'none'; object-src 'none';")
			w.Header().Add("X-Frame-Options", "SAMEORIGIN")
			w.Header().Add("X-Content-Type-Options", "nosniff")
			w.Header().Add("Strict-Transport-Security", "max-age=31536000;")
			h.ServeHTTP(w, r)
		})
	}
}
