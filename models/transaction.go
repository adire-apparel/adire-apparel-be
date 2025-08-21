package models

import (
	"adire-apparel/lib"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionStatus string
type TransactionType string
type PaymentMethod string

const (
	Credit     TransactionType = "credit"
	Debit      TransactionType = "debit"
	Refund     TransactionType = "refund"
	Cashback   TransactionType = "cashback"
	Withdrawal TransactionType = "withdrawal"
)

const (
	CompletedTransaction TransactionStatus = "completed"
	PendingTransaction   TransactionStatus = "pending"
	FailedTransaction    TransactionStatus = "failed"
	CancelledTransaction TransactionStatus = "cancelled"
)

const (
	PaymentMethodCard         PaymentMethod = "card"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
	PaymentMethodCash         PaymentMethod = "cash"
)

type TransactionModel struct {
	BaseModel
	UserId        string              `json:"user_id" gorm:"type:uuid;not null"`
	User          UserModel           `json:"user" gorm:"foreignKey:UserId;references:Id"`
	Type          TransactionType     `json:"type" gorm:"type:varchar(20);not null"`
	Status        TransactionStatus   `json:"status" gorm:"type:varchar(20);default:pending"`
	Amount        int64               `json:"amount" gorm:"not null"`
	Currency      string              `json:"currency" gorm:"not null;default:NGN"`
	Description   string              `json:"description"`
	Reference     string              `json:"reference" gorm:"uniqueIndex"`
	OrderId       string              `json:"order_id"`
	Order         OrderModel          `json:"order,omitempty" gorm:"foreignKey:OrderId"`
	Metadata      TransactionMetadata `json:"metadata" gorm:"type:json"`
	ProcessedAt   *time.Time          `json:"processed_at"`
	CreatedBy     string              `json:"created_by" gorm:"column:created_by;not null" validate:"required"`
	CreatedByUser UserModel           `json:"created_by_user" gorm:"foreignKey:CreatedBy;references:Id"`
	IsDeleted     bool                `json:"is_deleted" gorm:"default:false"`
}

type TransactionMetadata struct {
	PaymentMethod *string `json:"payment_method,omitempty"`
	Gateway       *string `json:"gateway,omitempty"`
	GatewayRef    *string `json:"gateway_ref,omitempty"`
	AdminUserId   *string `json:"admin_user_id,omitempty"`
	Reason        *string `json:"reason,omitempty"`
}

func (t *TransactionModel) BeforeCreate(tx *gorm.DB) error {
	if t.Id == uuid.Nil {
		t.Id = uuid.New()
	}

	now := time.Now()
	t.CreatedAt = now
	t.ModifiedAt = sql.NullTime{Time: now, Valid: true}
	t.Status = "pending"

	ref := lib.GenerateTransactionReference()
	t.Reference = ref

	return nil
}

func (t *TransactionModel) BeforeUpdate(tx *gorm.DB) error {
	t.ModifiedAt = sql.NullTime{Time: time.Now(), Valid: true}
	return nil
}
