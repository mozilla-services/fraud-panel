package main

import "net/http"

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
