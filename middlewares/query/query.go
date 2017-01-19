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

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/dictyBase/apihelpers/apherror"
	"github.com/manyminds/api2go"
)

type contextKey string

// String output the details of context key
func (c contextKey) String() string {
	return "pagination context key " + string(c)
}

var (
	// ContextKeyQueryParams is the key used for storing Params struct in
	// request context
	ContextKeyQueryParams = contextKey("jsparams")
	acceptH               = http.CanonicalHeaderKey("accept")
	contentType           = http.CanonicalHeaderKey("content-type")
	filterMediaType       = strconv.Quote(`application/vnd.api+json; supported-ext="dictybase/filtering-resouce"`)
	qregx                 = regexp.MustCompile(`^\w+\[(\w+)\]$`)
)

// Fields is the container for field names
type Fields struct {
	// Flag to indicate if the field's type name is from a relationship
	// resource
	Relationship bool
	names        []string
}

// GetAll returns all the fields name
func (fl *Fields) GetAll() []string {
	return fl.names
}

// Params is container for various query parameters
type Params struct {
	// contain include query paramters
	Includes []string
	// contain fields query paramters
	SparseFields map[string]*Fields
	// contain filter query parameters
	Filters map[string]string
	// check for presence of fields parameters
	HasSparseFields bool
	// check for presence of include parameters
	HasIncludes bool
	// check for presence of filter parameters
	HasFilters bool
}

func newParams() *Params {
	return &Params{
		Filters:      make(map[string]string),
		SparseFields: make(map[string]*Fields),
	}
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
		params := newParams()
		values := r.URL.Query()
		for k, v := range values {
			switch {
			case strings.HasPrefix(k, "filter"):
				// check for correct header
				if !validateHeader(w, r) {
					return
				}
				if m := qregx.FindStringSubmatch(k); m != nil {
					params.Filters[m[1]] = v[0]
					if !params.HasFilters {
						params.HasFilters = true
					}
				} else {
					apherror.JSONAPIError(
						w,
						apherror.ErrQueryParam.New(
							fmt.Sprintf("Unable to match filter query param %s", v[0]),
						),
					)
					return
				}
			case strings.HasPrefix(k, "fields"):
				if m := qregx.FindStringSubmatch(k); m != nil {
					f := &Fields{Relationship: false}
					if strings.Contains(v[0], ",") {
						f.names = strings.Split(v[0], ",")
						params.SparseFields[m[1]] = f
					} else {
						f.names = []string{v[0]}
						params.SparseFields[m[1]] = f
					}
					if !params.HasSparseFields {
						params.HasSparseFields = true
					}
				} else {
					apherror.JSONAPIError(
						w,
						apherror.ErrQueryParam.New(
							fmt.Sprintf("Unable to match fields query param %s", v[0]),
						),
					)
					return
				}
			case k == "include":
				if strings.Contains(v[0], ",") {
					params.Includes = strings.Split(v[0], ",")
				} else {
					params.Includes = []string{v[0]}
				}
				if !params.HasIncludes {
					params.HasIncludes = true
				}
			default:
				continue
			}
		}
		if params.HasFilters || params.HasSparseFields || params.HasIncludes {
			ctx := context.WithValue(r.Context(), ContextKeyQueryParams, params)
			fn(w, r.WithContext(ctx))
		} else {
			fn(w, r)
		}
	}
	return newFn
}

func queryParamError(w http.ResponseWriter, status int, title, detail string) {
	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(status)
	jsnErr := api2go.Error{
		Status: strconv.Itoa(status),
		Title:  title,
		Detail: detail,
		Meta: map[string]interface{}{
			"creator": "query middleware",
		},
	}
	err := json.NewEncoder(w).Encode(api2go.HTTPError{Errors: []api2go.Error{jsnErr}})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func validateHeader(w http.ResponseWriter, r *http.Request) bool {
	if r.Header.Get(acceptH) != filterMediaType {
		apherror.JSONAPIError(
			w,
			apherror.ErrNotAcceptable.New(
				fmt.Sprintf(
					"The given Accept header value %s is incorrect for filter query extension",
					r.Header.Get(acceptH),
				),
			),
		)
		return false
	}
	if r.Header.Get(acceptH) != r.Header.Get(contentType) {
		apherror.JSONAPIError(
			w,
			apherror.ErrUnsupportedMedia.New(
				fmt.Sprintf(
					"The given media type %s in Content-Type header is not supported",
					r.Header.Get(contentType),
				),
			),
		)
		return false
	}
	return true
}
