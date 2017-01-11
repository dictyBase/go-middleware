package query

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

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

func TestIncludeParams(t *testing.T) {
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
