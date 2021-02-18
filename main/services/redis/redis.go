package redis

import (
	"github.com/go-redis/redis"
	"os"
)

func BuildRedisClient() *redis.Client {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost:6379"
	}

	redis := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: "",
		DB:       0,
	})

	_, err := redis.Ping().Result()
	if err != nil {
		panic(err)
	}

	return redis
}
