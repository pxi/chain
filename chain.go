// Package chain implements minimal HTTP middleware chaining.
package chain

import "net/http"

// Func is a middleware function. It usually performs some middleware action
// and then passes control to the given handler.
type Func func(http.Handler) http.Handler

// Slice is a collection of Func's that are executed in order.
type Slice []Func

// Then returns an http.Handler that calls each Func in s and finally
// dispatches the request to the given handler.
func (s Slice) Then(h http.Handler) http.Handler {
	max := len(s) - 1
	for i := range s {
		h = s[max-i](h)
	}
	return h
}

// ThenFunc returns an http.Handler that calls each Func in s and finally
// dispatches the request to the given handler function.
func (s Slice) ThenFunc(h http.HandlerFunc) http.Handler { return s.Then(h) }
