package nocache

import (
	"net/http"
	"time"
)

// Unix epoch time
var epoch = time.Unix(0, 0).Format(time.RFC1123)

// Taken from https://github.com/mytrile/nocache
var noCacheHeaders = map[string]string{
	"Expires":         epoch,
	"Cache-Control":   "no-cache, no-store, no-transform, must-revalidate, private, max-age=0",
	"Pragma":          "no-cache",
	"X-Accel-Expires": "0",
}

// Handler as per http://wiki.nginx.org/HttpProxyModule. It sets:
//      Expires: Thu, 01 Jan 1970 00:00:00 UTC
//      Cache-Control: no-cache, private, max-age=0
//      X-Accel-Expires: 0
//      Pragma: no-cache (for HTTP/1.0 proxies/clients)
func Hanlder(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Set our NoCache headers
		for k, v := range noCacheHeaders {
			w.Header().Set(k, v)
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
