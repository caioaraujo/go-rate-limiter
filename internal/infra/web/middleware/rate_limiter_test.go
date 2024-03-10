package middleware

import (
	"github.com/magiconair/properties/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRateLimiterByIP(t *testing.T) {
	svr := httptest.NewServer(RateLimiter(SomeHandler))
	defer svr.Close()

	res, err := http.Get(svr.URL)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	assert.Equal(t, res.StatusCode, http.StatusOK)
	res2, err := http.Get(svr.URL)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	assert.Equal(t, res2.StatusCode, http.StatusOK)
	res3, err := http.Get(svr.URL)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	assert.Equal(t, res3.StatusCode, http.StatusTooManyRequests)
}

func SomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}
