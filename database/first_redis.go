package database

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
)

type FirstRedis struct {
}

const redisKey = "inputs"

var red, _ = NewRedisDB()
var fp = new(FirstPostgres)

func (fr *FirstRedis) checkCache(key string) (bool, error) {
	res, err := red.RedisClient.Exists(ctx, key).Result()
	fmt.Println(err, res)
	if err != nil || res == 0 {
		return false, err
	}
	return true, err
}

func (fr *FirstRedis) getFromCache() (interface{}, error) {
	result, err := red.RedisClient.ZRange(ctx, redisKey, 0, -1).Result()
	result = append(result, "storage: From cache")
	return result, err
}

func (fr *FirstRedis) SaveToCache(args [][]string) (result interface{}, err error) {
	err = red.RedisClient.Del(ctx, redisKey).Err()
	if err != nil {
		return nil, err
	}
	for _, val := range args {
		float, _ := strconv.ParseFloat(val[0], 64)
		opt := redis.Z{
			Score:  float,
			Member: val[1],
		}

		result = red.RedisClient.ZAdd(ctx, redisKey, opt).Args()
	}
	fmt.Println(result, err)
	return
}
