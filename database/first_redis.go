package database

import (
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
)

type FirstRedis struct {
}

const redisKey = "inputs"

var rDB, _ = NewRedisDB()
var fp = new(FirstPostgres)

func (fr *FirstRedis) checkCache(key string) (bool, error) {
	res, err := rDB.RedisClient.Exists(key).Result()
	fmt.Println(err, res)
	if err != nil || res == 0 {
		return false, err
	}
	return true, err
}

func (fr *FirstRedis) getFromCache() (interface{}, error) {

	result, err := rDB.RedisClient.ZRange(redisKey, 0, -1).Result()
	result = append(result, "storage: From cache")
	return result, err
}

func (fr *FirstRedis) SaveToCache(args [][]string) (result interface{}, err error) {
	for _, val := range args {
		float, _ := strconv.ParseFloat(val[0], 64)
		opt := redis.Z{
			Score:  float,
			Member: val[1],
		}
		result = rDB.RedisClient.ZAdd(redisKey, opt).Args()
	}

	fmt.Println(result, err)
	return
}

//func (fr *FirstRedis) SaveReorderCache(input model.ReorderInput) (int, error) {
//	tOneScore := rDB.RedisClient.ZScore(redisKey, input.Text1).Val()
//	tTwoScore := rDB.RedisClient.ZScore(redisKey, input.Text2).Val()
//	fmt.Println(tOneScore, tTwoScore)
//	var newOneScore = tOneScore
//	if tOneScore <= tTwoScore {
//		newOneScore = tTwoScore + 1
//	} else if tOneScore >= tTwoScore {
//		if tTwoScore == 1 {
//			newOneScore = tTwoScore
//			tTwoScore += 1
//		}
//		newOneScore = tTwoScore - 1
//		tTwoScore += 1
//	}
//	oneOpt := redis.Z{
//		Score:  newOneScore,
//		Member: input.Text1,
//	}
//	fmt.Println(oneOpt, tOneScore, tTwoScore, newOneScore)
//	res, err := rDB.RedisClient.ZAddXXCh(redisKey, oneOpt).Result()
//	if res <= 0 {
//		return 0, errors.New("reorder save cache error")
//	}
//	id, err := fp.ReorderSavePG(int(newOneScore), input.Text1)
//	if err != nil {
//		rDB.RedisClient.FlushDB()
//		return 0, err
//	}
//	return id, err
//}
