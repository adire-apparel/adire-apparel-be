package services

import (
	"adire-apparel/dto"
	"adire-apparel/models"
	"errors"

	"gorm.io/gorm"
)

type TransactionService struct {
	database *gorm.DB
}

func NewTransactionService(database *gorm.DB) *TransactionService {
	return &TransactionService{
		database: database,
	}
}

var (
	ErrTransactionNotFound = errors.New("transaction not found")
)

func (s *TransactionService) CreateTransaction(payload *dto.CreateTransactionDto) (*models.TransactionModel, error) {
	_, err := NewUserService(s.database).GetUserById(payload.UserId)
	if err != nil {
		return nil, err
	}

	admin, err := NewUserService(s.database).GetUserById(string(payload.CreatedBy))
	if err != nil {
		return nil, err
	}

	if admin.Status != "active" {
		return nil, gorm.ErrInvalidTransaction
	}

	_, err = NewOrderService(s.database).GetOrder(payload.OrderId)
	if err != nil {
		return nil, ErrOrderNotFound
	}

	transaction := &models.TransactionModel{
		UserId:      payload.UserId,
		CreatedBy:   payload.CreatedBy,
		OrderId:     payload.OrderId,
		Type:        payload.Type,
		Amount:      payload.Amount,
		Currency:    payload.Currency,
		Description: payload.Description,
	}

	if err := s.database.Create(&transaction).Error; err != nil {
		return nil, err
	}

	if err := s.database.Preload("User").Preload("CreatedByUser").Preload("Order").First(&transaction).Error; err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *TransactionService) UpdateTransaction(id string, payload *dto.UpdateTransactionDto) (*models.TransactionModel, error) {
	transaction := &models.TransactionModel{}
	if err := s.database.Where("id = ?", id).First(&transaction).Error; err != nil {
		return nil, ErrTransactionNotFound
	}

	transaction.Status = payload.Status
	transaction.Metadata = payload.Metadata

	if err := s.database.Save(transaction).Error; err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *TransactionService) GetTransactions(params *dto.Pagination) (*dto.PaginatedResponse[models.TransactionModel], error) {
	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Page <= 0 {
		params.Page = 1
	}

	var transactions []models.TransactionModel
	var totalItems int64

	query := s.database.Model(&models.TransactionModel{}).Where("is_deleted = ?", false)

	if err := query.Count(&totalItems).Error; err != nil {
		return &dto.PaginatedResponse[models.TransactionModel]{
			Data:       []models.TransactionModel{},
			Limit:      params.Limit,
			Page:       params.Page,
			TotalItems: 0,
			TotalPages: 0,
		}, err
	}

	offset := (params.Page - 1) * params.Limit

	if err := query.
		Preload("User").
		Preload("CreatedByUser").
		Preload("Order").
		Offset(offset).
		Limit(params.Limit).
		Order("created_at DESC").
		Find(&transactions).Error; err != nil {
		return &dto.PaginatedResponse[models.TransactionModel]{
			Data:       []models.TransactionModel{},
			Limit:      params.Limit,
			Page:       params.Page,
			TotalItems: 0,
			TotalPages: 0,
		}, err
	}

	totalPages := int64(0)
	if params.Limit > 0 {
		totalPages = (totalItems + int64(params.Limit) - 1) / int64(params.Limit)
	}

	return &dto.PaginatedResponse[models.TransactionModel]{
		Data:       transactions,
		TotalItems: int(totalItems),
		TotalPages: int(totalPages),
		Page:       params.Page,
		Limit:      params.Limit,
	}, nil
}

func (s *TransactionService) GetTransactionsByUser(id string, params *dto.Pagination) (*dto.PaginatedResponse[models.TransactionModel], error) {
	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Page <= 0 {
		params.Page = 1
	}

	var transactions []models.TransactionModel
	var totalItems int64

	query := s.database.Model(&models.TransactionModel{}).Where("user_id = ? AND is_deleted = ?", id, false)

	if err := query.Count(&totalItems).Error; err != nil {
		return &dto.PaginatedResponse[models.TransactionModel]{
			Data:       []models.TransactionModel{},
			Limit:      params.Limit,
			Page:       params.Page,
			TotalItems: 0,
			TotalPages: 0,
		}, err
	}

	offset := (params.Page - 1) * params.Limit

	if err := query.
		Preload("User").
		Preload("CreatedByUser").
		Preload("Order").
		Offset(offset).
		Limit(params.Limit).
		Order("created_at DESC").
		Find(&transactions).Error; err != nil {
		return &dto.PaginatedResponse[models.TransactionModel]{
			Data:       []models.TransactionModel{},
			Limit:      params.Limit,
			Page:       params.Page,
			TotalItems: 0,
			TotalPages: 0,
		}, err
	}

	totalPages := int64(0)
	if params.Limit > 0 {
		totalPages = (totalItems + int64(params.Limit) - 1) / int64(params.Limit)
	}

	return &dto.PaginatedResponse[models.TransactionModel]{
		Data:       transactions,
		TotalItems: int(totalItems),
		TotalPages: int(totalPages),
		Page:       params.Page,
		Limit:      params.Limit,
	}, nil
}

func (s *TransactionService) GetTransactionsByAdmin(id string, params *dto.Pagination) (*dto.PaginatedResponse[models.TransactionModel], error) {
	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Page <= 0 {
		params.Page = 1
	}

	var transactions []models.TransactionModel
	var totalItems int64

	query := s.database.Model(&models.TransactionModel{}).Where("created_by = ? AND is_deleted = ?", id, false)

	if err := query.Count(&totalItems).Error; err != nil {
		return &dto.PaginatedResponse[models.TransactionModel]{
			Data:       []models.TransactionModel{},
			Limit:      params.Limit,
			Page:       params.Page,
			TotalItems: 0,
			TotalPages: 0,
		}, err
	}

	offset := (params.Page - 1) * params.Limit

	if err := query.
		Preload("User").
		Preload("CreatedByUser").
		Preload("Order").
		Offset(offset).
		Limit(params.Limit).
		Order("created_at DESC").
		Find(&transactions).Error; err != nil {
		return &dto.PaginatedResponse[models.TransactionModel]{
			Data:       []models.TransactionModel{},
			Limit:      params.Limit,
			Page:       params.Page,
			TotalItems: 0,
			TotalPages: 0,
		}, err
	}

	totalPages := int64(0)
	if params.Limit > 0 {
		totalPages = (totalItems + int64(params.Limit) - 1) / int64(params.Limit)
	}

	return &dto.PaginatedResponse[models.TransactionModel]{
		Data:       transactions,
		TotalItems: int(totalItems),
		TotalPages: int(totalPages),
		Page:       params.Page,
		Limit:      params.Limit,
	}, nil
}

func (s *TransactionService) GetTransaction(id string) (*models.TransactionModel, error) {
	var transaction models.TransactionModel

	if err := s.database.Preload("User").Preload("CreatedByUser").Preload("Order").Where("id = ?", id).First(transaction).Error; err != nil {
		return nil, ErrProductNotFound

	}

	if transaction.IsDeleted {
		return nil, ErrTransactionNotFound
	}

	return &transaction, nil
}

func (s *TransactionService) DeleteTransaction(id string) error {
	transaction := &models.TransactionModel{}
	if err := s.database.Where("id = ?", id).First(transaction).Error; err != nil {
		return err
	}

	transaction.IsDeleted = true

	if err := s.database.Save(transaction).Error; err != nil {
		return err
	}

	return nil
}
