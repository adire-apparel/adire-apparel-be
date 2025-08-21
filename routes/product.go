package routes

import (
	"adire-apparel/handlers"

	"github.com/gin-gonic/gin"
)

func ProductRoutes(router *gin.RouterGroup) *gin.RouterGroup {
	product := router.Group("/products")
	handler := handlers.NewProductHandler()

	product.POST("", handler.CreateProduct())
	product.GET("", handler.GetProducts())
	product.GET("/:id", handler.GetProduct())
	product.PUT("/:id", handler.UpdateProduct())
	product.PUT("/add-images/:id", handler.AddImagesToProduct())
	product.PUT("/remove-images/:id", handler.RemoveImagesFromProduct())
	product.DELETE("/:id", handler.DeleteProduct())

	return product
}
