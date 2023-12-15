package database

import (
	"fmt"
	"redis-task/model"
)

type SecondRedis struct {
}

const structKey = "users"
const fieldKey = "user"

func (sr *SecondRedis) SaveStructToCache(input model.User) error {
	keyStr := fmt.Sprintf("%s:%s:%v", structKey, fieldKey, input.Id)
	result, err := red.RedisClient.HSet(ctx, keyStr, input).Result()
	if err != nil {
		return err
	}
	fmt.Println(result)
	return nil
}
func (sr *SecondRedis) GetStructFromCache() (users []model.UserResponse, err error) {
	var user model.UserResponse
	keysStr := fmt.Sprintf("%s:%s:%v", structKey, fieldKey, "*")
	keys, err := red.RedisClient.Keys(ctx, keysStr).Result()
	if err != nil {
		return nil, err
	}

	for _, val := range keys {
		err = red.RedisClient.HGetAll(ctx, val).Scan(&user)
		if err != nil {
			return
		}
		users = append(users, user)
	}
	return
}

func (sr *SecondRedis) UpdateStructOnCache(input model.UserUpdate) (user model.UserResponse, err error) {
	return model.UserResponse{}, err
}

//	func checkCacheField(field string) (bool, error) {
//		return red.RedisClient.HExists(structKey, field).Result()
//	}
//func checkStructCacheKey(key string) (int64, error) {
//	return red.RedisClient.Exists(key).Result()
//}
