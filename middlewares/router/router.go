// Package router provides convenient wrappers around httprouter to work with
// net/http handlers
package router

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const params = "params"

// Wrapper for httprouter
type RouterWrapper struct {
	Router *httprouter.Router
}

// Creates a new instance of wrapper
func NewRouter() *RouterWrapper {
	return &RouterWrapper{Router: httprouter.New()}
}

// Delete is a shortcut wrapper for httprouter's DELETE
// Example: r.Delete("/exercise", handler)
func (r *RouterWrapper) Delete(path string, fn http.HandlerFunc) {
	r.Router.DELETE(path, HandlerFunc(fn))
}

// Get is a shortcut wrapper for httprouter's GET
// Example: r.Get("/exercise", handler)
func (r *RouterWrapper) Get(path string, fn http.HandlerFunc) {
	r.Router.GET(path, HandlerFunc(fn))
}

// Get is a shortcut wrapper for httprouter's HEAD
// Example: r.HEAD("/exercise", handler)
func (r *RouterWrapper) Head(path string, fn http.HandlerFunc) {
	r.Router.HEAD(path, HandlerFunc(fn))
}

// Get is a shortcut wrapper for httprouter's OPTIONS
// Example: r.OPTIONS("/exercise", handler)
func (r *RouterWrapper) Options(path string, fn http.HandlerFunc) {
	r.Router.OPTIONS(path, HandlerFunc(fn))
}

// Get is a shortcut wrapper for httprouter's PATCH
// Example: r.PATCH("/exercise", handler)
func (r *RouterWrapper) Patch(path string, fn http.HandlerFunc) {
	r.Router.PATCH(path, HandlerFunc(fn))
}

// Get is a shortcut wrapper for httprouter's POST
// Example: r.POST("/exercise", handler)
func (r *RouterWrapper) Post(path string, fn http.HandlerFunc) {
	r.Router.POST(path, HandlerFunc(fn))
}

// Get is a shortcut wrapper for httprouter's PUT
// Example: r.PUT("/exercise", handler)
func (r *RouterWrapper) Put(path string, fn http.HandlerFunc) {
	r.Router.PUT(path, HandlerFunc(fn))
}

// HandlerFunc wraps around standard net/http handler function to work with
// httprouter. Generally this function is not called directly.
func HandlerFunc(fn http.HandlerFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := context.WithValue(r.Context(), params, p)
		fn(w, r.WithContext(ctx))
	}
}

// Return the httprouter.Params struct from the http request
func Params(r *http.Request) httprouter.Params {
	return r.Context().Value(params).(httprouter.Params)
}
