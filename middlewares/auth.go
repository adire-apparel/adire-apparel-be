package middlewares

import (
	"adire-apparel/config"
	"adire-apparel/database"
	"adire-apparel/lib"
	"adire-apparel/models"
	"adire-apparel/services"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	bearerPrefix = "Bearer "
)

func isOpenRoute(route, method string) bool {
	for _, openRoute := range config.AppConfig.NonAuthRoutes {
		if openRoute.Endpoint == route && openRoute.Method == method {
			return true
		}
	}

	return false
}

func isAdminRoute(route, method string) bool {
	for _, adminRoute := range config.AppConfig.AdminRoutes {
		if adminRoute.Endpoint == route && adminRoute.Method == method {
			return true
		}
	}
	return false
}

func extractBearerToken(authHeader string) (string, bool) {
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", false
	}
	token := strings.TrimPrefix(authHeader, bearerPrefix)
	trimmedToken := strings.TrimSpace(token)
	return trimmedToken, trimmedToken != ""
}

func hasAdminAccess(userRole models.Role) bool {
	return userRole == "admin" || userRole == "staff"
}

func AuthMiddleware() gin.HandlerFunc {
	authService := services.NewAuthService(database.GetDatabase())

	return func(ctx *gin.Context) {
		route := ctx.FullPath()
		method := ctx.Request.Method
		if isOpenRoute(route, method) {
			ctx.Next()
			return
		}

		authHeader := ctx.Request.Header.Get("Authorization")
		token, ok := extractBearerToken(authHeader)
		if !ok {
			ctx.Error(lib.NewApiErrror("No auth token found", http.StatusUnauthorized))
			ctx.Abort()
			return
		}

		claims, err := lib.ValidateToken(token)
		if err != nil {
			ctx.Error(lib.NewApiErrror("Invalid auth token", http.StatusUnauthorized))
			ctx.Abort()
			return
		}

		user, err := authService.FindUserById(claims.UserId.String())
		if err != nil {
			ctx.Error(lib.NewApiErrror("User not found", http.StatusNotFound))
			ctx.Abort()
			return
		}

		if isAdminRoute(route, method) && !hasAdminAccess(user.Role) {
			ctx.Error(lib.NewApiErrror("Unauthorized access", http.StatusForbidden))
			ctx.Abort()
			return
		}

		log.Printf("userId:%s and user role:%s", user.Id, user.Role)

		// Store user info in context
		// To get userId elsewhere: ctx.GetString(config.AppConfig.CurrentUserId)
		// To get user elsewhere: ctx.Get(config.AppConfig.CurrentUser)
		// To get role elsewhere: ctx.GetString(config.AppConfig.CurrentUserRole)
		ctx.Set(config.AppConfig.CurrentUser, user)
		ctx.Set(config.AppConfig.CurrentUserId, user.Id.String())
		ctx.Set(config.AppConfig.CurrentUserRole, user.Role)

		ctx.Next()
	}
}
