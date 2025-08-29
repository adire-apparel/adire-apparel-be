package dto

import (
	"adire-apparel/models"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type OrderPagination struct {
	Pagination
	Name *string `json:"name,omitempty"`
}

type AddressDto struct {
	Street  string  `json:"street" validate:"required,min=1,max=255"`
	City    string  `json:"city" validate:"required,min=1,max=100"`
	State   string  `json:"state" validate:"required,min=1,max=100"`
	ZipCode *string `json:"zip_code" validate:"omitempty,min=3,max=20"`
	Country string  `json:"country" validate:"required,len=2"`
}

type PricingDto struct {
	Base               int  `json:"base" validate:"required,min=0"`
	TaxRate            int  `json:"tax_rate" validate:"min=0,max=10000"`
	DiscountPercentage *int `json:"discount_percentage" validate:"omitempty,min=0,max=100"`
}

type CreateOrderItemDto struct {
	ProductId string     `json:"product_id" validate:"required,uuid"`
	VariantId *string    `json:"variant_id,omitempty" validate:"omitempty,uuid"`
	Quantity  int        `json:"quantity" validate:"required,min=1,max=1000"`
	Price     PricingDto `json:"price" validate:"required"`
}

type CreateOrderDto struct {
	Items           []CreateOrderItemDto `json:"items" validate:"required,min=1,max=50,dive"`
	Currency        string               `json:"currency" validate:"len=3"`
	ShippingAddress *AddressDto          `json:"shipping_address" validate:"required"`
	BillingAddress  *AddressDto          `json:"billing_address" validate:"required"`
	UserId          *string              `json:"user_id" validate:"required,uuid"`
	Notes           *string              `json:"notes,omitempty" validate:"max=1000"`
}

type UpdateOrderItemDto struct {
	Id        string     `json:"id,omitempty" validate:"omitempty,uuid"`
	ProductId string     `json:"product_id" validate:"required,uuid"`
	VariantId *string    `json:"variant_id,omitempty" validate:"omitempty,uuid"`
	Quantity  int        `json:"quantity" validate:"required,min=1,max=1000"`
	Price     PricingDto `json:"price" validate:"required"`
}

type UpdateOrderDto struct {
	Status          *models.OrderStatus  `json:"status,omitempty" validate:"omitempty,oneof=pending processing shipped delivered cancelled"`
	Items           []UpdateOrderItemDto `json:"items,omitempty" validate:"omitempty,min=1,max=50,dive"`
	Currency        *string              `json:"currency,omitempty" validate:"omitempty,len=3"`
	ShippingAddress *AddressDto          `json:"shipping_address,omitempty" validate:"omitempty"`
	BillingAddress  *AddressDto          `json:"billing_address,omitempty" validate:"omitempty"`
	Notes           *string              `json:"notes,omitempty" validate:"omitempty,max=1000"`
}

type UpdateOrderStatusDto struct {
	Status models.OrderStatus `json:"status" validate:"required,oneof=pending processing shipped delivered cancelled"`
	Notes  *string            `json:"notes,omitempty" validate:"omitempty,max=1000"`
}

type OrderResponseDto struct {
	Id              uuid.UUID              `json:"id"`
	UserId          uuid.UUID              `json:"user_id"`
	Items           []OrderItemResponseDto `json:"items"`
	Currency        string                 `json:"currency"`
	OrderDate       time.Time              `json:"order_date"`
	Status          models.OrderStatus     `json:"status"`
	SubtotalAmount  int                    `json:"subtotal_amount"`
	TotalAmount     int                    `json:"total_amount"`
	ShippingAddress AddressDto             `json:"shipping_address"`
	BillingAddress  AddressDto             `json:"billing_address"`
	TrackingNumber  *string                `json:"tracking_number,omitempty"`
	ShippedAt       *time.Time             `json:"shipped_at,omitempty"`
	DeliveredAt     *time.Time             `json:"delivered_at,omitempty"`
	CancelledAt     *time.Time             `json:"cancelled_at,omitempty"`
	Notes           *string                `json:"notes,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

type OrderItemResponseDto struct {
	Id        uuid.UUID  `json:"id"`
	ProductId uuid.UUID  `json:"product_id"`
	VariantId *uuid.UUID `json:"variant_id,omitempty"`
	Quantity  int        `json:"quantity"`
	Price     PricingDto `json:"price"`
	Subtotal  int        `json:"subtotal"`
}

type OrderListResponseDto struct {
	Id             uuid.UUID          `json:"id"`
	UserId         uuid.UUID          `json:"user_id"`
	User           models.UserModel   `json:"user"`
	Currency       string             `json:"currency"`
	OrderDate      time.Time          `json:"order_date"`
	Status         models.OrderStatus `json:"status"`
	TotalAmount    int                `json:"total_amount"`
	ItemCount      int                `json:"item_count"`
	TrackingNumber *string            `json:"tracking_number,omitempty"`
	CreatedAt      time.Time          `json:"created_at"`
}

type OrderFilterDto struct {
	Status   *models.OrderStatus `json:"status,omitempty" validate:"omitempty,oneof=pending processing shipped delivered cancelled"`
	UserId   *string             `json:"user_id,omitempty" validate:"omitempty,uuid"`
	DateFrom *time.Time          `json:"date_from,omitempty"`
	DateTo   *time.Time          `json:"date_to,omitempty"`
	Currency *string             `json:"currency,omitempty" validate:"omitempty,len=3"`
}

type OrderStatsDto struct {
	TotalOrders      int            `json:"total_orders"`
	TotalRevenue     int            `json:"total_revenue"`
	AverageOrderSize int            `json:"average_order_size"`
	OrdersByStatus   map[string]int `json:"orders_by_status"`
}

func (dto *CreateOrderDto) SetDefaults() {
	if dto.Currency == "" {
		dto.Currency = "NGN"
	}
	if dto.ShippingAddress.Country == "" {
		dto.ShippingAddress.Country = "NG"
	}
	if dto.BillingAddress.Country == "" {
		dto.BillingAddress.Country = "NG"
	}
}

func (dto *CreateOrderDto) Validate() error {
	if len(dto.Items) == 0 {
		return errors.New("order must contain at least one item")
	}

	for i, item := range dto.Items {
		if item.Quantity <= 0 {
			return fmt.Errorf("item %d: quantity must be greater than 0", i+1)
		}
		if item.Price.Base < 0 {
			return fmt.Errorf("item %d: price cannot be negative", i+1)
		}
	}

	return nil
}
