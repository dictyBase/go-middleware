// Package query support parsing of JSON API based query parameters.
// It supports the following parameters
//		include - /url?include=foo,bar,baz
//		fields(sparse fieldsets) - /url?fields[articles]=title,body&fields[people]=name
//		filter - /url?filter[name]=foo&filter[country]=argentina
// The include and fields are part of JSON API whereas filter is a custom
// extension for dictybase. For details look here
// https://github.com/json-api/json-api/blob/9c7a03dbc37f80f6ca81b16d444c960e96dd7a57/extensions/index.md#-extension-negotiation
// and here
// https://github.com/dictyBase/Migration/blob/master/Webservice-specs.md#filtering
// This middleware terminates the chain in case of incorrect or inappropriate
// http headers for filter query parameters.
package query

import "net/http"

type contextKey string

// String output the details of context key
func (c contextKey) String() string {
	return "pagination context key " + string(c)
}

var (
	ContextKeyQueryParams = contextKey("jsparams")
)

// Params is container for various query parameters
type Params struct {
	// contain include query paramters
	Includes []string
	// contain fields query paramters
	Fields map[string][]string
	// contain filter query parameters
	Filters map[string]string
}

// MiddlewareFn parses the includes, fields and filter query strings and stores
// it in request context under  ContextKeyQueryParam variable as a Params type
// For filter query parameters, the client should include the appropiate media
// type and media type parameters as described here
// https://github.com/dictyBase/Migration/blob/master/Webservice-specs.md#dictybase-specifications.
// Otherwise, the request never gets passed to the handler and either of
// 406(Not Acceptable) or 415(Unsupported Media Type) http status is returned.
func MiddlewareFn(fn http.HandlerFunc) http.HandlerFunc {
	newFn := func(w http.ResponseWriter, r *http.Request) {
		fn(w, r)
	}
	return newFn
}
