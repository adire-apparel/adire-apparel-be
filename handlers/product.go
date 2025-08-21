package handlers

import (
	"adire-apparel/constants"
	"adire-apparel/database"
	"adire-apparel/dto"
	"adire-apparel/lib"
	"adire-apparel/models"
	"adire-apparel/services"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	service *services.ProductService
}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{
		service: services.NewProductService(database.GetDatabase()),
	}
}

var (
	ErrInvalidColorName    = errors.New("invalid color name")
	ErrInvalidColorHex     = errors.New("invalid color hex")
	ErrInvalidSizes        = errors.New("invalid sizes")
	ErrInvalidBasePrice    = errors.New("base price is required")
	ErrInvalidTaxRate      = errors.New("tax rate is required")
	ErrInvalidCurrency     = errors.New("currency is required")
	ErrInvalidCurrencyType = errors.New("invalid currency is required")
)

func ValidateVariants(variants []models.VariantModel) error {
	for _, variant := range variants {
		if variant.Color.Name == "" {
			return ErrInvalidColorName
		}
		if variant.Color.Hex == "" {
			return ErrInvalidColorHex
		}
		if len(variant.Sizes) == 0 {
			return ErrInvalidSizes
		}
		for _, size := range variant.Sizes {
			if size == 0 {
				return ErrInvalidSizes
			}
		}
	}
	return nil
}

func ValidatePrice(price models.PricingModel) error {
	if price.Base == 0 {
		return ErrInvalidBasePrice
	}
	if price.TaxRate == 0 {
		return ErrInvalidTaxRate
	}
	if !constants.IsValidCurrency(price.Currency) {
		return ErrInvalidCurrency
	}
	return nil
}

func (h *ProductHandler) CreateProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload dto.CreateProductDto

		if err := ctx.ShouldBind(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		if err := ValidateVariants(payload.Variants); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		if err := ValidatePrice(payload.Price); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		product := dto.CreateProductDto{
			Name:        payload.Name,
			Description: payload.Description,
			CategoryId:  payload.CategoryId,
			Price:       payload.Price,
			Variants:    payload.Variants,
			CreatedBy:   payload.CreatedBy,
		}

		saved, err := h.service.CreateProduct(product)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Product created successfully", saved)
	}
}

func (h *ProductHandler) UpdateProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload dto.UpdateProductDto
		id := ctx.Param("id")

		if err := ctx.ShouldBind(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		if len(payload.Variants) == 0 {
			lib.BadRequest(ctx, "at least one variant is required", "400")
			return
		}

		if err := ValidateVariants(payload.Variants); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		if err := ValidatePrice(payload.Price); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		product := dto.UpdateProductDto{
			Name:        payload.Name,
			Description: payload.Description,
			CategoryId:  payload.CategoryId,
			Price:       payload.Price,
			Variants:    payload.Variants,
			Status:      payload.Status,
		}

		updated, err := h.service.UpdateProduct(id, product)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Product updated successfully", updated)
	}
}

func (h *ProductHandler) GetProducts() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var params dto.ProductPagination

		if err := ctx.ShouldBind(&params); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		products, err := h.service.GetProducts(params)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Products retreived successfully", products)
	}
}

func (h *ProductHandler) GetProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		product, err := h.service.GetProduct(id)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Product retrieved successfully", product)
	}
}

func (h *ProductHandler) DeleteProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		err := h.service.DeleteProduct(id)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Product deleted successfully", nil)
	}
}

func (h *ProductHandler) AddImagesToProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		const (
			maxMemory = 8 << 20 // 8 MB
			bucket    = "products"
		)

		contentType := ctx.ContentType()
		if !strings.Contains(contentType, "multipart/form-data") {
			lib.BadRequest(ctx, "invalid content type", "400")
			return
		}

		if err := ctx.Request.ParseMultipartForm(maxMemory); err != nil {
			lib.BadRequest(ctx, "Failed to parse multipart form: "+err.Error(), "400")
			return
		}

		form, err := ctx.MultipartForm()
		if err != nil {
			lib.BadRequest(ctx, "Failed to get multipart form: "+err.Error(), "400")
			return
		}

		files := form.File["images"]
		if len(files) == 0 {
			lib.BadRequest(ctx, "No images provided in the request", "400")
			return
		}

		imageUrls, err := lib.ImageUploader(files, bucket)
		if err != nil {
			lib.BadRequest(ctx, "Unable to upload images", "400")
			return
		}

		id := ctx.Param("id")
		if err := h.service.AddImagesToProduct(id, imageUrls); err != nil {
			lib.BadRequest(ctx, "Failed to add images to product", "400")
			return
		}

		lib.Success(ctx, "Images added successfully", nil)
	}
}

func (h *ProductHandler) RemoveImagesFromProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var images dto.RemoveImagesDto
		id := ctx.Param("id")

		if err := ctx.ShouldBind(&images); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		err := h.service.RemoveImagesFromProduct(id, images)
		if err != nil {
			lib.InternalServerError(ctx, "Failed to upload images: "+err.Error())
			return
		}

		lib.Success(ctx, "Image removed successfully", nil)
	}
}
