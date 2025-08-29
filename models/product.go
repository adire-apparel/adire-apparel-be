package models

import (
	"adire-apparel/lib"
	"database/sql/driver"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

type ProductStatus string

const (
	Active       ProductStatus = "active"
	OutOfStock   ProductStatus = "out_of_stock"
	Discontinued ProductStatus = "discontinued"
	ComingSoon   ProductStatus = "coming_soon"
)

type ProductModel struct {
	BaseModel
	Name          string         `json:"name" gorm:"not null" validate:"required"`
	Description   string         `json:"description" gorm:"type:text" validate:"required"`
	CategoryId    string         `json:"category_id" gorm:"not null" validate:"required"`
	Category      CategoryModel  `json:"category" gorm:"foreignKey:CategoryId"`
	ImageUrls     StringArray    `json:"image_urls" gorm:"type:json" validate:"required"`
	Price         PricingModel   `json:"price" gorm:"embedded;embeddedPrefix:price_" validate:"required"`
	Sku           string         `json:"sku" gorm:"uniqueIndex;not null" validate:"required"`
	Slug          string         `json:"slug" gorm:"uniqueIndex;not null" validate:"required"`
	Status        ProductStatus  `json:"status" gorm:"type:varchar(20);default:active" validate:"required,oneof=active out_of_stock discontinued coming_soon"`
	Variants      []VariantModel `json:"variants" gorm:"foreignKey:ProductId"`
	CreatedBy     string         `json:"created_by" gorm:"column:created_by;not null" validate:"required"`
	CreatedByUser UserModel      `json:"created_by_user" gorm:"foreignKey:CreatedBy;references:Id"`
	IsDeleted     bool           `json:"is_deleted" gorm:"default:false"`
	Ratings       []RatingModel  `json:"ratings" gorm:"foreignKey:ProductId"`
}

type CategoryModel struct {
	BaseModel
	Name          string         `json:"name" gorm:"not null" validate:"required"`
	Description   string         `json:"description" gorm:"type:text" validate:"required"`
	CreatedBy     string         `json:"created_by" gorm:"column:created_by;not null" validate:"required"`
	CreatedByUser UserModel      `json:"created_by_user" gorm:"foreignKey:CreatedBy;references:Id"`
	Products      []ProductModel `json:"products,omitempty" gorm:"foreignKey:CategoryId"`
}

type VariantModel struct {
	BaseModel
	ProductId string       `json:"product_id" gorm:"not null"`
	Product   ProductModel `json:"-" gorm:"foreignKey:ProductId"`
	Color     Color        `json:"color" gorm:"embedded;embeddedPrefix:color_"`
	Sizes     IntArray     `json:"sizes" gorm:"type:json"`
	Stock     int          `json:"stock" gorm:"default:0"`
	ImageUrls StringArray  `json:"image_urls" gorm:"type:json" validate:"required"`
}

type Color struct {
	Name string `json:"name" gorm:"column:name"`
	Hex  string `json:"hex" gorm:"column:hex"`
}

type PricingModel struct {
	Base               int    `json:"base" gorm:"column:base"`
	Currency           string `json:"currency" gorm:"column:currency;default:NGN"`
	TaxRate            int    `json:"tax_rate" gorm:"column:tax_rate;default:0"`
	DiscountPercentage int    `json:"discount_percentage,omitempty" gorm:"column:discount_percentage;default:0"`
}

type RatingModel struct {
	BaseModel
	Rating    string       `json:"rating" gorm:"column:rating"`
	Comment   string       `json:"comment" gorm:"column:comment"`
	ProductId string       `json:"product_id" gorm:"column:product_id;not null" validate:"required"`
	Product   ProductModel `json:"product,omitempty" gorm:"foreignKey:ProductId;references:Id"`
	UserId    string       `json:"user_id" gorm:"column:user_id;not null" validate:"required"`
	User      UserModel    `json:"user" gorm:"foreignKey:UserId;references:Id"`
}

type StringArray []string

func (sa *StringArray) Scan(value interface{}) error {
	if value == nil {
		*sa = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, sa)
}

func (sa StringArray) Value() (driver.Value, error) {
	if len(sa) == 0 {
		return nil, nil
	}
	return json.Marshal(sa)
}

type IntArray []int

func (ia *IntArray) Scan(value interface{}) error {
	if value == nil {
		*ia = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, ia)
}

func (ia IntArray) Value() (driver.Value, error) {
	if len(ia) == 0 {
		return nil, nil
	}
	return json.Marshal(ia)
}

func (p *ProductModel) BeforeCreate(tx *gorm.DB) error {
	if p.Status == "" {
		p.Status = Active
		p.Sku = lib.CreateSku("ADP-")
		p.Slug = lib.ToKebabCase(p.Name)
	}
	return nil
}

func (p *ProductModel) BeforeUpdate(tx *gorm.DB) error {
	return nil
}
