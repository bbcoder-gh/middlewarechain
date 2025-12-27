// middlewarechain chains different middlewares to a handler which can be used across a variety of routers
package middlewarechain

import "net/http"

// Middleware defines a function to process middleware
type Middleware func(http.HandlerFunc) http.HandlerFunc

// Chain applies multiple middlewares to a http.HandlerFunc and returns the final http.HandlerFunc
func Chain(h http.HandlerFunc, middlewares ...Middleware) (aggregateHandler http.HandlerFunc) {
	aggregateHandler = h

	for i := len(middlewares) - 1; i >= 0; i-- {
		aggregateHandler = middlewares[i](aggregateHandler)
	}
	return
}
