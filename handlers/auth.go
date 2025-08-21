package handlers

import (
	"adire-apparel/config"
	"adire-apparel/database"
	"adire-apparel/dto"
	"adire-apparel/lib"
	"adire-apparel/services"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/oov/gothic"
)

type AuthHandler struct {
	service services.AuthService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		service: *services.NewAuthService(database.GetDatabase()),
	}
}

func (h *AuthHandler) Signup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload dto.CreateUserDto

		if err := ctx.ShouldBind(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		user := dto.CreateUserDto{
			Name:     payload.Name,
			Email:    payload.Email,
			Role:     payload.Role,
			Password: payload.Password,
		}

		err := h.service.Signup(user)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Account created successfully. Check your mail for the confirmation message", nil)
	}
}

func (h *AuthHandler) Signin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload dto.SigninDto

		if err := ctx.ShouldBind(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		user, err := h.service.Signin(payload)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "User signed in successfully.", user)
	}
}

func (h *AuthHandler) SigninWithOauth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := gothic.BeginAuth(ctx.Param("provider"), ctx.Writer, ctx.Request)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
			return
		}
	}
}

func (h *AuthHandler) SigninWithOauthCallback() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		provider := ctx.Param("provider")
		token := ctx.Query("token")
		var user *goth.User

		if provider != "google" {
			req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
			if err != nil {
				lib.InternalServerError(ctx, "Internal server error,"+err.Error())
				return
			}
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			client := &http.Client{}
			response, err := client.Do(req)
			if err != nil {
				lib.InternalServerError(ctx, "Internal server error,"+err.Error())
				return
			}

			data, err := io.ReadAll(response.Body)
			if err != nil {
				lib.InternalServerError(ctx, "Internal server error,"+err.Error())
				return
			}

			var result dto.GoogleModel
			if err := json.Unmarshal(data, &result); err != nil {
				lib.InternalServerError(ctx, "Internal server error,"+err.Error())
				return
			}

			user = &goth.User{
				Provider: provider,
				Name:     result.Name,
				Email:    result.Email,
				UserID:   result.UserID,
			}
		}

		newUser, err := h.service.SigninWithOauth(user)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		setAuthToken(ctx, newUser.Token)
		lib.Success(ctx, "User signed in successfully.", newUser)
	}
}

func (h *AuthHandler) Refresh() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.GetString(config.AppConfig.CurrentUserId)

		token, err := h.service.Refresh(id)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
		}

		lib.Success(ctx, "Check your mail for instructions to rest your password", token)
	}
}

func (h *AuthHandler) VerifyEmail() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload dto.VerifyEmailDto

		if err := ctx.ShouldBind(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		err := h.service.VerifyEmail(payload)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Password changed successfully.", nil)
	}
}

func (h *AuthHandler) VerifyPhone() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload dto.VerifyPhoneDto

		if err := ctx.ShouldBind(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		err := h.service.VerifyPhone(payload)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Password changed successfully.", nil)
	}
}

func (h *AuthHandler) ChangePassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload dto.ChangePasswordDto
		id := ctx.GetString(config.AppConfig.CurrentUserId)

		if err := ctx.ShouldBind(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		err := h.service.ChangePassword(id, payload)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Password changed successfully.", nil)
	}
}

func (h *AuthHandler) ForgotPassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		lib.Success(ctx, "Check your mail for instructions to rest your password", nil)
	}
}

func (h *AuthHandler) ResetPassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		lib.Success(ctx, "Your password has been reset. Proceed to sign in.", nil)
	}
}

func setAuthToken(ctx *gin.Context, token string) {
	cookie := &http.Cookie{Name: string(config.AppConfig.JWTSecret), Value: token, Path: "/"}
	http.SetCookie(ctx.Writer, cookie)
}
