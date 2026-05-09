package handler

import (
	"cashflow_gin/api"
	"cashflow_gin/services"
	"cashflow_gin/utils"
	"context"
)

type TransactionAPI struct {
	Service services.TransactionService
}

func (c *TransactionAPI) GetTransactions(ctx context.Context, request api.GetTransactionsRequestObject) (api.GetTransactionsResponseObject, error) {
	userID, userRole, err := utils.GetUserInfo(ctx)
	if err != nil {
		status := false
		msg := "Gagal Auth: " + err.Error()
		return api.GetTransactions500JSONResponse{
			Status:  &status,
			Message: &msg,
		}, nil
	}

	transactions, err := c.Service.GetAll(ctx, userID, userRole)
	if err != nil {
		status := false
		msg := "Gagal Auth: " + err.Error()
		return api.GetTransactions500JSONResponse{
			Status:  &status,
			Message: &msg,
		}, nil
	}

	Status := true
	Message := "Get Transactions Success"
	var res []api.TransactionRes
	for _, v := range *transactions {
		res = append(res, api.TransactionRes{
			Amount: &v.Amount,
			Category: &api.CategoryRes{
				Id:      &v.Category.ID,
				Name:    &v.Category.Name,
				GroupId: &v.Category.GroupID,
				Type:    &v.Category.Type,
			},
			Date:        &v.Date,
			Description: &v.Description,
			Id:          &v.ID,
			Title:       &v.Title,
			User: &api.UserRes{
				Id:       &v.User.ID,
				Email:    &v.User.Email,
				Username: &v.User.Username,
			},
		})
	}
	return api.GetTransactions200JSONResponse{
		Data:    &res,
		Status:  &Status,
		Message: &Message,
	}, nil
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
