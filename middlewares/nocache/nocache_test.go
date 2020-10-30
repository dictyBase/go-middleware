package nocache

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type testHandler struct{}

func (h *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello world")
}

func TestHandler(t *testing.T) {
	ts := httptest.NewServer(Handler(&testHandler{}))
	defer ts.Close()
	assert := require.New(t)
	res, err := ts.Client().Get(ts.URL)
	assert.NoError(err, "expect no error from http get call")
	defer res.Body.Close()
	assert.Equal(res.StatusCode, http.StatusOK, "should be successful http request")
	assert.Equal(res.Header.Get("Pragma"), "no-cache", "should match Pragma header value")
	assert.Equal(
		res.Header.Get("Expires"),
		time.Unix(0, 0).Format(time.RFC1123),
		"should match Expires header value",
	)
	assert.Equal(
		res.Header.Get("Cache-Control"),
		"no-cache, no-store, no-transform, must-revalidate, private, max-age=0",
		"should match Cache-Control header value",
	)
	assert.Equal(
		res.Header.Get("X-Accel-Expires"),
		"0",
		"should match X-Accel-Expires header",
	)
}
