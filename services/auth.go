package services

import (
	"adire-apparel/config"
	"adire-apparel/dto"
	"adire-apparel/lib"
	"adire-apparel/models"
	"errors"
	"fmt"
	"time"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/google"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
	ErrUserExists      = errors.New("user already exists")
	ErrEmailSendFailed = errors.New("failed to send email")
	ErrInvalidOtp      = errors.New("invalid otp")
	ErrInvalidProvider = errors.New("invalid provider")
	ErrPasswordIsSame  = errors.New("new password is same as old password")
)

type AuthService struct {
	database *gorm.DB
}

func NewAuthService(database *gorm.DB) *AuthService {
	return &AuthService{
		database: database,
	}
}

type SigninResponse struct {
	User  models.UserModel `json:"user"`
	Token string           `json:"token"`
}

func (s *AuthService) Signup(payload dto.CreateUserDto) error {
	exists, err := s.userExistsByEmail(payload.Email)
	if err != nil {
		return err
	}

	if exists {
		return ErrUserExists
	}

	if !lib.ValidatePassword(payload.Password) {
		return ErrInvalidPassword
	}

	hashedPassword, err := lib.HashPassword(payload.Password)
	if err != nil {
		return err
	}

	otp := lib.GenerateOtp()

	user := models.UserModel{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: hashedPassword,
		Role:     payload.Role,
		Otp:      otp,
	}

	if err := s.database.Create(&user).Error; err != nil {
		return err
	}

	go func() {
		lib.SendEmail(lib.EmailDto{
			To:       []string{payload.Email},
			Subject:  "Welcome to Adire Apparel",
			Template: "welcome",
			Data: map[string]interface{}{
				"name": payload.Name,
				"otp":  otp,
			},
		})
	}()

	return nil
}

func (s *AuthService) Signin(payload dto.SigninDto) (*SigninResponse, error) {
	user, err := s.FindUserByEmail(payload.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if user.LockedUntil != nil && user.LockedUntil.After(time.Now()) {
		return nil, fmt.Errorf("account locked until %s", user.LockedUntil.Format("2006-01-02 15:04:05"))
	}

	err = lib.ComparePassword(payload.Password, user.Password)

	if err != nil {
		user.LoginAttempts += 1
		if user.LoginAttempts >= 3 {
			lockedUntil := time.Now().Add(time.Minute * 15)
			user.LockedUntil = &lockedUntil
		}
		err = s.database.Save(&user).Error
		if err != nil {
			return nil, err
		}

		return nil, ErrInvalidPassword
	}

	user.LoginAttempts = 0
	user.LockedUntil = nil
	user.LastLoggedIn = time.Now()

	token, err := lib.GenerateToken(user.Id)
	if err != nil {
		return nil, err
	}

	err = s.database.Save(&user).Error
	if err != nil {
		return nil, err
	}

	user.Password = ""
	user.Otp = ""

	return &SigninResponse{
		User:  user,
		Token: token,
	}, nil
}

func (s *AuthService) SigninWithOauth(payload *goth.User) (*models.UserModel, error) {
	if payload.Provider != "google" && payload.Provider != "facebook" {
		return nil, ErrInvalidProvider
	}

	return nil, nil
}

func (s *AuthService) Refresh(id string) (*string, error) {
	return nil, nil
}

func (s *AuthService) VerifyEmail(payload dto.VerifyEmailDto) error {
	user, err := s.FindUserByEmail(payload.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	if user.Otp != payload.Otp {
		return ErrInvalidOtp
	}

	user.Otp = ""
	user.EmailVerified = true

	err = s.database.Save(&user).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) VerifyPhone(payload dto.VerifyPhoneDto) error {
	user, err := s.FindUserByPhone(payload.Phone)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	if user.Otp != payload.Otp {
		return ErrInvalidOtp
	}

	user.Otp = ""
	user.PhoneVerified = true

	err = s.database.Save(&user).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) ChangePassword(id string, payload dto.ChangePasswordDto) error {
	user, err := s.FindUserById(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	if lib.ComparePassword(payload.NewPassword, user.Password) == nil {
		return ErrPasswordIsSame
	}

	if !lib.ValidatePassword(payload.NewPassword) {
		return ErrInvalidPassword
	}

	hashedPassword, err := lib.HashPassword(payload.NewPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword

	if err := s.database.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

func (s *AuthService) ForgotPassword(email string) error {
	user, err := s.FindUserByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	otp := lib.GenerateOtp()
	user.Otp = otp

	err = s.database.Save(&user).Error
	if err != nil {
		return err
	}

	go func() {
		lib.SendEmail(lib.EmailDto{
			To:       []string{email},
			Subject:  "Forgot Password",
			Template: "forgot-password",
			Data: map[string]interface{}{
				"otp": otp,
			},
		})
	}()

	return nil
}

func (s *AuthService) ResetPassword(payload dto.ResetPasswordDto) error {
	user, err := s.FindUserByEmail(payload.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	if user.Otp != payload.Otp {
		return ErrInvalidOtp
	}

	hashedPassword, err := lib.HashPassword(payload.Password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	user.Otp = ""

	err = s.database.Save(&user).Error
	if err != nil {
		return err
	}

	go func() {
		lib.SendEmail(lib.EmailDto{
			To:       []string{user.Email},
			Subject:  "Password Reset",
			Template: "password-reset",
			Data: map[string]interface{}{
				"name": user.Name,
			},
		})
	}()

	return nil
}

func (s *AuthService) FindUserByEmail(email string) (models.UserModel, error) {
	var user models.UserModel

	if err := s.database.Where("email = ?", email).First(&user).Error; err != nil {
		return models.UserModel{}, err
	}

	return user, nil
}

func (s *AuthService) FindUserByPhone(phone string) (models.UserModel, error) {
	var user models.UserModel

	if err := s.database.Where("phone = ?", phone).First(&user).Error; err != nil {
		return models.UserModel{}, err
	}

	return user, nil
}

func (s *AuthService) FindUserById(id string) (models.UserModel, error) {
	var user models.UserModel

	if err := s.database.Where("id = ?", id).First(&user).Error; err != nil {
		return models.UserModel{}, err
	}

	return user, nil
}

func InitializeProvider() {
	appConfig := config.AppConfig

	store := sessions.NewCookieStore([]byte(appConfig.GoogleAuthId))
	store.MaxAge(int(time.Hour) * 30)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = !appConfig.IsDevMode

	gothic.Store = store
	goth.UseProviders(
		google.New(
			appConfig.GoogleAuthId,
			appConfig.GoogleAuthSecret,
			fmt.Sprintf("%s/auth/google/callback", appConfig.ApiUrl),
			"email", "profile"),
		facebook.New(
			appConfig.FacebookAuthId,
			appConfig.FacebookAuthSecret,
			fmt.Sprintf("%s/auth/facebook/callback", appConfig.ApiUrl),
			"email", "public_profile"),
	)
}

func (s *AuthService) userExistsByEmail(email string) (bool, error) {
	var count int64
	err := s.database.Model(&models.UserModel{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}
