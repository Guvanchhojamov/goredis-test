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

const (
	tableFirst = "table_first"
	orderBy    = "order_id"
	decrease   = "-1"
	increase   = "+1"
)

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
	query := fmt.Sprintf(`SELECT ft.order_id, ft.text FROM %s ft ORDER BY %s ASC`, tableFirst, orderBy)
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

func (fp *FirstPostgres) ReorderInput(input model.ReorderInput) (data interface{}, err error) {
	isChanged, err := fp.reorderSavePG(input)
	if err != nil || !isChanged {
		return
	}
	red.RedisClient.Del(ctx, redisKey)
	data, err = fp.GetData()
	if err != nil {
		return
	}
	return
}

func (fp *FirstPostgres) reorderSavePG(input model.ReorderInput) (isChanged bool, err error) {
	textOrder, err := getTextOrder(input)
	if err != nil {
		return false, err
	}
	fmt.Println(err, textOrder)
	if textOrder == input.Order || input.Order <= 0 {
		return false, errors.New("same orders or order is negative")
	}
	var operator string
	var betweenStr string
	if textOrder < input.Order {
		// #case when order of text small then input order  (Ex: text: "a"=1 order:5 ===>  1<5)
		betweenStr = fmt.Sprintf("%d AND %d", textOrder+1, input.Order)
		operator = decrease
	} else {
		// #case when order of text bigger then input order  (Ex: text: "a"=4 order:1 ===>  4>1)
		betweenStr = fmt.Sprintf("%d AND %d", input.Order, textOrder-1)
		operator = increase
	}
	return changePosition(input, betweenStr, operator)
}

func getTextOrder(input model.ReorderInput) (textOrder int, err error) {
	query := fmt.Sprintf(`SELECT %s FROM %s WHERE text=$1`, orderBy, tableFirst)
	fmt.Println(query)
	err = db.PDb.QueryRow(ctx, query, input.Text).Scan(&textOrder)
	return
}

//func getBetweenVal(input model.ReorderInput, tVal int) (betweenVals []int, betStr string, err error) {
//	if tVal < input.Order {
//		betStr = fmt.Sprintf("%d AND %d", tVal+1, input.Order)
//	} else {
//		betStr = fmt.Sprintf("%d AND %d", input.Order, tVal-1)
//	}
//	fmt.Println(betStr)
//	var betweenVal int
//	bquery := fmt.Sprintf(`SELECT %s FROM %s WHERE %s BETWEEN %s`, orderBy, tableFirst, orderBy, betStr)
//	rows, err := db.PDb.Query(ctx, bquery)
//	for rows.Next() {
//		err = rows.Scan(&betweenVal)
//		if err != nil {
//			return
//		}
//		betweenVals = append(betweenVals, betweenVal)
//	}
//	return
//}

func changePosition(input model.ReorderInput, betweenStr string, operation string) (isCahnged bool, err error) {
	var query string
	query = fmt.Sprintf(`UPDATE %s SET %s=COALESCE(%s,0)%s WHERE %s BETWEEN %s AND text<>$1`, tableFirst, orderBy, orderBy, operation, orderBy, betweenStr)
	fmt.Println(query)
	trx, err := db.PDb.Begin(ctx)
	if err != nil {
		return false, err
	}
	if _, err = trx.Exec(ctx, query, input.Text); err != nil {
		return false, err
	}
	query = fmt.Sprintf(`UPDATE %s SET %s=$1 WHERE text=$2`, tableFirst, orderBy)
	if _, err = trx.Exec(ctx, query, input.Order, input.Text); err != nil {
		return false, err
	}
	return true, trx.Commit(ctx)
}
