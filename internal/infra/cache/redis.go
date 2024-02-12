package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func Set(key, value string) error {
	client := GetClient()
	err := client.Set(ctx, key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func Get(value string) (string, error) {
	client := GetClient()
	value, err := client.Get(ctx, "key").Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

func GetClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return client
}
