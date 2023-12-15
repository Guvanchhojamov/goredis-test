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
	fmt.Println(keyStr)
	result, err := red.RedisClient.HSet(ctx, keyStr, input).Result()
	if err != nil {
		return err
	}
	fmt.Println(result)
	return nil
}
func (sr *SecondRedis) GetStructFromCache() (users []model.UserResponse, err error) {
	return nil, err
}

func (sr *SecondRedis) UpdateStructOnCache(input model.UserUpdate) (user model.UserResponse, err error) {
	return model.UserResponse{}, err
}

//func checkCacheField(field string) (bool, error) {
//	return red.RedisClient.HExists(structKey, field).Result()
//}
//func checkStructCacheKey() (int64, error) {
//	return red.RedisClient.Exists(structKey).Result()
//}
