package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

type Role string
type UserStatus string

const (
	Customer Role = "customer"
	Staff    Role = "staff"
	Admin    Role = "admin"
)

const (
	ActiveUser    UserStatus = "active"
	InactiveUser  UserStatus = "inactive"
	SuspendedUser UserStatus = "suspended"
	PendingUser   UserStatus = "pending"
)

type UserModel struct {
	BaseModel
	Name              string            `json:"name" gorm:"not null" validate:"required"`
	Email             string            `json:"email" gorm:"not null" validate:"required unique"`
	Password          string            `json:"-" gorm:"not null" validate:"required"`
	Phone             string            `json:"phone,omitempty" gorm:""`
	Role              Role              `json:"role" gorm:"not null" validate:"required"`
	LastLoggedIn      time.Time         `json:"last_logged_in" gorm:"type:TIMESTAMP with time zone;null"`
	ImageUrl          string            `json:"image_url,omitempty" gorm:""`
	Status            UserStatus        `json:"status" gorm:"type:varchar(20);default:active" validate:"required,oneof=active inactive suspended pending"`
	EmailVerified     bool              `json:"email_verified" gorm:"default:false"`
	PhoneVerified     bool              `json:"phone_verified" gorm:"default:false"`
	TwoFactorEnabled  bool              `json:"two_factor_enabled" gorm:"default:false"`
	LoginAttempts     int               `json:"-" gorm:"default:0"`
	LockedUntil       *time.Time        `json:"-"`
	PreferredLanguage string            `json:"preferred_language" gorm:"default:en"`
	PreferredCurrency string            `json:"preferred_currency" gorm:"default:NGN"`
	Timezone          string            `json:"timezone" gorm:"default:UTC"`
	Orders            []OrderModel      `json:"orders,omitempty" gorm:"foreignKey:UserId"`
	Wallet            *WalletModel      `json:"wallet" gorm:"foreignKey:UserId"`
	Notifications     NotificationArray `json:"notifications" gorm:"type:json"`
	Preferences       UserPreferences   `json:"preferences" gorm:"type:json"`
	Otp               string            `json:"-"`
	Token             string            `json:"-"`
	IsDeleted         bool              `json:"-" gorm:"default:false"`
}

type WalletModel struct {
	BaseModel
	UserId            string                   `json:"user_id" gorm:"not null;uniqueIndex"`
	User              UserModel                `json:"user" gorm:"foreignKey:UserId"`
	Balance           int64                    `json:"balance" gorm:"default:0"`
	Currency          string                   `json:"currency" gorm:"default:NGN;not null"`
	Status            UserStatus               `json:"status" gorm:"type:varchar(20);default:active"`
	DailyLimit        *int64                   `json:"daily_limit"`
	MonthlyLimit      *int64                   `json:"monthly_limit"`
	LastTransactionAt *time.Time               `json:"last_transaction_at"`
	Transactions      []WalletTransactionModel `json:"transactions,omitempty" gorm:"foreignKey:WalletId"`
}

type WalletTransactionModel struct {
	BaseModel
	WalletId      string              `json:"wallet_id" gorm:"not null"`
	Wallet        WalletModel         `json:"wallet" gorm:"foreignKey:WalletId"`
	Type          TransactionType     `json:"type" gorm:"type:varchar(20);not null"`
	Status        TransactionStatus   `json:"status" gorm:"type:varchar(20);default:pending"`
	Amount        int64               `json:"amount" gorm:"not null"`
	Currency      string              `json:"currency" gorm:"not null;default:NGN"`
	Description   string              `json:"description"`
	Reference     *string             `json:"reference" gorm:"uniqueIndex"`
	OrderId       *string             `json:"order_id"`
	Order         *OrderModel         `json:"order,omitempty" gorm:"foreignKey:OrderId"`
	BalanceBefore int64               `json:"balance_before" gorm:"not null"`
	BalanceAfter  int64               `json:"balance_after" gorm:"not null"`
	Metadata      TransactionMetadata `json:"metadata" gorm:"type:json"`
	ProcessedAt   *time.Time          `json:"processed_at"`
	ExpiresAt     *time.Time          `json:"expires_at"`
}

type UserPreferences struct {
	EmailNotifications     bool   `json:"email_notifications"`
	PushNotifications      bool   `json:"push_notifications"`
	SMSNotifications       bool   `json:"sms_notifications"`
	MarketingEmails        bool   `json:"marketing_emails"`
	NewsletterSubscription bool   `json:"newsletter_subscription"`
	Theme                  string `json:"theme"`
	DisplayCurrency        string `json:"display_currency"`
}

type NotificationModel struct {
	Id        string    `json:"id"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Read      bool      `json:"read"`
	CreatedAt time.Time `json:"created_at"`
}

type NotificationArray []NotificationModel

func (na *NotificationArray) Scan(value interface{}) error {
	if value == nil {
		*na = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, na)
}

func (na NotificationArray) Value() (driver.Value, error) {
	if len(na) == 0 {
		return nil, nil
	}
	return json.Marshal(na)
}

func (up *UserPreferences) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, up)
}

func (up UserPreferences) Value() (driver.Value, error) {
	return json.Marshal(up)
}

func (tm *TransactionMetadata) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, tm)
}

func (tm TransactionMetadata) Value() (driver.Value, error) {
	return json.Marshal(tm)
}

func (u *UserModel) BeforeCreate(tx *gorm.DB) error {
	if u.Role == "" {
		u.Role = "customer"
	}
	if u.Status == "" {
		u.Status = ActiveUser
	}
	return nil
}

func (u *UserModel) IsLocked() bool {
	return u.LockedUntil != nil && u.LockedUntil.After(time.Now())
}

func (u *UserModel) CanLogin() bool {
	return u.Status == ActiveUser && !u.IsLocked()
}

func (w *WalletModel) BeforeCreate(tx *gorm.DB) error {
	if w.Currency == "" {
		w.Currency = "NGN"
	}
	if w.Status == "" {
		w.Status = ActiveUser
	}
	return nil
}

func (w *WalletModel) CanTransact() bool {
	return w.Status == ActiveUser
}

func (w *WalletModel) HasSufficientBalance(amount int64) bool {
	return w.Balance >= amount
}

func (w *WalletModel) WithinDailyLimit(amount int64) bool {
	if w.DailyLimit == nil {
		return true
	}
	return amount <= *w.DailyLimit
}

func (w *WalletModel) WithinMonthlyLimit(amount int64) bool {
	if w.MonthlyLimit == nil {
		return true
	}
	return amount <= *w.MonthlyLimit
}
