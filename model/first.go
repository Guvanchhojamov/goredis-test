package model

type Inputs struct {
	Text []string `json:"text" db:"text" binding:"required"`
}

type ReorderInput struct {
	Text  string `json:"text" db:"text" binding:"required"`
	Order int    `json:"order" db:"order_id" binding:"required" `
}
