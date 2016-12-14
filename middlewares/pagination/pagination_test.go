package pagination

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func get(w http.ResponseWriter, r *http.Request) {
	prop, ok := r.Context().Value(ContextKeyPagination).(*Props)
	if ok {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Got page size %d and page number %d", prop.Entries, prop.Current)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprint(w, "No pagination property")
}

func TestMiddlewareFn(t *testing.T) {
	testHandlerFn := http.HandlerFunc(get)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://dictybase.org/example?page[number]=6&page[size]=10", nil)
	MiddlewareFn(testHandlerFn).ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected http response %d", w.Code)
	}
	expBody := "Got page size 10 and page number 6"
	if strings.Compare(expBody, w.Body.String()) != 0 {
		t.Fatalf("unexpected http response body %s\n", w.Body.String())
	}
}
