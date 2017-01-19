package query

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dictyBase/modware/modwaretest"
)

// IndentJSON uniformly indent the json byte
func IndentJSON(b []byte) []byte {
	var out bytes.Buffer
	_ = json.Indent(&out, b, "", " ")
	return bytes.TrimSpace(out.Bytes())
}

func getfilter(w http.ResponseWriter, r *http.Request) {
	p, ok := r.Context().Value(ContextKeyQueryParams).(*Params)
	if ok {
		if p.HasFilters {
			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(p.Filters)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "no filters")
		}
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprint(w, "No query params")
}

func getfields(w http.ResponseWriter, r *http.Request) {
	p, ok := r.Context().Value(ContextKeyQueryParams).(*Params)
	if ok {
		if p.HasSparseFields {
			w.WriteHeader(http.StatusOK)
			m := make(map[string][]string)
			for k, v := range p.SparseFields {
				m[k] = v.GetAll()
			}
			err := json.NewEncoder(w).Encode(m)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "no fields")
		}
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprint(w, "No query params")
}

func getinclude(w http.ResponseWriter, r *http.Request) {
	p, ok := r.Context().Value(ContextKeyQueryParams).(*Params)
	if ok {
		if p.HasIncludes {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Includes\t%s", strings.Join(p.Includes, ":"))
		} else {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "no includes")
		}
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprint(w, "No query params")
}

func TestFilterError(t *testing.T) {
	testHandlerFn := http.HandlerFunc(getfilter)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(
		"GET",
		"http://dictybase.org/example?filter[name]=foo&filter[country]=argentina",
		nil,
	)
	MiddlewareFn(testHandlerFn).ServeHTTP(w, r)
	if w.Code != http.StatusNotAcceptable {
		t.Fatalf("unexpected http response %d with body %s", w.Code, w.Body.String())
	}
	jsonBlob := []byte(`{
				"errors":[
					{
						"status":"406","title":"Accept header is not acceptable",
						"detail":"Not Acceptable: The given Accept header value  is incorrect for filter query extension",
						"meta":{"creator":"api error helper"}
					}
				]
			}`)
	expJSON := IndentJSON(jsonBlob)
	matchJSON := IndentJSON(w.Body.Bytes())
	if bytes.Compare(expJSON, matchJSON) != 0 {
		t.Fatalf("expected \n%s response does not match with \n%s\n", string(expJSON), string(matchJSON))
	}
}

func TestFilter(t *testing.T) {
	testHandlerFn := http.HandlerFunc(getfilter)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(
		"GET",
		"http://dictybase.org/example?filter[name]=foo&filter[country]=argentina",
		nil,
	)
	r.Header.Set(acceptH, filterMediaType)
	r.Header.Set(contentType, filterMediaType)
	MiddlewareFn(testHandlerFn).ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected http response %d with body %s", w.Code, w.Body.String())
	}
	jsonBlob := []byte(`{"country":"argentina", "name": "foo"}`)
	expJSON := modwaretest.IndentJSON(jsonBlob)
	matchJSON := modwaretest.IndentJSON(w.Body.Bytes())
	if bytes.Compare(expJSON, matchJSON) != 0 {
		t.Fatalf("expected \n%s response does not match with \n%s\n", string(expJSON), string(matchJSON))
	}
}

func TestFields(t *testing.T) {
	testHandlerFn := http.HandlerFunc(getfields)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(
		"GET",
		"http://dictybase.org/example?fields[name]=lola,bantu&fields[blogs]=title",
		nil,
	)
	MiddlewareFn(testHandlerFn).ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected http response %d with body %s", w.Code, w.Body.String())
	}
	jsonBlob := []byte(`{"blogs": ["title"],"name":["lola","bantu"]}`)
	expJSON := modwaretest.IndentJSON(jsonBlob)
	matchJSON := modwaretest.IndentJSON(w.Body.Bytes())
	if bytes.Compare(expJSON, matchJSON) != 0 {
		t.Fatalf("expected \n%s response does not match with \n%s\n", string(expJSON), string(matchJSON))
	}
}

func TestInclude(t *testing.T) {
	testHandlerFn := http.HandlerFunc(getinclude)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://dictybase.org/example?include=foo,bar,baz", nil)
	MiddlewareFn(testHandlerFn).ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected http response %d", w.Code)
	}
	expBody := fmt.Sprintf("Includes\t%s", strings.Join([]string{"foo", "bar", "baz"}, ":"))
	if strings.Compare(expBody, w.Body.String()) != 0 {
		t.Fatalf("unexpected http response body %s\n", w.Body.String())
	}
}
