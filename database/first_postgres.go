package database

import (
	"errors"
	"fmt"
	"redis-task/model"
	"slices"
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

func (fp *FirstPostgres) ReorderInputs(input model.ReorderInput) (int, error) {
	//ok, _ := fr.checkCache(redisKey)
	id, err := fp.reorderSavePG(input)
	if err != nil {
		return 0, err
	}
	return id, err

}

func (fp *FirstPostgres) reorderSavePG(input model.ReorderInput) (orderId int, err error) {
	tVal, err := getTextVal(input)
	if err != nil {
		return 0, err
	}
	fmt.Println(err, tVal)
	if tVal == input.Order || input.Order <= 0 {
		return 1, errors.New("same orders or order is negative")
	}
	betweenVals, err := getBetweenVal(input, tVal)
	//if len(betweenVals) == 1 {
	//	changeOnOne(input, betweenVals)
	//}
	if err != nil {
		return
	}
	if tVal < input.Order {
		//todo call func @decBetweens
		_, err := decBetweens(input, betweenVals)
		if err != nil {
			return 0, err
		}

	} else {
		//todo call func @incBetweens
		_, err := incBetweens(input, betweenVals)
		if err != nil {
			return 0, err
		}
	}
	data, err := fp.GetData()
	if err != nil {
		return 0, err
	}
	fmt.Println(betweenVals, data)
	return
}

func getTextVal(input model.ReorderInput) (tVal int, err error) {
	query := fmt.Sprintf(`SELECT %s FROM %s WHERE text=$1`, orderBy, tableFirst)
	fmt.Println(query)
	err = db.PDb.QueryRow(ctx, query, input.Text).Scan(&tVal)
	return
}

func getBetweenVal(input model.ReorderInput, tVal int) (betweenVals []int, err error) {
	var betStr string
	if tVal < input.Order {
		betStr = fmt.Sprintf("%d AND %d", tVal+1, input.Order)
	} else {
		betStr = fmt.Sprintf("%d AND %d", input.Order, tVal-1)
	}

	var betweenVal int
	bquery := fmt.Sprintf(`SELECT %s FROM %s WHERE %s BETWEEN %s`, orderBy, tableFirst, orderBy, betStr)
	rows, err := db.PDb.Query(ctx, bquery)
	for rows.Next() {
		err = rows.Scan(&betweenVal)
		if err != nil {
			return
		}
		betweenVals = append(betweenVals, betweenVal)
	}
	return
}

func decBetweens(input model.ReorderInput, betweens []int) (int, error) {
	var query string

	for _, val := range betweens {
		query += fmt.Sprintf(`UPDATE %s SET %s=%d WHERE %s=%d;`, tableFirst, orderBy, val-1, orderBy, val)
	}
	query += fmt.Sprintf(`UPDATE %s SET %s=%d WHERE text=$1;`, tableFirst, orderBy, input.Order)
	fmt.Println(query)
	_, err := db.PDb.Query(ctx, query, input.Text)
	if err != nil {
		return 0, err
	}
	return 1, err
}

func incBetweens(input model.ReorderInput, betweens []int) (int, error) {
	var querys string
	trx, err := db.PDb.Begin(ctx)
	slices.Reverse(betweens) // Reverse the slice for correct incrementing
	for _, val := range betweens {
		querys = fmt.Sprintf(`UPDATE %s SET %s=%d WHERE %s=%d`, tableFirst, orderBy, val+1, orderBy, val)
		fmt.Println(querys)
		_, err := trx.Exec(ctx, querys)
		if err != nil {
			return 0, err
		}
	}
	queryInc := fmt.Sprintf(`UPDATE %s SET %s=%d WHERE text=$1`, tableFirst, orderBy, input.Order)
	fmt.Println(queryInc)
	_, err = trx.Exec(ctx, queryInc, input.Text)
	if err != nil {
		return 0, err
	}
	return 1, trx.Commit(ctx)
}

//func changeOnOne(input model.ReorderInput, betweens []int) {
//	var queryo string
//	trx, err := db.PDb.Begin(ctx)
//	queryo = fmt.Sprintf(`UPDATE %s SET %s=%d WHERE %s=%d`, tableFirst, orderBy, val+1, orderBy, val)
//	fmt.Println(queryo)
//	_, err := trx.Exec(ctx, querys)
//	if err != nil {
//		return 0, err
//	}
//	queryInc := fmt.Sprintf(`UPDATE %s SET %s=%d WHERE text=$1`, tableFirst, orderBy, input.Order)
//	fmt.Println(queryInc)
//	_, err = trx.Exec(ctx, queryInc, input.Text)
//	if err != nil {
//		return 0, err
//	}
//	return 1, trx.Commit(ctx)
//}
