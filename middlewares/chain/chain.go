// Package chain provides types and functions to compose net/http compatible middlewares.
// It is completely borrowed from alice(https://godoc.org/github.com/justinas/alice) except it only
// works with http.HandlerFunc.
package chain

import (
	"net/http"
)

// MiddlewareFn is a signature type for net/http middleware
type MiddlewareFn func(http.HandlerFunc) http.HandlerFunc

// Chain keeps a stack of net/http middlewares
type Chain struct {
	middlewares []MiddlewareFn
}

// NewChain is a constructor for a new chain
func NewChain(m ...MiddlewareFn) Chain {
	return Chain{append([]MiddlewareFn(nil), m...)}
}

// ThenFunc terminates a middleware chain by passing the final http.HandlerFunc
func (c Chain) ThenFunc(fn http.HandlerFunc) http.HandlerFunc {
	idx := len(c.middlewares) - 1
	for i := idx; i >= 0; i-- {
		fn = c.middlewares[i](fn)
	}
	return fn
}

// Append takes a new http.HandlerFunc, appends and returns a new chain from the existing one
func (c Chain) Append(m ...MiddlewareFn) Chain {
	nc := make([]MiddlewareFn, 0, len(c.middlewares)+len(m))
	nc = append(nc, c.middlewares...)
	return Chain{append(nc, m...)}
}

// Extend takes a new Chain, appends and returns a new chain from the existing one
func (c Chain) Extend(chain Chain) Chain {
	return c.Append(chain.middlewares...)
}
