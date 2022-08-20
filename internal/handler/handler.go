package handler

import (
	"dealer/internal/models"
	"dealer/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	orderProcessor service.OrderProcessorInterface
}

func NewHandler(orderProcessor service.OrderProcessorInterface) *Handler {
	return &Handler{
		orderProcessor: orderProcessor,
	}
}

func (h *Handler) NewOrder(ctx *gin.Context) {
	var req *models.OrderRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	order := &models.Order{
		OrderType:      req.OrderType,
		Quantity:       req.Quantity,
		RemainQuantity: req.Quantity,
		PriceType:      req.PriceType,
		Price:          req.Price,
	}
	err := h.orderProcessor.NewOrder(ctx, order)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, order)
}

func (h *Handler) CancelOrder(ctx *gin.Context) {
	var req *models.CancelOrderRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	err := h.orderProcessor.CancelOrder(ctx, req.ID)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}
