package routes

import (
	"adire-apparel/handlers"

	"github.com/gin-gonic/gin"
)

func TransactionRoutes(router *gin.RouterGroup) *gin.RouterGroup {
	transaction := router.Group("/transactions")
	handler := handlers.NewTransactionHandler()

	transaction.POST("", handler.CreateTransaction())
	transaction.GET("", handler.GetTransactions())
	transaction.GET("/user/:id", handler.GetTransactions())
	transaction.GET("/admin/:id", handler.GetTransactions())
	transaction.GET("/:id", handler.GetTransaction())
	transaction.PUT("/:id", handler.UpdateTransaction())
	transaction.DELETE("/:id", handler.DeleteTransaction())

	return transaction
}
