package middleware

import "net/http"

// Middleware is a function that wraps an http.Handler with additional behavior
type Middleware func(http.Handler) http.Handler

// ChainMiddlewares applies multiple middlewares to a handler in order.
// The first middleware in the list runs first (outermost).
func ChainMiddlewares(h http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}
