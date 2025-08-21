package handlers

import (
	"adire-apparel/database"
	"adire-apparel/dto"
	"adire-apparel/lib"
	"adire-apparel/services"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	service *services.CategoryService
}

func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{
		service: services.NewCategoryService(database.GetDatabase()),
	}
}

func (h *CategoryHandler) CreateCategory() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload dto.CreateCategoryDto

		if err := ctx.ShouldBind(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		category := dto.CreateCategoryDto{
			Name:        payload.Name,
			Description: payload.Description,
			CreatedBy:   payload.CreatedBy,
		}

		saved, err := h.service.CreateCategory(category)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Category created successfully", saved)
	}
}

func (h *CategoryHandler) UpdateCategory() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload dto.CreateCategoryDto
		id := ctx.Param("id")

		if err := ctx.ShouldBind(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		category, err := h.service.UpdateCategory(id, payload)

		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Category updated successfully", category)
	}
}

func (h *CategoryHandler) GetCategories() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var params dto.Pagination

		if err := ctx.ShouldBind(&params); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		categories, err := h.service.GetCategories(params)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Categories retrived successfully", categories)
	}
}

func (h *CategoryHandler) GetCategory() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		category, err := h.service.GetCategory(id)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Category retrieved successfully", category)
	}
}

func (h *CategoryHandler) DeleteCategory() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		err := h.service.DeleteCategory(id)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Category deleted successfully", nil)
	}
}
