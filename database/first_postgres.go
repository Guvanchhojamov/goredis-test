package database

import (
	"context"
	"errors"
	"fmt"
	"redis-task/model"
	"strconv"
)

type FirstPostgres struct {
}

var fr = new(FirstRedis)
var db, _ = NewPostgresDB()

const tableFirst = "table_first"
const orderBy = "order_id"

func (fp *FirstPostgres) GetData() (interface{}, error) {

	ok, err := fr.checkCache(redisKey)
	fmt.Println(ok, err)
	if err != nil {
		return nil, err
	}
	if ok {
		data, err := fr.getFromCache()
		return data, err
	}
	data, err := getFromPG()
	return data, err
}

func getFromPG() (result []string, errr error) {
	var orderId string
	var text string
	var forCache [][]string
	query := fmt.Sprintf(`SELECT  ft.order_id, ft.text FROM %s ft ORDER BY %s ASC`, tableFirst, orderBy)
	rows, err := db.PDb.Query(context.Background(), query)
	if err == nil {
		for rows.Next() {
			rows.Scan(&orderId, &text)
			forCache = append(forCache, []string{orderId, text})
			result = append(result, text)
		}
		if len(result) == 0 {
			return
		}
		cache, err := fr.SaveToCache(forCache)
		if err != nil && cache == nil {
			result = append(result, "status:Not saved from DB to Cache")
			return
		}
	}
	result = append(result, "status:from DB saved to Cache")
	return
}

func (fp *FirstPostgres) SaveData(input model.Inputs) (int, error) {
	orderId, cacheString, err := savePG(input)
	fmt.Println(cacheString)
	if err != nil {
		return 0, err
	}
	redResult, err := fr.SaveToCache(cacheString)
	if redResult == 0 && err != nil {
		return 0, errors.New("redis save err")
	}
	return orderId, err
}
func savePG(input model.Inputs) (int, [][]string, error) {
	var orderId string
	var err error
	var forCache [][]string
	for _, val := range input.Text {
		query := fmt.Sprintf(`INSERT INTO %s (text) VALUES ('%s') RETURNING order_id`, tableFirst, val)
		err = db.PDb.QueryRow(context.Background(), query).Scan(&orderId)
		forCache = append(forCache, []string{orderId, val})
	}
	id, _ := strconv.Atoi(orderId)
	return id, forCache, err
}

//	func (fp *FirstPostgres) ReorderInputs(input model.ReorderInput) (interface{}, error) {
//		ok, _ := fr.checkCache(redisKey)
//		if ok {
//			_, err := fr.SaveReorderCache(input)
//			data, err := fp.GetData()
//			return data, err
//		} else {
//			_, _ = fp.GetData()
//			_, err := fr.SaveReorderCache(input)
//			data, err := fp.GetData()
//			return data, err
//		}
//	}
func (fp *FirstPostgres) ReorderInputs(input model.ReorderInput) (interface{}, error) {
	//ok, _ := fr.checkCache(redisKey)
	return nil, nil

}

func (fp *FirstPostgres) ReorderSavePG(newScore int, member string) (orderId int, err error) {
	query := fmt.Sprintf(`UPDATE %s SET order_id = $1  WHERE text=$2 RETURNING order_id`, tableFirst)
	err = db.PDb.QueryRow(context.Background(), query, newScore, member).Scan(&orderId)
	if err != nil {
		return 0, err
	}
	return
}
