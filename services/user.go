package services

import (
	"adire-apparel/dto"
	"adire-apparel/models"
	"errors"
	"strings"

	"gorm.io/gorm"
)

type UserService struct {
	database *gorm.DB
}

func NewUserService(database *gorm.DB) *UserService {
	return &UserService{
		database: database,
	}
}

var (
	ErrEmailExists       = errors.New("a user with this email already exists")
	ErrPhoneExists       = errors.New("a user with this phone number already exists")
	ErrUserNotAuthorized = errors.New("user is not authorized to perform this action")
)

func (s *UserService) GetAllUsers(params dto.UserPagination) (*dto.PaginatedResponse[models.UserModel], error) {
	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Page <= 0 {
		params.Page = 1
	}

	var users []models.UserModel
	var totalItems int64

	query := s.database.Model(&models.UserModel{})

	if params.Name != nil && strings.TrimSpace(*params.Name) != "" {
		namePattern := "%" + strings.ToLower(strings.TrimSpace(*params.Name)) + "%"
		query = query.Where("LOWER(name) LIKE ?", namePattern)
	}

	if params.Email != nil && strings.TrimSpace(*params.Email) != "" {
		emailPattern := "%" + strings.ToLower(strings.TrimSpace(*params.Email)) + "%"
		query = query.Where("LOWER(email) LIKE ?", emailPattern)
	}

	if params.Status != nil {
		query = query.Where("status = ?", *params.Status)
	}

	if params.IsDeleted != nil {
		query = query.Where("is_deleted = ?", *params.IsDeleted)
	} else {
		query = query.Where("is_deleted = ?", false)
	}

	if err := query.Count(&totalItems).Error; err != nil {
		return &dto.PaginatedResponse[models.UserModel]{
			Data:       []models.UserModel{},
			Limit:      params.Limit,
			Page:       params.Page,
			TotalItems: 0,
			TotalPages: 0,
		}, err
	}

	offset := (params.Page - 1) * params.Limit

	if err := query.Offset(offset).
		Limit(params.Limit).
		Find(&users).Error; err != nil {
		return nil, err
	}

	totalPages := int64(0)
	if params.Limit > 0 {
		totalPages = (totalItems + int64(params.Limit) - 1) / int64(params.Limit)
	}

	return &dto.PaginatedResponse[models.UserModel]{
		Data:       users,
		TotalItems: int(totalItems),
		TotalPages: int(totalPages),
		Page:       params.Page,
		Limit:      params.Limit,
	}, nil
}

func (s *UserService) GetUserById(id string) (*models.UserModel, error) {
	var user models.UserModel
	if err := s.database.Where("id = ? AND is_deleted = ?", id, false).First(&user).Error; err != nil {
		return nil, ErrUserNotFound
	}

	return &user, nil
}

func (s *UserService) GetUserByEmail(email string) (*models.UserModel, error) {
	var user models.UserModel
	if err := s.database.Where("email = ? AND is_deleted = ?", email, false).First(&user).Error; err != nil {
		return nil, ErrUserNotFound
	}

	return &user, nil
}

func (s *UserService) UpdateUser(id string, update *dto.EditUserDto) (*models.UserModel, error) {
	user := &models.UserModel{}
	if err := s.database.Where("id = ?", id).First(user).Error; err != nil {
		return nil, err
	}

	if emailExists(s.database, *update.Email) && user.Email != *update.Email {
		return nil, ErrEmailExists
	}

	if phoneExists(s.database, *update.Phone) && user.Phone != *update.Phone {
		return nil, ErrPhoneExists
	}

	if update.Name != nil && *update.Name != "" && *update.Name != user.Name {
		user.Name = *update.Name
	}

	if update.Phone != nil && *update.Phone != "" && *update.Phone != user.Phone {
		user.Phone = *update.Phone
	}
	if update.Email != nil && *update.Email != "" && *update.Email != user.Email {
		user.Email = *update.Email
	}

	err := s.database.Save(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) UpdateUserPreferences(id string, payload dto.EditUserPreferencesDto) (*models.UserPreferences, error) {
	user, err := NewAuthService(s.database).FindUserById(id)
	if err != nil {
		return nil, err
	}

	if payload.DisplayCurrency != "" && payload.DisplayCurrency != user.Preferences.DisplayCurrency {
		user.Preferences.DisplayCurrency = payload.DisplayCurrency
	}
	if payload.EmailNotifications != user.Preferences.EmailNotifications {
		user.Preferences.EmailNotifications = payload.EmailNotifications
	}
	if payload.MarketingEmails != user.Preferences.MarketingEmails {
		user.Preferences.MarketingEmails = payload.MarketingEmails
	}
	if payload.NewsletterSubscription != user.Preferences.NewsletterSubscription {
		user.Preferences.NewsletterSubscription = payload.NewsletterSubscription
	}
	if payload.PushNotifications != user.Preferences.PushNotifications {
		user.Preferences.PushNotifications = payload.PushNotifications
	}
	if payload.SMSNotifications != user.Preferences.SMSNotifications {
		user.Preferences.SMSNotifications = payload.SMSNotifications
	}
	if payload.Theme != "" && payload.Theme != user.Preferences.Theme {
		user.Preferences.Theme = payload.Theme
	}

	if err = s.database.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user.Preferences, nil
}

func (s *UserService) DeleteUser(id string) error {
	user, err := NewAuthService(s.database).FindUserById(id)
	if err != nil {
		return ErrUserNotFound
	}

	user.IsDeleted = true

	if err := s.database.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

func (s *UserService) UploadImage(id, image string) error {
	user, err := NewAuthService(s.database).FindUserById(id)
	if err != nil {
		return ErrUserNotFound
	}

	user.ImageUrl = image

	if err := s.database.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

func emailExists(db *gorm.DB, email string) bool {
	var user models.UserModel
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return false
	}

	return true
}

func phoneExists(db *gorm.DB, phone string) bool {
	var user models.UserModel
	if err := db.Where("phone = ?", phone).First(&user).Error; err != nil {
		return false
	}

	return true
}
