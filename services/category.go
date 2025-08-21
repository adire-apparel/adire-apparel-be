package services

import (
	"adire-apparel/dto"
	"adire-apparel/models"

	"gorm.io/gorm"
)

type CategoryService struct {
	database *gorm.DB
}

func NewCategoryService(database *gorm.DB) *CategoryService {
	return &CategoryService{
		database: database,
	}
}

func (s *CategoryService) CreateCategory(payload dto.CreateCategoryDto) (*models.CategoryModel, error) {
	exist := checkIfCategoryExists(s.database, payload.Name)
	if exist {
		return nil, gorm.ErrDuplicatedKey
	}

	user, err := NewUserService(s.database).GetUserById(payload.CreatedBy)
	if err != nil {
		return nil, err
	}

	if user.Role != "admin" {
		return nil, gorm.ErrInvalidTransaction
	}

	category := models.CategoryModel{
		Name:        payload.Name,
		Description: payload.Description,
		CreatedBy:   user.Id.String(),
	}

	if err := s.database.Create(&category).Error; err != nil {
		return nil, err
	}

	if err := s.database.Preload("CreatedByUser", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name, image_url, role, status")
	}).First(&category, category.Id).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

func (s *CategoryService) UpdateCategory(id string, payload dto.CreateCategoryDto) (*models.CategoryModel, error) {
	category := &models.CategoryModel{}
	if err := s.database.Where("id = ?", id).First(category).Error; err != nil {
		return nil, err
	}

	category.Name = payload.Name
	category.Description = payload.Description

	err := s.database.Save(category).Error
	if err != nil {
		return nil, err
	}

	if err := s.database.Preload("CreatedByUser", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name, image_url, role, status")
	}).First(&category, category.Id).Error; err != nil {
		return nil, err
	}

	return category, nil
}

func (s *CategoryService) GetCategories(params dto.Pagination) (*dto.PaginatedResponse[models.CategoryModel], error) {
	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Page <= 0 {
		params.Page = 1
	}

	var categories []models.CategoryModel
	var totalItems int64

	if err := s.database.Model(&models.CategoryModel{}).Count(&totalItems).Error; err != nil {
		return &dto.PaginatedResponse[models.CategoryModel]{
			Data:       []models.CategoryModel{},
			Limit:      params.Limit,
			Page:       params.Page,
			TotalItems: 0,
			TotalPages: 0,
		}, err
	}

	offset := (params.Page - 1) * params.Limit

	if err := s.database.Model(&models.CategoryModel{}).Preload("CreatedByUser", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name, image_url, role, status")
	}).Offset(offset).Limit(params.Limit).Find(&categories).Error; err != nil {
		return &dto.PaginatedResponse[models.CategoryModel]{
			Data:       []models.CategoryModel{},
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

	return &dto.PaginatedResponse[models.CategoryModel]{
		Data:       categories,
		TotalItems: int(totalItems),
		TotalPages: int(totalPages),
		Page:       params.Page,
		Limit:      params.Limit,
	}, nil
}

func (s *CategoryService) GetCategory(id string) (*models.CategoryModel, error) {
	var category models.CategoryModel

	if err := s.database.Preload("CreatedByUser", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name, image_url, role, status")
	}).Where("id = ?", id).First(&category).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

func (s *CategoryService) DeleteCategory(id string) error {
	if err := s.database.Where("id = ?", id).Delete(&models.CategoryModel{}).Error; err != nil {
		return err
	}

	return nil
}

func checkIfCategoryExists(database *gorm.DB, name string) bool {
	var category models.CategoryModel
	if err := database.Where("name = ?", name).First(&category).Error; err != nil {
		return false
	}

	return true
}
