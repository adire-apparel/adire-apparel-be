package dto

import "adire-apparel/models"

type UserPagination struct {
	Pagination
	Name      *string            `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Email     *string            `json:"email,omitempty" validate:"omitempty,email,max=255"`
	Status    *models.UserStatus `json:"status,omitempty" validate:"omitempty,oneof=active inactive suspended"`
	IsDeleted *bool              `json:"is_deleted,omitempty"`
}

type CreateUserDto struct {
	Email    string      `json:"email"`
	Name     string      `json:"name"`
	Password string      `json:"password"`
	Role     models.Role `json:"role"`
}

type EditUserDto struct {
	Email *string `json:"email"`
	Name  *string `json:"name"`
	Phone *string `json:"phone"`
}

type SigninDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SigninWithOauthDto struct {
	Provider string `json:"provider" gorm:"not null" validate:"required,oneof=google facebook"`
	Token    string `json:"token"`
}

type VerifyEmailDto struct {
	Email string `json:"email"`
	Otp   string `json:"otp"`
}

type VerifyPhoneDto struct {
	Phone string `json:"phone"`
	Otp   string `json:"otp"`
}

type ResetPasswordDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Otp      string `json:"otp"`
}

type ChangePasswordDto struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type EditUserPreferencesDto struct {
	EmailNotifications     bool   `json:"email_notifications"`
	PushNotifications      bool   `json:"push_notifications"`
	SMSNotifications       bool   `json:"sms_notifications"`
	MarketingEmails        bool   `json:"marketing_emails"`
	NewsletterSubscription bool   `json:"newsletter_subscription"`
	Theme                  string `json:"theme"`
	DisplayCurrency        string `json:"display_currency"`
}

type GoogleModel struct {
	Email     string `json:"email"`
	Name      string `json:"name"`
	FirstName string `json:"given_name"`
	LastName  string `json:"family_name"`
	UserID    string `json:"sub"`
}

type FacebookModel struct {
	Email     string `json:"email"`
	Name      string `json:"name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserID    string `json:"id"`
}
