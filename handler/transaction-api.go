package handler

import (
	"cashflow_gin/api"
	"cashflow_gin/services"
	"context"
)

type TransactionAPI struct {
	Service services.TransactionService
}

func (c *TransactionAPI) GetTransactions(ctx context.Context, request api.GetTransactionsRequestObject) (api.GetTransactionsResponseObject, error) {
	return api.GetTransactions200JSONResponse{}, nil
}

func (c *TransactionAPI) CreateTransaction(ctx context.Context, request api.CreateTransactionRequestObject) (api.CreateTransactionResponseObject, error) {
	return api.CreateTransaction201JSONResponse{}, nil
}

func (c *TransactionAPI) DeleteTransaction(ctx context.Context, request api.DeleteTransactionRequestObject) (api.DeleteTransactionResponseObject, error) {
	return api.DeleteTransaction200JSONResponse{}, nil
}

func (c *TransactionAPI) FindTransactionById(ctx context.Context, request api.FindTransactionByIdRequestObject) (api.FindTransactionByIdResponseObject, error) {
	return api.FindTransactionById200JSONResponse{}, nil
}

func (c *TransactionAPI) UpdateTransaction(ctx context.Context, request api.UpdateTransactionRequestObject) (api.UpdateTransactionResponseObject, error) {
	return api.UpdateTransaction201JSONResponse{}, nil
}
