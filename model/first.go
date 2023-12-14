package model

type Input struct {
	OrderId string `json:"orderId" db:"order_id"`
	Text    string `json:"text" db:"text" binding:"required"`
}

type Inputs struct {
	Text []string `json:"text" db:"text" binding:"required"`
}

type ReorderInput struct {
	Text1   string `json:"text1" db:"text" binding:"required"`
	OrderId string `json:"order_id" db:"order_id" binding:"required" `
}
