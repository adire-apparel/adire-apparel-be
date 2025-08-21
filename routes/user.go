package routes

import (
	"adire-apparel/handlers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.RouterGroup) *gin.RouterGroup {
	user := router.Group("/users")
	handler := handlers.NewUserHandler()

	user.GET("", handler.GetUsers())
	user.GET("/:id", handler.GetUser())
	user.PUT("/:id", handler.UpdateUser())
	user.PUT("/preferences/:id", handler.UpdateUserPreferences())
	user.PUT("/image/:id", handler.UploadImage())
	user.DELETE("/:id", handler.DeleteUser())

	return user
}
