package model

type User struct {
	Id       int     `json:"id" redis:"id" binding:"required"`
	Username string  `json:"username" redis:"username" binding:"required"`
	Age      float64 `json:"age" redis:"age"`
	Address  string  `json:"address" redis:"address"`
}

type UserResponse struct {
	Id       int     `json:"id" redis:"id" `
	Username string  `json:"username" redis:"username"`
	Age      float64 `json:"age" redis:"age"`
	Address  string  `json:"address" redis:"address"`
}

type UserUpdate struct {
	Username string  `json:"username" redis:"username"`
	Age      float64 `json:"age" redis:"age" `
	Address  string  `json:"address" redis:"address"`
}
