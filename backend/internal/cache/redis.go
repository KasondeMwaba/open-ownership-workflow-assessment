package cache

import "github.com/redis/go-redis/v9"

func OpenRedis(redisURL string) *redis.Client {
	options, err := redis.ParseURL(redisURL)
	if err != nil {
		options = &redis.Options{Addr: "localhost:6379"}
	}
	return redis.NewClient(options)
}
