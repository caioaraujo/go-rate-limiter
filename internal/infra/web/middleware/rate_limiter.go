package middleware

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/caioaraujo/go-rate-limiter/internal/infra/cache"

	"encoding/json"
	"net/http"
)

func RateLimiter(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	type userInfo struct {
		timesVisit int
		lastSeen   time.Time
	}
	var (
		users   = make(map[string]*userInfo)
		message = "you have reached the maximum number of requests or actions allowed within a certain time frame"
	)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tempoBloqueio := getIntFromCache("TEMPO_BLOQUEIO")
		maxReqAllowed := getIntFromCache("MAX_REQ_PERM")
		metodoBloqueio := getStringFromCache("MET_BLOQUEIO")

		key := ""
		if metodoBloqueio == "IP" {
			key = getUserIP(r)
		} else if metodoBloqueio == "TOKEN" {
			key = r.Header.Get("API_KEY")
		} else {
			panic("Valor nao permitido para metodo de bloqueio!")
		}

		isBlocked := isUserBlocked(key)

		if isBlocked {
			w.WriteHeader(http.StatusTooManyRequests)
			err := json.NewEncoder(w).Encode(&message)
			if err != nil {
				panic(err)
			}
			return
		}

		if currentInfo, ok := users[key]; !ok {
			users[key] = &userInfo{timesVisit: 1, lastSeen: time.Now()}
		} else {
			if time.Since(currentInfo.lastSeen) > time.Duration(tempoBloqueio)*time.Second {
				users[key] = &userInfo{timesVisit: 1, lastSeen: time.Now()}
				next(w, r)
			}
			currentInfo.timesVisit += 1
			if currentInfo.timesVisit > maxReqAllowed {
				blockUser(key, tempoBloqueio)
				delete(users, key)
				w.WriteHeader(http.StatusTooManyRequests)
				err := json.NewEncoder(w).Encode(&message)
				if err != nil {
					panic(err)
				}
				return
			}
		}

		next(w, r)

	})
}

func isUserBlocked(key string) bool {
	var currentBlockedUser string
	currentBlockedUser = getStringFromCache(key)
	if currentBlockedUser == "1" {
		return true
	}
	return false
}

func blockUser(key string, timeLimit int) {
	limitSeconds := fmt.Sprintf("%ds", timeLimit)
	duration, err := time.ParseDuration(limitSeconds)
	if err != nil {
		panic(err)
	}
	setInCache(key, "1", duration)
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
