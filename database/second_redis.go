package database

import (
	"fmt"
	"redis-task/model"
)

type SecondRedis struct {
}

func NewSecondRedis() *SecondRedis {
	return &SecondRedis{}
}

const (
	structKey = "users"
	fieldKey  = "user"
)

func (sr *SecondRedis) SaveStructToCache(input model.User) error {
	keyStr := fmt.Sprintf("%s:%s:%v", structKey, fieldKey, input.Id)
	fmt.Println(input)
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

func (sr *SecondRedis) UpdateStructOnCache(input model.UserUpdate, id string) (err error) {
	keyStr := fmt.Sprintf("%s:%s:%s", structKey, fieldKey, id)
	updateStruct := generateUpdateString(input)
	fullStr := fmt.Sprintf("%s %s", keyStr, updateStruct)
	fmt.Println(fullStr)
	res, err := red.RedisClient.HSet(ctx, keyStr, updateStruct).Result()
	fmt.Println(res, err)
	if err != nil {
		return
	}
	return
}
func generateUpdateString(input model.UserUpdate) interface{} {
	var updateStruct = make(map[string]interface{})
	if input.Username != nil {
		updateStruct["username"] = *input.Username
	}
	if input.Age != nil {
		updateStruct["age"] = *input.Age
	}
	if input.Address != nil {
		updateStruct["address"] = *input.Address
	}
	return updateStruct
}
