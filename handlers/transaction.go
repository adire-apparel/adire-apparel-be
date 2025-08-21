package handlers

import (
	"adire-apparel/database"
	"adire-apparel/dto"
	"adire-apparel/lib"
	"adire-apparel/services"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	service *services.TransactionService
}

func NewTransactionHandler() *TransactionHandler {
	return &TransactionHandler{
		service: services.NewTransactionService(database.GetDatabase()),
	}
}

func (h TransactionHandler) CreateTransaction() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload dto.CreateTransactionDto

		if err := ctx.ShouldBind(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		transaction := dto.CreateTransactionDto{
			UserId:      payload.UserId,
			Type:        payload.Type,
			Status:      payload.Status,
			Amount:      payload.Amount,
			Currency:    payload.Currency,
			Description: payload.Description,
			OrderId:     payload.OrderId,
			Metadata:    payload.Metadata,
			CreatedBy:   payload.CreatedBy,
		}

		created, err := h.service.CreateTransaction(&transaction)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Transaction created successfully", created)
	}
}

func (h TransactionHandler) UpdateTransaction() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload dto.UpdateTransactionDto
		id := ctx.Param("id")

		if err := ctx.ShouldBind(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		transaction := dto.UpdateTransactionDto{
			Status:   payload.Status,
			Metadata: payload.Metadata,
		}

		saved, err := h.service.UpdateTransaction(id, &transaction)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Transaction created successfully", saved)
	}
}

func (h TransactionHandler) GetTransactions() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var params dto.Pagination

		if err := ctx.ShouldBind(&params); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		transactions, err := h.service.GetTransactions(&params)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Transaction created successfully", transactions)
	}
}

func (h TransactionHandler) GetTransactionsByUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		var params dto.Pagination

		if err := ctx.ShouldBind(&params); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		transactions, err := h.service.GetTransactionsByUser(id, &params)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Transaction created successfully", transactions)
	}
}

func (h TransactionHandler) GetTransactionsByAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		var params dto.Pagination

		if err := ctx.ShouldBind(&params); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		transactions, err := h.service.GetTransactionsByAdmin(id, &params)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Transaction created successfully", transactions)
	}
}

func (h TransactionHandler) GetTransaction() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		transaction, err := h.service.GetTransaction(id)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Transaction created successfully", transaction)
	}
}

func (h TransactionHandler) DeleteTransaction() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		err := h.service.DeleteTransaction(id)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Transaction created successfully", nil)
	}
}
