package database

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
)

type RedisDB struct {
	RedisClient *redis.Client
}

var ctx = context.Background()

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
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("redis ping error: %v\n", err)
		return nil, err
	}
	return &RedisDB{RedisClient: redisClient}, nil
}
