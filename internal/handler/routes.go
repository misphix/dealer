package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router gin.IRouter, handler *Handler) {
	router.GET("status", status)
	v1Group := router.Group("v1")
	order := v1Group.Group("order")
	order.POST("", handler.NewOrder)
}

func status(ctx *gin.Context) {
	ctx.String(http.StatusOK, "ok")
}
