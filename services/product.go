package services

import (
	"adire-apparel/dto"
	"adire-apparel/models"
	"errors"
	"fmt"
	"log"
	"strings"

	"gorm.io/gorm"
)

type ProductService struct {
	database *gorm.DB
}

func NewProductService(database *gorm.DB) *ProductService {
	return &ProductService{
		database: database,
	}
}

var (
	ErrProductNotFound  = errors.New("product not found")
	ErrNoImagesToRemove = errors.New("no images specified for removal")
	ErrImageNotFound    = errors.New("one or more images not found on product")
)

func (s *ProductService) CreateProduct(payload dto.CreateProductDto) (*models.ProductModel, error) {
	category, err := NewCategoryService(s.database).GetCategory(payload.CategoryId)
	if err != nil {
		return nil, err
	}

	exist := checkIfProductExists(s.database, payload.Name)
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

	log.Println("here validating user...")

	product := models.ProductModel{
		Name:        payload.Name,
		Description: payload.Description,
		Price:       payload.Price,
		Category:    *category,
		Variants:    payload.Variants,
		CategoryId:  payload.CategoryId,
		CreatedBy:   payload.CreatedBy,
	}

	if err := s.database.Create(&product).Error; err != nil {
		return nil, err
	}

	if err := s.database.Preload("CreatedByUser").First(&product, product.Id).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (s *ProductService) UpdateProduct(id string, payload dto.UpdateProductDto) (*models.ProductModel, error) {
	product := &models.ProductModel{}
	if err := s.database.Where("id = ?", id).First(product).Error; err != nil {
		return nil, err
	}

	if product.CategoryId != payload.CategoryId {
		category, err := NewCategoryService(s.database).GetCategory(payload.CategoryId)

		if err != nil {
			return nil, err
		}

		product.Category = *category
		product.CategoryId = payload.CategoryId
	}

	product.Name = payload.Name
	product.Description = payload.Description
	product.Price = payload.Price
	product.Variants = payload.Variants
	product.Status = payload.Status

	if err := s.database.Save(product).Error; err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) GetProducts(params dto.ProductPagination) (*dto.PaginatedResponse[models.ProductModel], error) {
	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Page <= 0 {
		params.Page = 1
	}

	var products []models.ProductModel
	var totalItems int64

	query := s.database.Model(&models.ProductModel{}).Where("is_deleted = ?", false)

	if params.Name != "" {
		lowerCaseName := strings.ToLower(params.Name)
		query = query.Where("name = ?", lowerCaseName)
	}

	if params.CategoryId != "" {
		query = query.Where("category_id = ?", params.CategoryId)
	}

	if params.MinPrice > 0 {
		query = query.Where("price >= ?", params.MinPrice)
	}

	if params.MaxPrice > 0 {
		query = query.Where("price <= ?", params.MaxPrice)
	}

	if err := query.Count(&totalItems).Error; err != nil {
		return &dto.PaginatedResponse[models.ProductModel]{
			Data:       []models.ProductModel{},
			Limit:      params.Limit,
			Page:       params.Page,
			TotalItems: 0,
			TotalPages: 0,
		}, err
	}

	offset := (params.Page - 1) * params.Limit

	if err := query.
		Preload("Category").
		Preload("CreatedByUser").
		Preload("Variants").
		Offset(offset).
		Limit(params.Limit).
		Order("created_at DESC").
		Find(&products).Error; err != nil {
		return &dto.PaginatedResponse[models.ProductModel]{
			Data:       []models.ProductModel{},
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

	return &dto.PaginatedResponse[models.ProductModel]{
		Data:       products,
		TotalItems: int(totalItems),
		TotalPages: int(totalPages),
		Page:       params.Page,
		Limit:      params.Limit,
	}, nil
}

func (s *ProductService) GetProduct(id string) (*models.ProductModel, error) {
	var product models.ProductModel

	if err := s.database.Where("id = ?", id).First(&product).Error; err != nil {
		return nil, ErrProductNotFound
	}

	if product.IsDeleted {
		return nil, ErrProductNotFound
	}

	return &product, nil
}

func (s *ProductService) DeleteProduct(id string) error {
	product := &models.ProductModel{}
	if err := s.database.Where("id = ?", id).First(product).Error; err != nil {
		return err
	}

	product.IsDeleted = true

	if err := s.database.Save(product).Error; err != nil {
		return err
	}

	return nil
}

func (s *ProductService) AddImagesToProduct(id string, images []string) error {
	var product models.ProductModel
	if err := s.database.Where("id = ?", id).First(&product).Error; err != nil {
		return err
	}

	product.ImageUrls = images

	if err := s.database.Save(&product).Error; err != nil {
		return err
	}

	return nil
}

func (s *ProductService) RemoveImagesFromProduct(id string, payload dto.RemoveImagesDto) error {
	if id == "" {
		return errors.New("product ID cannot be empty")
	}
	if len(payload.Images) == 0 {
		return ErrNoImagesToRemove
	}

	var product models.ProductModel
	if err := s.database.Where("id = ?", id).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProductNotFound
		}
		return fmt.Errorf("failed to fetch product: %w", err)
	}

	if len(product.ImageUrls) == 0 {
		return errors.New("product has no images to remove")
	}

	removeMap := make(map[string]bool)
	for _, image := range payload.Images {
		removeMap[image] = true
	}

	existingImageMap := make(map[string]bool)
	for _, image := range product.ImageUrls {
		existingImageMap[image] = true
	}

	for _, imageToRemove := range payload.Images {
		if !existingImageMap[imageToRemove] {
			return fmt.Errorf("image '%s' not found on product", imageToRemove)
		}
	}

	var newImages []string
	for _, existingImage := range product.ImageUrls {
		if !removeMap[existingImage] {
			newImages = append(newImages, existingImage)
		}
	}

	if len(newImages) == 0 {
		return errors.New("cannot remove all images from product")
	}

	product.ImageUrls = newImages

	if err := s.database.Save(&product).Error; err != nil {
		return fmt.Errorf("failed to update product images: %w", err)
	}

	return nil
}

func checkIfProductExists(database *gorm.DB, name string) bool {
	var category models.CategoryModel
	if err := database.Where("name = ?", name).First(&category).Error; err != nil {
		return false
	}

	return true
}
