// Package router provides convenient wrappers around httprouter to work with
// net/http handlers
package router

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type contextKey string

// String output the details of context key
func (c contextKey) String() string {
	return "router context key " + string(c)
}

// contextKeyParams is the key used for stroing router Param structure
var (
	contextKeyParams = contextKey("params")
)

// Wrapper is a wrapper for httprouter
type Wrapper struct {
	Router *httprouter.Router
}

// NewRouter is the constructor
func NewRouter() *Wrapper {
	return &Wrapper{Router: httprouter.New()}
}

// Delete is a shortcut wrapper for httprouter's DELETE
// Example: r.Delete("/exercise", handler)
func (r *Wrapper) Delete(path string, fn http.HandlerFunc) {
	r.Router.DELETE(path, HandlerFunc(fn))
}

// Get is a shortcut wrapper for httprouter's GET
// Example: r.Get("/exercise", handler)
func (r *Wrapper) Get(path string, fn http.HandlerFunc) {
	r.Router.GET(path, HandlerFunc(fn))
}

// Head is a shortcut wrapper for httprouter's HEAD
// Example: r.HEAD("/exercise", handler)
func (r *Wrapper) Head(path string, fn http.HandlerFunc) {
	r.Router.HEAD(path, HandlerFunc(fn))
}

// Options is a shortcut wrapper for httprouter's OPTIONS
// Example: r.OPTIONS("/exercise", handler)
func (r *Wrapper) Options(path string, fn http.HandlerFunc) {
	r.Router.OPTIONS(path, HandlerFunc(fn))
}

// Patch is a shortcut wrapper for httprouter's PATCH
// Example: r.PATCH("/exercise", handler)
func (r *Wrapper) Patch(path string, fn http.HandlerFunc) {
	r.Router.PATCH(path, HandlerFunc(fn))
}

// Post is a shortcut wrapper for httprouter's POST
// Example: r.POST("/exercise", handler)
func (r *Wrapper) Post(path string, fn http.HandlerFunc) {
	r.Router.POST(path, HandlerFunc(fn))
}

// Put is a shortcut wrapper for httprouter's PUT
// Example: r.PUT("/exercise", handler)
func (r *Wrapper) Put(path string, fn http.HandlerFunc) {
	r.Router.PUT(path, HandlerFunc(fn))
}

// HandlerFunc wraps around standard net/http handler function to work with
// httprouter. Generally this function is not called directly.
func HandlerFunc(fn http.HandlerFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := context.WithValue(r.Context(), contextKeyParams, p)
		fn(w, r.WithContext(ctx))
	}
}

// Params returns the httprouter.Params struct from the http request
func Params(r *http.Request) httprouter.Params {
	return r.Context().Value(contextKeyParams).(httprouter.Params)
}
