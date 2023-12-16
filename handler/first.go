package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"redis-task/database"
	"redis-task/model"
)

var fdb = new(database.FirstPostgres)

func (h *Handler) saveInput(ctx *gin.Context) {
	var input model.Inputs
	err := ctx.Bind(&input)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, "error no valid object")
		return
	}
	orderId, err := fdb.SaveData(input)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, "err: "+err.Error())
		return
	}
	ctx.JSON(http.StatusOK, fmt.Sprintf("Saved to DB and Cache! orderId: %v", orderId))
	return
}

func (h *Handler) getInput(ctx *gin.Context) {
	data, err := fdb.GetData()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, data)
	return
}
func (h *Handler) reorderInput(ctx *gin.Context) {
	var input model.ReorderInput
	err := ctx.BindJSON(&input)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]interface{}{"error": "error no valid Reorder object"})
		return
	}
	data, err := fdb.ReorderInput(input)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, data)
	return
}
