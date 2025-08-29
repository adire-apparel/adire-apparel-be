package handlers

import (
	"adire-apparel/database"
	"adire-apparel/dto"
	"adire-apparel/lib"
	"adire-apparel/services"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	service *services.OrderService
}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{
		service: services.NewOrderService(database.GetDatabase()),
	}
}

func (h *OrderHandler) CreateOrder() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload dto.CreateOrderDto

		if err := ctx.ShouldBind(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		order, err := h.service.CreateOrder(payload)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Order created successfully", order)
	}
}

func (h *OrderHandler) UpdateOrder() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload dto.UpdateOrderDto

		if err := ctx.ShouldBind(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		lib.Success(ctx, "Order updated successfully", nil)
	}
}

func (h *OrderHandler) GetOrders() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var params dto.Pagination

		if err := ctx.ShouldBind(&params); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		orders, err := h.service.GetOrders(params)

		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Orders retrieved successfully", orders)
	}
}

func (h *OrderHandler) GetOrdersByUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var params dto.Pagination

		if err := ctx.ShouldBind(&params); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		id := ctx.Param("id")
		orders, err := h.service.GetOrdersByUser(id, params)

		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Orders retrieved successfully", orders)
	}
}

func (h *OrderHandler) GetOrder() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		order, err := h.service.GetOrder(id)
		if err != nil {
			ctx.Error(err)
			return
		}

		lib.Success(ctx, "Order retrieved successfully", order)
	}
}

func (h *OrderHandler) DeleteOrder() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		err := h.service.DeleteOrder(id)
		if err != nil {
			ctx.Error(err)
			return
		}

		lib.Success(ctx, "Order deleted successfully", nil)
	}
}
