package middlewares

import (
	"adire-apparel/config"
	"adire-apparel/database"
	"adire-apparel/lib"
	"adire-apparel/models"
	"adire-apparel/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RoleMiddleware(role models.Role) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, exists := ctx.Get(config.AppConfig.CurrentUserId)
		if !exists {
			ctx.Error(lib.NewApiErrror("You need to login to access this endpoint", http.StatusUnauthorized))
			ctx.Abort()
			return
		}

		service := services.NewUserService(database.GetDatabase())
		user, err := service.GetUserById(id.(string))
		if err != nil {
			ctx.Error(lib.NewApiErrror("Invalid user id", http.StatusUnauthorized))
			ctx.Abort()
			return
		}

		if string(user.Role) != string(role) {
			ctx.Error(lib.NewApiErrror("Invalid user type. You're not authorized to view this endpoint.", http.StatusUnauthorized))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
