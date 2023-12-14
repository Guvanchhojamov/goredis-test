package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"redis-task/database"
	"redis-task/model"
)

var fpr = new(database.FirstPostgres)

func (h *Handler) saveInput(ctx *gin.Context) {
	var input model.Inputs
	err := ctx.Bind(&input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "error no valid object")
		return
	}
	fmt.Sprintf("map: %v", input)
	orderId, err := fpr.SaveData(input)
	fmt.Println(orderId, err)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "err: "+err.Error())
	}
	ctx.JSON(http.StatusOK, fmt.Sprintf("Saved to DB and Cache! orderId: %v", orderId))
}

func (h *Handler) getInput(ctx *gin.Context) {
	data, err := fpr.GetData()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	}
	ctx.JSON(http.StatusOK, data)
}
func (h *Handler) reorderInput(ctx *gin.Context) {
	var input model.ReorderInput
	err := ctx.BindJSON(&input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "error no valid Reorder object")
		return
	}
	data, err := fpr.ReorderInputs(input)
	if err != nil {
		return
	}
	ctx.JSON(http.StatusOK, data)
}
