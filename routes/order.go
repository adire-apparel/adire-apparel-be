package routes

import (
	"adire-apparel/handlers"

	"github.com/gin-gonic/gin"
)

func OrderRoutes(router *gin.RouterGroup) *gin.RouterGroup {
	order := router.Group("/orders")
	handler := handlers.NewOrderHandler()

	order.POST("", handler.CreateOrder())
	order.PUT("/:id", handler.UpdateOrder())
	order.GET("", handler.GetOrders())
	order.GET("/user/:id", handler.GetOrdersByUser())
	order.GET("/:id", handler.GetOrder())
	order.DELETE("/:id", handler.DeleteOrder())

	return order
}
