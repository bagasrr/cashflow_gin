package controllers

import (
	"cashflow_gin/dto/request"
	"cashflow_gin/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	service services.TransactionService
}

func NewTransactionController(service services.TransactionService) *TransactionController {
	return &TransactionController{service: service}
}

// Create godoc
// @Summary      Create Transaction
// @Description  Membuat transaksi baru.
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Param        request body request.CreateTransactionRequest true "request body"
// @Success      201 {object} response.BaseResponse{data=response.TransactionResponse}
// @Failure      400 {object} response.BaseResponse
// @Failure      500 {object} response.BaseResponse
// @Security 	 BearerAuth
// @Router       /transactions [post]
func (c *TransactionController) Create(ctx *gin.Context) {
	var input request.CreateTransactionRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		c.sendError(ctx, http.StatusBadRequest, "Invalid input data", err)
		return
	}

	// 1 baris untuk ambil UserID
	userID, err := c.getUserID(ctx)
	if err != nil {
		c.sendError(ctx, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	newTransaction, err := c.service.Create(userID, input)
	if err != nil {
		c.sendError(ctx, http.StatusInternalServerError, "Failed to create transaction", err)
		return
	}

	c.sendSuccess(ctx, "Transaction created successfully", newTransaction)
}

// FindAll godoc
// @Summary      Find All Transactions
// @Description  Mendapatkan semua transaksi.
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Success      200 {object} response.BaseResponse{data=[]response.TransactionResponse}
// @Failure      500 {object} response.BaseResponse
// @Security 	 BearerAuth
// @Router       /transactions [get]
func (c *TransactionController) FindAll(ctx *gin.Context) {
	// Note: Harusnya FindAll juga butuh userID kan? Transaction itu private per user.
	// Tapi aku ikutin logic code aslimu dulu.
	transactions, err := c.service.GetAll()
	if err != nil {
		c.sendError(ctx, http.StatusInternalServerError, "Failed to retrieve transactions", err)
		return
	}

	c.sendSuccess(ctx, "Transactions retrieved successfully", transactions)
}

// GetTransactionByID godoc
// @Summary      Get Transaction By ID
// @Description  Mendapatkan transaksi berdasarkan ID.
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Param        id path string true "Transaction ID"
// @Success      200 {object} response.BaseResponse{data=response.TransactionResponse}
// @Failure      400 {object} response.BaseResponse
// @Failure      500 {object} response.BaseResponse
// @Security 	 BearerAuth
// @Router       /transactions/{id} [get]
func (c *TransactionController) GetTransactionByID(ctx *gin.Context) {
	// Reuse helper getParamID
	transactionID, err := c.getParamID(ctx, "id")
	if err != nil {
		c.sendError(ctx, http.StatusBadRequest, "Invalid transaction ID", err)
		return
	}

	userID, err := c.getUserID(ctx)
	if err != nil {
		c.sendError(ctx, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	transaction, err := c.service.GetTransactionByID(userID, transactionID)
	if err != nil {
		c.sendError(ctx, http.StatusInternalServerError, "Failed to retrieve transaction", err)
		return
	}

	c.sendSuccess(ctx, "Transaction retrieved successfully", transaction)
}

// UpdateTransaction godoc
// @Summary      Update Transaction
// @Description  Memperbarui transaksi berdasarkan ID.
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Param        id path string true "Transaction ID"
// @Param        request body request.UpdateTransactionRequest true "request body"
// @Success      200 {object} response.BaseResponse{data=response.TransactionResponse}
// @Failure      400 {object} response.BaseResponse
// @Failure      500 {object} response.BaseResponse
// @Security 	 BearerAuth
// @Router       /transactions/{id}/update [patch]
func (c *TransactionController) UpdateTransaction(ctx *gin.Context) {
	transactionID, err := c.getParamID(ctx, "id")
	if err != nil {
		c.sendError(ctx, http.StatusBadRequest, "Invalid transaction ID", err)
		return
	}

	var input request.UpdateTransactionRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		c.sendError(ctx, http.StatusBadRequest, "Invalid input format", err)
		return
	}

	userID, err := c.getUserID(ctx)
	if err != nil {
		c.sendError(ctx, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	updatedTransaction, err := c.service.UpdateTransaction(userID, transactionID, input)
	if err != nil {
		c.sendError(ctx, http.StatusInternalServerError, "Failed to update transaction", err)
		return
	}

	c.sendSuccess(ctx, "Transaction updated successfully", updatedTransaction)
}

// SoftDeleteTransaction godoc
// @Summary      Soft Delete Transaction
// @Description  Melakukan soft delete pada transaksi berdasarkan ID.
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Param        id path string true "Transaction ID"
// @Param        walletid path string true "Wallet ID"
// @Success      200 {object} response.BaseResponse
// @Failure      400 {object} response.BaseResponse
// @Failure      500 {object} response.BaseResponse
// @Security 	 BearerAuth
// @Router       /transactions/{id}/wallet/{walletid}/soft-delete [patch]
func (c *TransactionController) SoftDeleteTransaction(ctx *gin.Context) {
	transactionID, err := c.getParamID(ctx, "id")
	if err != nil {
		c.sendError(ctx, http.StatusBadRequest, "Invalid transaction ID", err)
		return
	}

	walletID, err := c.getParamID(ctx, "walletid")
	if err != nil {
		c.sendError(ctx, http.StatusBadRequest, "Invalid wallet ID", err)
		return
	}

	userID, err := c.getUserID(ctx)
	if err != nil {
		c.sendError(ctx, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	err = c.service.SoftDeleteTransaction(userID, transactionID, walletID)
	if err != nil {
		c.sendError(ctx, http.StatusInternalServerError, "Failed to soft delete transaction", err)
		return
	}

	c.sendSuccess(ctx, "Transaction soft deleted successfully", nil)
}
