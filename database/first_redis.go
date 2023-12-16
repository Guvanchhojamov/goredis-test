package database

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
)

type FirstRedis struct {
}

const inputCacheKey = "inputs"

var red, _ = NewRedisDB()
var fp = new(FirstPostgres)

func (fr *FirstRedis) checkCache(key string) (bool, error) {
	res, err := red.RedisClient.Exists(ctx, key).Result()
	if err != nil || res == 0 {
		return false, err
	}
	return true, err
}

func (fr *FirstRedis) getFromCache() (interface{}, error) {
	result, err := red.RedisClient.ZRange(ctx, inputCacheKey, 0, -1).Result()
	result = append(result, "storage: From cache")
	return result, err
}

func (fr *FirstRedis) SaveToCache(args [][]string) (result interface{}, err error) {
	err = red.RedisClient.Del(ctx, inputCacheKey).Err()
	if err != nil {
		return
	}
	for _, val := range args {
		score, _ := strconv.ParseFloat(val[0], 64)
		opt := redis.Z{
			Score:  score,
			Member: val[1],
		}
		result = red.RedisClient.ZAdd(ctx, inputCacheKey, opt).Args()
	}
	fmt.Println(result, err)
	return
}
