package cache

import (
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheInterface interface {
	Connect() *redis.Client
	Set(client *redis.Client, key, value string, expiration time.Duration) error
	Get(client *redis.Client, key string) (string, error)
}
