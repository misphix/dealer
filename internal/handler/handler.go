package handler

import (
	"dealer/internal/models"
	"dealer/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	dealer service.DealerInterface
}

func NewHandler(dealer service.DealerInterface) *Handler {
	return &Handler{
		dealer: dealer,
	}
}

func (h *Handler) NewOrder(ctx *gin.Context) {
	var order *models.Order
	if err := ctx.ShouldBind(&order); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	deals, err := h.dealer.ProcessOrder(ctx, order)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, deals)
}
