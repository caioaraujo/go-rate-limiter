package middleware

import (
	"golang.org/x/time/rate"

	"encoding/json"
	"net/http"
)

func RateLimiter(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	limiter := rate.NewLimiter(2, 4)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			message := "you have reached the maximum number of requests or actions allowed within a certain time frame"
			w.WriteHeader(http.StatusTooManyRequests)
			err := json.NewEncoder(w).Encode(&message)
			if err != nil {
				panic(err)
			}
			return
		} else {
			next(w, r)
		}
	})
}
