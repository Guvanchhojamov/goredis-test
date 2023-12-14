package database

import (
	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type RedisDB struct {
	RedisClient *redis.Client
}

func NewRedisDB() (*RedisDB, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: "",
		DB:       0,
	})
	if err := redisClient.Ping().Err(); err != nil {
		log.Fatalf("redis ping error: %v\n", err)
		return nil, err
	}
	return &RedisDB{RedisClient: redisClient}, nil
}
