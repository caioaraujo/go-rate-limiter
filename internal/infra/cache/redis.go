package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func Set(client *redis.Client, key, value string, expiration time.Duration) error {
	var ctx = context.Background()
	err := client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func Get(client *redis.Client, value string) (string, error) {
	var ctx = context.Background()
	value, err := client.Get(ctx, value).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

func Connect() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return client
}
