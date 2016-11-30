// Package cors is a thin wrapper for original
// cors(https://godoc.org/github.com/rs/cors) package which provides a
// http.HandlerFunc compatible function for middleware chaining.
package cors

import (
	"net/http"

	"github.com/dictyBase/go-middlewares/middlewares/chain"
	"github.com/rs/cors"
)

func CorsAdapter(c *cors.Cors) chain.MiddlewareFn {
	fn := func(fn http.HandlerFunc) http.HandlerFunc {
		ifn := func(w http.ResponseWriter, r *http.Request) {
			c.HandlerFunc(w, r)
			fn(w, r)
		}
		return ifn
	}
	return fn
}
