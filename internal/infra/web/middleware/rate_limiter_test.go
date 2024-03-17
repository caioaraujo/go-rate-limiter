package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/caioaraujo/go-rate-limiter/internal/infra/cache"
	"github.com/magiconair/properties/assert"
)

func TestRateLimiterBlockByMaxReqAllowed(t *testing.T) {
	svr := httptest.NewServer(RateLimiter(SomeHandler))
	defer svr.Close()

	// max duas requisicoes
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
	fmt.Printf("Status response 1: %v", res.StatusCode)
	assert.Equal(t, res.StatusCode, http.StatusOK)

	// segunda requisicao deve passar
	res2, err := http.Get(svr.URL)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer res2.Body.Close()
	fmt.Printf("Status response 2: %v", res2.StatusCode)
	assert.Equal(t, res2.StatusCode, http.StatusOK)

	// terceira requisicao deve bloquear..
	res3, err := http.Get(svr.URL)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	assert.Equal(t, res3.StatusCode, http.StatusTooManyRequests)
	fmt.Printf("Status response 3: %v", res3.StatusCode)
	defer res3.Body.Close()

	time.Sleep(1 * time.Second)

	// uma quarta requisicao deve passar apos o tempo de desbloqueio
	res4, err := http.Get(svr.URL)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	assert.Equal(t, res4.StatusCode, http.StatusOK)
	fmt.Printf("Status response 4: %v", res4.StatusCode)
	defer res4.Body.Close()
}

func TestRateLimiterBlockByIP(t *testing.T) {
	svr := httptest.NewServer(RateLimiter(SomeHandler))
	defer svr.Close()

	// max duas requisicoes por IP
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

	// uma quarta requisicao, com outro IP, deve passar
	req4, err := http.NewRequest("GET", svr.URL, nil)
	req4.Header.Set("X-Real-Ip", "0.0.0.0")
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	httpclient := &http.Client{}
	res4, err := httpclient.Do(req4)
	assert.Equal(t, res4.StatusCode, http.StatusOK)
	defer res4.Body.Close()
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
