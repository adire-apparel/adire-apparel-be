package routes

import (
	"adire-apparel/handlers"

	"github.com/gin-gonic/gin"
)

func CategoryRoutes(router *gin.RouterGroup) *gin.RouterGroup {
	product := router.Group("/categories")
	handler := handlers.NewCategoryHandler()

	product.POST("", handler.CreateCategory())
	product.GET("", handler.GetCategories())
	product.GET("/:id", handler.GetCategory())
	product.PUT("/:id", handler.UpdateCategory())
	product.DELETE("/:id", handler.DeleteCategory())

	return product
}
