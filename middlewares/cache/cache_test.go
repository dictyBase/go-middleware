package cache

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
	curr := time.Now()
	c := NewHTTPCache(11, curr.AddDate(0, 11, 0))
	ts := httptest.NewServer(c.Handler(&testHandler{}))
	assert := require.New(t)
	res, err := ts.Client().Get(ts.URL)
	assert.NoError(err, "expect no error from http get call")
	defer res.Body.Close()
	assert.Equal(res.StatusCode, http.StatusOK, "should be successful http request")
	assert.Equal(
		res.Header.Get("Expires"),
		curr.AddDate(0, 11, 0).Format(http.TimeFormat),
		"should match Expires header value",
	)
	assert.Equal(
		res.Header.Get("Cache-Control"),
		fmt.Sprintf(
			"public, max-age=%d",
			int(time.Until(curr.AddDate(0, 11, 0)).Seconds()),
		),
		"should match Cache-Control header value",
	)
}
