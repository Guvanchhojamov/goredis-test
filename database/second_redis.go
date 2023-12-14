package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"redis-task/model"
)

type SecondRedis struct {
}

var red, _ = NewRedisDB()

const structKey = "structs"
const fieldKey = "user"

func (sr *SecondRedis) SaveStructToCache(input model.User) error {
	data, err := json.Marshal(input)
	field := fmt.Sprintf("%s:%d", fieldKey, input.Id)
	if ok, _ := checkCacheField(field); ok {
		return errors.New("error user already in cache")
	}
	result, err := red.RedisClient.HSetNX(structKey, field, data).Result()
	if err != nil {
		return err
	}
	if !result {
		return errors.New("error not saved to cacche")
	}
	fmt.Println(result, err)
	return err
}
func (sr *SecondRedis) GetStructFromCache() (users []model.UserResponse, err error) {
	key, err := checkStructCacheKey()
	if err != nil {
		return nil, err
	}
	if key == 0 {
		return nil, errors.New("no any saved struct")
	}
	sLen := red.RedisClient.HLen(structKey).Val()
	var user model.UserResponse
	var id int
	for i := 0; i < int(sLen); i++ {
		id++
		result, err := red.RedisClient.HGet(structKey, fmt.Sprintf("%s:%d", fieldKey, id)).Result()
		if err != nil {
			return nil, errors.New("error get from cache")
		}
		err = json.Unmarshal([]byte(result), &user)
		users = append(users, user)
	}
	fmt.Println(users, err)
	return
}

func (sr *SecondRedis) UpdateStructOnCache(input model.UserUpdate) (user model.UserResponse, err error) {
	key, err := checkStructCacheKey()
	if err != nil || key == 0 {
		return user, errors.New("key not found in cache" + err.Error())
	}
	field := fmt.Sprintf("%s:%d", fieldKey, input.Id)
	fmt.Println(field)
	if ok, _ := checkCacheField(field); !ok {
		return user, errors.New("user not found in cache")
	}
	data, err := json.Marshal(input)
	if err != nil {
		return
	}
	err = red.RedisClient.HSet(structKey, field, data).Err()
	if err != nil {
		return user, errors.New("update error" + err.Error())
	}

	res, err := red.RedisClient.HGet(structKey, field).Result()
	err = json.Unmarshal([]byte(res), &user)
	return
}

func checkCacheField(field string) (bool, error) {
	return red.RedisClient.HExists(structKey, field).Result()
}
func checkStructCacheKey() (int64, error) {
	return red.RedisClient.Exists(structKey).Result()
}
