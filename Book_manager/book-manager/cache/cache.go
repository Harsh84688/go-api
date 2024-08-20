package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client
var ctx = context.Background()

func InitCache() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // No password set
		DB:       0,  // use default DB
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic("Could not connect to Redis")
	}
}

func SetCachedResult(key string, result string, expiration time.Duration) error {
	err := rdb.Set(ctx, key, result, expiration).Err()
	return err
}

// Retrieve cached result from Redis
func GetCachedResult(key string) (string, error) {
	data, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		// Key does not exist
		return "", nil
	} else if err != nil {
		return "", err
	}
	return data, nil
}

func DeleteCachedResult(key string) {
	err := rdb.Del(ctx, key).Err()
	if err != nil {
		panic("Could not delete cached result")
	}
}
