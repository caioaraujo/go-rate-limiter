package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/caioaraujo/go-rate-limiter/internal/infra/cache"
	"github.com/stretchr/testify/assert"
)

func TestRateLimiterBlockByMaxReqAllowed(t *testing.T) {
	svr := httptest.NewServer(RateLimiter(SomeHandler))
	defer svr.Close()

	// max duas requisicoes em 1 seg
	err := setConfig("2", "IP", "1")
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	// primeira requisicao deve passar
	res, err := http.Get(svr.URL)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer res.Body.Close()
	assert.Equal(t, res.StatusCode, http.StatusOK)

	// segunda requisicao deve passar
	res2, err := http.Get(svr.URL)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer res2.Body.Close()
	assert.Equal(t, res2.StatusCode, http.StatusOK)

	// terceira requisicao deve bloquear..
	res3, err := http.Get(svr.URL)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	assert.Equal(t, res3.StatusCode, http.StatusTooManyRequests)
	defer res3.Body.Close()

	time.Sleep(1 * time.Second)

	// uma quarta requisicao deve passar apos o tempo de desbloqueio
	res4, err := http.Get(svr.URL)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	assert.Equal(t, res4.StatusCode, http.StatusOK)
	defer res4.Body.Close()
}

func TestRateLimiterBlockByIP(t *testing.T) {
	svr := httptest.NewServer(RateLimiter(SomeHandler))
	defer svr.Close()

	// max duas requisicoes por IP em 1 seg
	err := setConfig("2", "IP", "1")
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	// primeira requisicao deve passar
	res, err := doRequest("X-Real-Ip", "192.168.222.2", svr)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer res.Body.Close()
	assert.Equal(t, res.StatusCode, http.StatusOK)

	// segunda requisicao deve passar
	res2, err := doRequest("X-Real-Ip", "192.168.222.2", svr)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer res2.Body.Close()
	assert.Equal(t, res2.StatusCode, http.StatusOK)

	// terceira requisicao deve bloquear..
	res3, err := doRequest("X-Real-Ip", "192.168.222.2", svr)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	assert.Equal(t, res3.StatusCode, http.StatusTooManyRequests)
	defer res3.Body.Close()

	// uma quarta requisicao, com outro IP, deve passar
	res4, err := doRequest("X-Real-Ip", "0.0.0.0", svr)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	assert.Equal(t, res4.StatusCode, http.StatusOK)
	defer res4.Body.Close()
}

func TestRateLimiterBlockByToken(t *testing.T) {
	svr := httptest.NewServer(RateLimiter(SomeHandler))
	defer svr.Close()

	// max duas requisicoes por Token em 1 seg
	err := setConfig("2", "TOKEN", "1")
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	// primeira requisicao deve passar
	res1, err := doRequest("API_KEY", "abc123", svr)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	assert.Equal(t, res1.StatusCode, http.StatusOK)
	defer res1.Body.Close()

	// segunda requisicao deve passar
	res2, err := doRequest("API_KEY", "abc123", svr)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	assert.Equal(t, res2.StatusCode, http.StatusOK)
	defer res2.Body.Close()

	// terceira requisicao deve bloquear..
	res3, err := doRequest("API_KEY", "abc123", svr)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	assert.Equal(t, res3.StatusCode, http.StatusTooManyRequests)
	defer res3.Body.Close()

	// uma quarta requisicao, com outro token, deve passar
	res4, err := doRequest("API_KEY", "abc321", svr)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	assert.Equal(t, res4.StatusCode, http.StatusOK)
	defer res4.Body.Close()
}

func doRequest(key, value string, svr *httptest.Server) (*http.Response, error) {
	httpclient := &http.Client{}
	req, err := http.NewRequest("GET", svr.URL, nil)
	req.Header.Set(key, value)
	if err != nil {
		return nil, err
	}
	return httpclient.Do(req)
}

func setConfig(maxReqAllowed, blockMethod, blockTime string) error {
	client := cache.Connect()
	err := cache.Set(client, "MAX_REQ_PERM", maxReqAllowed, 0)
	if err != nil {
		return err
	}
	err = cache.Set(client, "MET_BLOQUEIO", blockMethod, 0)
	if err != nil {
		return err
	}
	err = cache.Set(client, "TEMPO_BLOQUEIO", blockTime, 0)
	if err != nil {
		return err
	}
	return nil
}

func SomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}
