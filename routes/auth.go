package routes

import (
	"adire-apparel/handlers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.RouterGroup) *gin.RouterGroup {
	auth := router.Group("/auth")
	handler := handlers.NewAuthHandler()

	auth.POST("/signup", handler.Signup())
	auth.POST("/signin", handler.Signin())
	auth.POST("/refresh", handler.Refresh())
	auth.POST("/verify-email", handler.VerifyEmail())
	auth.POST("/verify-phone", handler.VerifyPhone())
	auth.POST("/change-password", handler.ChangePassword())
	auth.POST("/forgot-password", handler.ForgotPassword())
	auth.POST("/reset-password", handler.ResetPassword())
	auth.GET("/facebook", handler.SigninWithOauth())
	auth.GET("/google", handler.SigninWithOauth())

	return auth
}
