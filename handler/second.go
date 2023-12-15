package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"redis-task/database"
	"redis-task/model"
)

var sr = database.NewSecondRedis()

func (h *Handler) saveStruct(ctx *gin.Context) {
	var input model.User
	err := ctx.BindJSON(&input)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, "error no valid values")
		return
	}
	err = sr.SaveStructToCache(input)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, "saved")
	return
}

func (h *Handler) getStruct(ctx *gin.Context) {
	data, err := sr.GetStructFromCache()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, map[string]interface{}{"data": data})
}
func (h *Handler) updateStruct(ctx *gin.Context) {
	var updateInput model.UserUpdate
	var id = ctx.Param("id")
	err := ctx.BindJSON(&updateInput)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}
	err = sr.UpdateStructOnCache(updateInput, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, map[string]interface{}{"data": "updated"})
	return
}
