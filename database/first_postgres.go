package database

import (
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
	rows, err := db.PDb.Query(ctx, query)
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
	var lastorderId int
	trx, err := db.PDb.Begin(ctx)
	if err != nil {
		return 0, nil, err
	}
	query := fmt.Sprintf(`SELECT COALESCE(MAX(%s) , 0) FROM %s`, orderBy, tableFirst)
	row := trx.QueryRow(ctx, query)
	if err = row.Scan(&lastorderId); err != nil {
		trx.Rollback(ctx)
		return 0, nil, err
	}
	fmt.Printf("last: %d", lastorderId)
	for _, val := range input.Text {
		lastorderId++
		query := fmt.Sprintf(`INSERT INTO %s (text,order_id) VALUES ('%s', %d) RETURNING order_id`, tableFirst, val, lastorderId)
		err = trx.QueryRow(ctx, query).Scan(&orderId)
		if err != nil {
			trx.Rollback(ctx)
			return 0, nil, err
		}
		forCache = append(forCache, []string{orderId, val})
	}

	id, _ := strconv.Atoi(orderId)
	return id, forCache, trx.Commit(ctx)
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
func (fp *FirstPostgres) ReorderInputs(input model.ReorderInput) (int, error) {
	//ok, _ := fr.checkCache(redisKey)
	return nil, nil

}

func (fp *FirstPostgres) reorderSavePG(input model.ReorderInput) (orderId int, err error) {
	trx, err := db.PDb.Begin(ctx)
	if err != nil {
		return
	}
	if err != nil {
		return
	}
	var tVal int
	query := fmt.Sprintf(`SELECT %s FROM %s WHERE text=$1 RETURNING %s`, orderBy, tableFirst, orderBy)
	if err = trx.QueryRow(ctx, query, input.Text).Scan(&tVal); err != nil {
		trx.Rollback()
		return
	}
	return 0, trx.Commit(ctx)
}
