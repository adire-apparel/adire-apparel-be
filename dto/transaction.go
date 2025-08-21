package dto

import (
	"adire-apparel/models"
)

type CreateTransactionDto struct {
	UserId      string                     `json:"user_id" gorm:"type:uuid;not null"`
	Type        models.TransactionType     `json:"type" gorm:"type:varchar(20);not null"`
	Status      models.TransactionStatus   `json:"status" gorm:"type:varchar(20);default:pending"`
	Amount      int64                      `json:"amount" gorm:"not null"`
	Currency    string                     `json:"currency" gorm:"not null;default:NGN"`
	Description string                     `json:"description"`
	OrderId     string                     `json:"order_id"`
	Metadata    models.TransactionMetadata `json:"metadata" gorm:"type:json"`
	CreatedBy   string                     `json:"created_by" gorm:"column:created_by;not null" validate:"required"`
}

type UpdateTransactionDto struct {
	Status   models.TransactionStatus   `json:"status" gorm:"type:varchar(20);default:pending"`
	Metadata models.TransactionMetadata `json:"metadata" gorm:"type:json"`
}
