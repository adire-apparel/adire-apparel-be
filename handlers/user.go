package handlers

import (
	"adire-apparel/database"
	"adire-apparel/dto"
	"adire-apparel/lib"
	"adire-apparel/services"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		service: services.NewUserService(database.GetDatabase()),
	}
}

func (h UserHandler) UpdateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload dto.EditUserDto
		id := ctx.Param("id")

		if err := ctx.ShouldBind(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		user := dto.EditUserDto{
			Email: payload.Email,
			Name:  payload.Name,
			Phone: payload.Phone,
		}

		updated, err := h.service.UpdateUser(id, &user)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error, "+err.Error())
			return
		}

		lib.Success(ctx, "User updated successfully", updated)
	}
}

func (h UserHandler) UpdateUserPreferences() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload dto.EditUserPreferencesDto
		id := ctx.Param("id")

		if err := ctx.ShouldBind(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		user := dto.EditUserPreferencesDto{
			EmailNotifications:     payload.EmailNotifications,
			PushNotifications:      payload.PushNotifications,
			SMSNotifications:       payload.SMSNotifications,
			MarketingEmails:        payload.MarketingEmails,
			NewsletterSubscription: payload.NewsletterSubscription,
			Theme:                  payload.Theme,
			DisplayCurrency:        payload.DisplayCurrency,
		}

		updated, err := h.service.UpdateUserPreferences(id, user)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "User updated successfully", updated)
	}
}

func (h UserHandler) GetUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var params dto.UserPagination

		if err := ctx.ShouldBind(&params); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		users, err := h.service.GetAllUsers(params)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Users fetched successfully", users)
	}
}

func (h UserHandler) GetUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		user, err := h.service.GetUserById(id)

		if err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		lib.Success(ctx, "User fetched successfully", user)
	}
}

func (h UserHandler) DeleteUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		err := h.service.DeleteUser(id)

		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "User deleted successfully", nil)
	}
}

func (h UserHandler) UploadImage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		const (
			maxMemory = 8 << 20 // 8 MB
			bucket    = "users"
		)

		contentType := ctx.ContentType()
		if !strings.Contains(contentType, "multipart/form-data") {
			lib.BadRequest(ctx, "Invalid content type", "400")
			return
		}

		if err := ctx.Request.ParseMultipartForm(maxMemory); err != nil {
			lib.BadRequest(ctx, "Failed to parse multipart form: "+err.Error(), "400")
			return
		}

		form, err := ctx.MultipartForm()
		if err != nil {
			lib.BadRequest(ctx, "Failed to parse multipart form: "+err.Error(), "400")
			return
		}

		file := form.File["image"][0]
		if file == nil {
			lib.BadRequest(ctx, "No images provided in the request", "400")
			return
		}

		imageUrl, err := lib.SingleImageUploader(file, bucket)
		if err != nil {
			lib.InternalServerError(ctx, "Unable to upload image")
			return
		}

		id := ctx.Param("id")
		err = h.service.UploadImage(id, imageUrl)

		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Image uploaded successfully", nil)
	}
}
