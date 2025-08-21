package models

import (
	"adire-apparel/lib"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderStatus string

const (
	Pending    OrderStatus = "pending"
	Processing OrderStatus = "processing"
	Shipped    OrderStatus = "shipped"
	Delivered  OrderStatus = "delivered"
	Cancelled  OrderStatus = "cancelled"
)

type OrderItemModel struct {
	BaseModel
	OrderId   uuid.UUID     `json:"order_id" gorm:"not null"`
	Order     OrderModel    `json:"order,omitempty" gorm:"foreignKey:OrderId"`
	ProductId string        `json:"product_id" gorm:"not null"`
	Product   ProductModel  `json:"product" gorm:"foreignKey:ProductId"`
	VariantId *string       `json:"variant_id,omitempty" gorm:"default:null"`
	Variant   *VariantModel `json:"variant,omitempty" gorm:"foreignKey:VariantId"`
	Quantity  int           `json:"quantity" gorm:"not null;check:quantity > 0" validate:"required,min=1"`
	Price     PricingModel  `json:"price" gorm:"embedded;embeddedPrefix:price_" validate:"required"`
}

type OrderModel struct {
	BaseModel
	UserId          string           `json:"user_id" gorm:"not null"`
	User            UserModel        `json:"user" gorm:"foreignKey:UserId"`
	Items           []OrderItemModel `json:"items" gorm:"foreignKey:OrderId"`
	Currency        string           `json:"currency" gorm:"type:varchar(3);default:NGN" validate:"required,len=3"`
	OrderDate       time.Time        `json:"order_date" gorm:"not null;default:CURRENT_TIMESTAMP"`
	Status          OrderStatus      `json:"status" gorm:"type:varchar(20);default:pending" validate:"required,oneof=pending processing shipped delivered cancelled"`
	SubtotalAmount  int              `json:"subtotal_amount" gorm:"not null;check:subtotal_amount >= 0" validate:"required,min=0"`
	TotalAmount     int              `json:"total_amount" gorm:"not null;check:total_amount >= 0" validate:"required,min=0"`
	ShippingAddress Address          `json:"shipping_address" gorm:"embedded;embeddedPrefix:shipping_"`
	BillingAddress  Address          `json:"billing_address" gorm:"embedded;embeddedPrefix:billing_"`
	TrackingNumber  string           `json:"tracking_number,omitempty" gorm:"type:varchar(100)"`
	ShippedAt       *time.Time       `json:"shipped_at,omitempty"`
	DeliveredAt     *time.Time       `json:"delivered_at,omitempty"`
	CancelledAt     *time.Time       `json:"cancelled_at,omitempty"`
	Notes           string           `json:"notes,omitempty" gorm:"type:text"`
}

type Address struct {
	Street  string  `json:"street" gorm:"column:street"`
	City    string  `json:"city" gorm:"column:city"`
	State   string  `json:"state" gorm:"column:state"`
	ZipCode *string `json:"zip_code" gorm:"column:zip_code"`
	Country string  `json:"country" gorm:"column:country;default:US"`
}

func (o *OrderModel) BeforeCreate(tx *gorm.DB) error {
	if o.OrderDate.IsZero() {
		o.OrderDate = time.Now()
	}
	if o.Status == "" {
		o.Status = Pending
	}
	if o.Currency == "" {
		o.Currency = "NGN"
	}
	o.TrackingNumber = lib.GenerateTrackingNumber()

	o.calculateAmounts()

	return nil
}

func (o *OrderModel) BeforeUpdate(tx *gorm.DB) error {
	now := time.Now()

	if o.Status == Shipped && o.ShippedAt == nil {
		o.ShippedAt = &now
	}
	if o.Status == Delivered && o.DeliveredAt == nil {
		o.DeliveredAt = &now
	}
	if o.Status == Cancelled && o.CancelledAt == nil {
		o.CancelledAt = &now
	}

	return nil
}

func (o *OrderModel) calculateAmounts() {
	subtotal := 0
	total := 0

	for _, item := range o.Items {
		itemSubtotal := item.Price.Base * item.Quantity

		if item.Price.DiscountPercentage > 0 {
			discount := (itemSubtotal * item.Price.DiscountPercentage) / 100
			itemSubtotal -= discount
		}

		subtotal += itemSubtotal

		itemTotal := itemSubtotal
		if item.Price.TaxRate > 0 {
			tax := (itemSubtotal * item.Price.TaxRate) / 100
			itemTotal += tax
		}

		total += itemTotal
	}

	o.SubtotalAmount = subtotal
	o.TotalAmount = total
}

func (o *OrderModel) CalculateTotal() int {
	total := 0
	for _, item := range o.Items {
		itemTotal := item.Price.Base * item.Quantity

		if item.Price.DiscountPercentage > 0 {
			discount := (itemTotal * item.Price.DiscountPercentage) / 100
			itemTotal -= discount
		}

		if item.Price.TaxRate > 0 {
			tax := (itemTotal * item.Price.TaxRate) / 100
			itemTotal += tax
		}

		total += itemTotal
	}
	return total
}

func (o *OrderModel) CanBeCancelled() bool {
	return o.Status == Pending || o.Status == Processing
}

func (o *OrderModel) CanBeShipped() bool {
	return o.Status == Processing
}

func (oi *OrderItemModel) GetSubtotal() int {
	subtotal := oi.Price.Base * oi.Quantity

	if oi.Price.DiscountPercentage > 0 {
		discount := (subtotal * oi.Price.DiscountPercentage) / 100
		subtotal -= discount
	}

	if oi.Price.TaxRate > 0 {
		tax := (subtotal * oi.Price.TaxRate) / 100
		subtotal += tax
	}

	return subtotal
}
