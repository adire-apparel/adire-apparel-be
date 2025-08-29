package dto

import (
	"adire-apparel/models"
)

type ProductPagination struct {
	Pagination
	Name       *string  `json:"name"`
	CategoryId *string  `json:"category_id" form:"category_id"`
	MinPrice   *float64 `json:"min_price" form:"min_price"`
	MaxPrice   *float64 `json:"max_price" form:"max_price"`
	SortBy     *string  `json:"sort_by" form:"sort_by"`
	SortOrder  *string  `json:"sort_order" form:"sort_order"`
}

type CreateProductDto struct {
	Name        string                `json:"name" validate:"required"`
	Description string                `json:"description" validate:"required"`
	CategoryId  string                `json:"category_id" validate:"required"`
	Price       models.PricingModel   `json:"price" validate:"required"`
	Variants    []models.VariantModel `json:"variants"`
	CreatedBy   string                `json:"created_by"`
}

type UpdateProductDto struct {
	Name        string                `json:"name"`
	Description string                `json:"description"`
	CategoryId  string                `json:"category_id"`
	Price       models.PricingModel   `json:"price"`
	Variants    []models.VariantModel `json:"variants"`
	Status      models.ProductStatus  `json:"status"`
}

type RemoveImagesDto struct {
	Images []string `json:"images"`
}
