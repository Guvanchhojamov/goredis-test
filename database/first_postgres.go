package database

import (
	"errors"
	"fmt"
	"redis-task/model"
)

type FirstPostgres struct {
}

var firstRedis = new(FirstRedis)
var db, _ = NewPostgresDB()

const (
	tableFirst = "table_first"
	orderField = "order_id"
	decrease   = "-1"
	increase   = "+1"
)

func (fp *FirstPostgres) GetData() (interface{}, error) {
	ok, err := firstRedis.checkCache(inputCacheKey)
	fmt.Println(ok, err)
	if err != nil {
		return nil, err
	}
	if ok {
		data, err := firstRedis.getFromCache()
		return data, err
	}
	data, err := getFromDataBase()
	return data, err
}

func getFromDataBase() (result []string, err error) {
	var (
		orderId  string
		text     string
		forCache [][]string
	)
	query := fmt.Sprintf(`SELECT ft.order_id, ft.text FROM %s ft ORDER BY %s ASC`, tableFirst, orderField)
	rows, err := db.PDb.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		rows.Scan(&orderId, &text)
		forCache = append(forCache, []string{orderId, text})
		result = append(result, text)
	}
	if len(result) == 0 {
		return result, err
	}
	cache, err := firstRedis.SaveToCache(forCache)
	if err != nil && cache == nil {
		result = append(result, "status:Not saved from DB to Cache")
		return result, err
	}

	result = append(result, "status:from DB saved to Cache")
	return result, err
}

func (fp *FirstPostgres) SaveData(input model.Inputs) (int, error) {
	id, err := saveToDataBase(input)
	if err != nil {
		return 0, err
	}
	err = red.RedisClient.Del(ctx, inputCacheKey).Err()
	if err != nil {
		return 0, err
	}
	_, err = fp.GetData()
	return id, err
}
func saveToDataBase(input model.Inputs) (id int, err error) {
	trx, err := db.PDb.Begin(ctx)
	if err != nil {
		return
	}
	query := fmt.Sprintf(`SELECT COALESCE(MAX(%s) , 0) FROM %s`, orderField, tableFirst)
	var maxOrder int
	if err = trx.QueryRow(ctx, query).Scan(&maxOrder); err != nil {
		trx.Rollback(ctx)
		return
	}
	insertValuesStr := generateInsertValues(input, maxOrder)
	query = fmt.Sprintf(`INSERT INTO %s (text,order_id) VALUES %s`, tableFirst, insertValuesStr)
	fmt.Println(query)
	_, err = trx.Exec(ctx, query)
	if err != nil {
		trx.Rollback(ctx)
		return
	}
	return 1, trx.Commit(ctx)
}

func (fp *FirstPostgres) ReorderInput(input model.ReorderInput) (data interface{}, err error) {
	isChanged, err := fp.reorderSaveToDatabase(input)
	if err != nil || !isChanged {
		return
	}
	red.RedisClient.Del(ctx, inputCacheKey)
	data, err = fp.GetData()
	if err != nil {
		return
	}
	return
}

func (fp *FirstPostgres) reorderSaveToDatabase(input model.ReorderInput) (isChanged bool, err error) {
	textOrder, err := getTextOrder(input)
	if err != nil {
		return false, err
	}
	fmt.Println(err, textOrder)
	if textOrder == input.Order || input.Order <= 0 {
		return false, errors.New("same orders or order is negative")
	}
	var (
		operator   string
		betweenStr string
	)
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
	query := fmt.Sprintf(`SELECT %s FROM %s WHERE text=$1`, orderField, tableFirst)
	fmt.Println(query)
	err = db.PDb.QueryRow(ctx, query, input.Text).Scan(&textOrder)
	return
}

func changePosition(input model.ReorderInput, betweenStr string, operation string) (isCahnged bool, err error) {
	var query string
	query = fmt.Sprintf(`UPDATE %s SET %s=COALESCE(%s,0)%s WHERE %s BETWEEN %s AND text<>$1`, tableFirst, orderField, orderField, operation, orderField, betweenStr)
	fmt.Println(query)
	trx, err := db.PDb.Begin(ctx)
	if err != nil {
		return false, err
	}
	if _, err = trx.Exec(ctx, query, input.Text); err != nil {
		return false, err
	}
	query = fmt.Sprintf(`UPDATE %s SET %s=$1 WHERE text=$2`, tableFirst, orderField)
	if _, err = trx.Exec(ctx, query, input.Order, input.Text); err != nil {
		return false, err
	}
	return true, trx.Commit(ctx)
}

func generateInsertValues(input model.Inputs, maxOrder int) (insertValuesStr string) {
	for i, val := range input.Text {
		maxOrder++
		if i == len(input.Text)-1 {
			insertValuesStr += fmt.Sprintf(`('%s', %d)`, val, maxOrder) // poslednaya bez zapyataya
			return
		}
		insertValuesStr += fmt.Sprintf(`('%s', %d),`, val, maxOrder)
	}
	return
}
