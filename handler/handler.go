package handler

import (
	"github.com/gin-gonic/gin"
)

type Handler struct {
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	first := router.Group("first")
	{
		first.GET("", h.getInput)
		first.POST("", h.saveInput)
		first.PATCH("", h.reorderInput)
	}
	second := router.Group("/second")
	{
		second.GET("", h.getStruct)
		second.POST("", h.saveStruct)
		second.PATCH(":id", h.updateStruct)
	}
	return router
}
