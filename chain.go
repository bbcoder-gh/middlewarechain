// middlewarechain chains different middlewares to a handler which can be used across a variety of routers
package middlewarechain

import "net/http"

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Chain(h http.HandlerFunc, middlewares ...Middleware) (aggregateHandler http.HandlerFunc) {
	aggregateHandler = h

	for _, m := range middlewares {
		aggregateHandler = m(aggregateHandler)
	}
	return
}
