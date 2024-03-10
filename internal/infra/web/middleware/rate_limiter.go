package middleware

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/caioaraujo/go-rate-limiter/internal/infra/cache"

	"golang.org/x/time/rate"

	"encoding/json"
	"net/http"
)

func RateLimiter(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	maxReqAllowed := getIntFromCache("MAX_REQ_PERM")
	tempoBloqueio := getIntFromCache("TEMPO_BLOQUEIO")
	metodoBloqueio := getStringFromCache("MET_BLOQUEIO")

	limiter := rate.NewLimiter(1, maxReqAllowed)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		currentToken := r.Header.Get("API_KEY")
		currentIP := getUserIP(r)

		isBlocked := isUserBlocked(currentIP, currentToken, metodoBloqueio)

		if isBlocked {
			message := "you have reached the maximum number of requests or actions allowed within a certain time frame"
			w.WriteHeader(http.StatusTooManyRequests)
			err := json.NewEncoder(w).Encode(&message)
			if err != nil {
				panic(err)
			}
			return
		}

		if !limiter.Allow() {
			message := "you have reached the maximum number of requests or actions allowed within a certain time frame"
			w.WriteHeader(http.StatusTooManyRequests)
			err := json.NewEncoder(w).Encode(&message)
			if err != nil {
				panic(err)
			}
			setUserLimit(currentIP, currentToken, metodoBloqueio, tempoBloqueio)
			return
		} else {
			next(w, r)
		}
	})
}

func isUserBlocked(ip, token, blockMethod string) bool {
	var currentBlockedUser string
	if blockMethod == "IP" {
		currentBlockedUser = getStringFromCache(ip)
	} else if blockMethod == "TOKEN" {
		currentBlockedUser = getStringFromCache(token)
	} else {
		panic("Método de bloqueio não configurado")
	}
	if currentBlockedUser == "1" {
		return true
	}
	return false
}

func setUserLimit(ip, token, blockMethod string, timeLimit int) {
	limitSeconds := fmt.Sprintf("%ds", timeLimit)
	duration, err := time.ParseDuration(limitSeconds)
	if err != nil {
		panic(err)
	}
	if blockMethod == "IP" {
		setInCache(ip, "1", duration)
	}
	if blockMethod == "TOKEN" {
		setInCache(token, "1", duration)
	}
}

func getStringFromCache(key string) string {
	client := cache.Connect()
	value, ok := cache.Get(client, key)
	if ok != nil {
		return ""
	}
	return value
}

func getIntFromCache(key string) int {
	value := getStringFromCache(key)
	if value == "" {
		return 0
	}
	valueInt, err := strconv.Atoi(value)
	if err != nil {
		panic(err)
	}
	return valueInt
}

func setInCache(key, value string, duration time.Duration) {
	client := cache.Connect()
	err := cache.Set(client, key, value, duration)
	if err != nil {
		panic(err)
	}
}

func getUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	IPSplitted := strings.Split(IPAddress, ":")
	return IPSplitted[0]
}
