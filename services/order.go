package services

import (
	"adire-apparel/dto"
	"adire-apparel/models"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderService struct {
	database *gorm.DB
}

func NewOrderService(database *gorm.DB) *OrderService {
	return &OrderService{
		database: database,
	}
}

var (
	ErrOrderNotFound    = errors.New("order not found")
	ErrInvalidOrderData = errors.New("invalid order data")
)

func (s *OrderService) CreateOrder(payload dto.CreateOrderDto) (*models.OrderModel, error) {
	payload.SetDefaults()

	if err := s.validateUser(payload.UserId); err != nil {
		return nil, err
	}

	orderItems, err := s.buildOrderItems(payload.Items)
	if err != nil {
		return nil, err
	}

	order := s.buildOrderModel(payload, orderItems)

	tx := s.database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	if err := s.loadOrderRelations(tx, &order); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &order, nil
}

func (s *OrderService) UpdateOrder(id string, payload dto.UpdateOrderDto) (*models.OrderModel, error) {
	order, err := s.GetOrder(id)
	if err != nil {
		return nil, err
	}

	tx := s.database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	updates := make(map[string]interface{})

	if payload.Status != nil {
		updates["status"] = *payload.Status
	}
	if payload.Currency != nil {
		updates["currency"] = *payload.Currency
	}
	if payload.Notes != nil {
		updates["notes"] = *payload.Notes
	}
	if payload.ShippingAddress != nil {
		updates["shipping_address"] = models.Address{
			Street:  payload.ShippingAddress.Street,
			City:    payload.ShippingAddress.City,
			State:   payload.ShippingAddress.State,
			ZipCode: payload.ShippingAddress.ZipCode,
			Country: payload.ShippingAddress.Country,
		}
	}
	if payload.BillingAddress != nil {
		updates["billing_address"] = models.Address{
			Street:  payload.BillingAddress.Street,
			City:    payload.BillingAddress.City,
			State:   payload.BillingAddress.State,
			ZipCode: payload.BillingAddress.ZipCode,
			Country: payload.BillingAddress.Country,
		}
	}

	if len(updates) > 0 {
		if err := tx.Model(&order).Updates(updates).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to update order: %w", err)
		}
	}

	if len(payload.Items) > 0 {
		if err := s.updateOrderItems(tx, order.Id.String(), payload.Items); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := s.loadOrderRelations(tx, order); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return order, nil
}

func (s *OrderService) UpdateOrderStatus(id string, payload dto.UpdateOrderStatusDto) (*models.OrderModel, error) {
	order, err := s.GetOrder(id)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{
		"status": payload.Status,
	}

	if payload.Notes != nil {
		updates["notes"] = *payload.Notes
	}

	if err := s.database.Model(&order).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update order status: %w", err)
	}

	if err := s.loadOrderRelations(s.database, order); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) CancelOrder(id string) error {
	order, err := s.GetOrder(id)
	if err != nil {
		return err
	}

	if order.Status == "delivered" {
		return errors.New("cannot cancel delivered order")
	}

	updates := map[string]interface{}{
		"status":       "cancelled",
		"cancelled_at": "NOW()",
	}

	if err := s.database.Model(&order).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	return nil
}

func (s *OrderService) GetOrders(params dto.Pagination) (*dto.PaginatedResponse[models.OrderModel], error) {
	params = s.normalizePagination(params)

	var orders []models.OrderModel
	var totalItems int64

	query := s.database.Model(&models.OrderModel{})

	if err := query.Count(&totalItems).Error; err != nil {
		return s.emptyPaginatedResponse(params), fmt.Errorf("failed to count orders: %w", err)
	}

	offset := (params.Page - 1) * params.Limit

	if err := query.
		Preload("Items").
		Preload("User").
		Offset(offset).
		Limit(params.Limit).
		Order("created_at DESC").
		Find(&orders).Error; err != nil {
		return s.emptyPaginatedResponse(params), fmt.Errorf("failed to fetch orders: %w", err)
	}

	totalPages := s.calculateTotalPages(totalItems, params.Limit)

	return &dto.PaginatedResponse[models.OrderModel]{
		Data:       orders,
		TotalItems: int(totalItems),
		TotalPages: totalPages,
		Page:       params.Page,
		Limit:      params.Limit,
	}, nil
}

func (s *OrderService) GetOrdersByUser(userId string, params dto.Pagination) (*dto.PaginatedResponse[models.OrderModel], error) {
	if _, err := uuid.Parse(userId); err != nil {
		return s.emptyPaginatedResponse(params), errors.New("invalid user ID format")
	}

	params = s.normalizePagination(params)

	var orders []models.OrderModel
	var totalItems int64

	query := s.database.Model(&models.OrderModel{}).Where("user_id = ?", userId)

	if err := query.Count(&totalItems).Error; err != nil {
		return s.emptyPaginatedResponse(params), fmt.Errorf("failed to count user orders: %w", err)
	}

	offset := (params.Page - 1) * params.Limit

	if err := query.
		Preload("Items").
		Offset(offset).
		Limit(params.Limit).
		Order("created_at DESC").
		Find(&orders).Error; err != nil {
		return s.emptyPaginatedResponse(params), fmt.Errorf("failed to fetch user orders: %w", err)
	}

	totalPages := s.calculateTotalPages(totalItems, params.Limit)

	return &dto.PaginatedResponse[models.OrderModel]{
		Data:       orders,
		TotalItems: int(totalItems),
		TotalPages: totalPages,
		Page:       params.Page,
		Limit:      params.Limit,
	}, nil
}

func (s *OrderService) GetOrder(id string) (*models.OrderModel, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, errors.New("invalid order ID format")
	}

	var order models.OrderModel

	if err := s.database.
		Preload("Items").
		Preload("User").
		Where("id = ?", id).
		First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to fetch order: %w", err)
	}

	return &order, nil
}

func (s *OrderService) DeleteOrder(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return errors.New("invalid order ID format")
	}

	result := s.database.Where("id = ?", id).Delete(&models.OrderModel{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete order: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrOrderNotFound
	}

	return nil
}

func (s *OrderService) validateUser(userId string) error {
	if _, err := uuid.Parse(userId); err != nil {
		return errors.New("invalid user ID format")
	}

	userService := NewUserService(s.database)
	if _, err := userService.GetUserById(userId); err != nil {
		return ErrUserNotFound
	}

	return nil
}

func (s *OrderService) buildOrderItems(items []dto.CreateOrderItemDto) ([]models.OrderItemModel, error) {
	var orderItems []models.OrderItemModel
	productService := NewProductService(s.database)

	for _, item := range items {
		if _, err := uuid.Parse(item.ProductId); err != nil {
			return nil, fmt.Errorf("invalid product ID format: %s", item.ProductId)
		}

		if _, err := productService.GetProduct(item.ProductId); err != nil {
			return nil, ErrProductNotFound
		}

		discountPercentage := 0
		if item.Price.DiscountPercentage != nil {
			discountPercentage = *item.Price.DiscountPercentage
		}

		orderItem := models.OrderItemModel{
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
			VariantId: item.VariantId,
			Price: models.PricingModel{
				Base:               item.Price.Base,
				TaxRate:            item.Price.TaxRate,
				DiscountPercentage: discountPercentage,
			},
		}

		orderItems = append(orderItems, orderItem)
	}

	return orderItems, nil
}

func (s *OrderService) buildOrderModel(payload dto.CreateOrderDto, items []models.OrderItemModel) models.OrderModel {
	return models.OrderModel{
		UserId:   payload.UserId,
		Items:    items,
		Currency: payload.Currency,
		ShippingAddress: models.Address{
			Street:  payload.ShippingAddress.Street,
			City:    payload.ShippingAddress.City,
			State:   payload.ShippingAddress.State,
			ZipCode: payload.ShippingAddress.ZipCode,
			Country: payload.ShippingAddress.Country,
		},
		BillingAddress: models.Address{
			Street:  payload.BillingAddress.Street,
			City:    payload.BillingAddress.City,
			State:   payload.BillingAddress.State,
			ZipCode: payload.BillingAddress.ZipCode,
			Country: payload.BillingAddress.Country,
		},
		Notes: payload.Notes,
	}
}

func (s *OrderService) updateOrderItems(tx *gorm.DB, orderId string, items []dto.UpdateOrderItemDto) error {
	if err := tx.Where("order_id = ?", orderId).Delete(&models.OrderItemModel{}).Error; err != nil {
		return fmt.Errorf("failed to delete existing order items: %w", err)
	}

	productService := NewProductService(tx)
	for _, item := range items {
		if _, err := productService.GetProduct(item.ProductId); err != nil {
			return ErrProductNotFound
		}

		discountPercentage := 0
		if item.Price.DiscountPercentage != nil {
			discountPercentage = *item.Price.DiscountPercentage
		}

		orderItem := models.OrderItemModel{
			OrderId:   uuid.MustParse(orderId),
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
			VariantId: item.VariantId,
			Price: models.PricingModel{
				Base:               item.Price.Base,
				TaxRate:            item.Price.TaxRate,
				DiscountPercentage: discountPercentage,
			},
		}

		if err := tx.Create(&orderItem).Error; err != nil {
			return fmt.Errorf("failed to create order item: %w", err)
		}
	}

	return nil
}

func (s *OrderService) loadOrderRelations(db *gorm.DB, order *models.OrderModel) error {
	return db.
		Preload("Items").
		Preload("User").
		First(order, order.Id).Error
}

func (s *OrderService) normalizePagination(params dto.Pagination) dto.Pagination {
	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit > 100 {
		params.Limit = 100
	}
	return params
}

func (s *OrderService) calculateTotalPages(totalItems int64, limit int) int {
	if limit <= 0 {
		return 0
	}
	return int((totalItems + int64(limit) - 1) / int64(limit))
}

func (s *OrderService) emptyPaginatedResponse(params dto.Pagination) *dto.PaginatedResponse[models.OrderModel] {
	return &dto.PaginatedResponse[models.OrderModel]{
		Data:       []models.OrderModel{},
		Limit:      params.Limit,
		Page:       params.Page,
		TotalItems: 0,
		TotalPages: 0,
	}
}
