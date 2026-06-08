package handler

import (
	"cashflow_gin/api"
	"cashflow_gin/services"
	"cashflow_gin/utils"
	"context"

	"github.com/google/uuid"
)

type TransactionAPI struct {
	Service services.TransactionService
}

func (c *TransactionAPI) GetTransactions(ctx context.Context, request api.GetTransactionsRequestObject) (api.GetTransactionsResponseObject, error) {
	userID, userRole, err := utils.GetUserInfo(ctx)
	if err != nil {

		return api.GetTransactions500JSONResponse{
			Status:  utils.BoolPtr(false),
			Message: utils.StringPtr("Cannot get context"),
			Errors:  utils.StringPtr("ERR: " + err.Error()),
		}, nil
	}

	transactions, err := c.Service.GetAll(ctx, userID, userRole)
	if err != nil {

		return api.GetTransactions500JSONResponse{
			Status:  utils.BoolPtr(false),
			Message: utils.StringPtr("Cannot get transactions"),
			Errors:  utils.StringPtr("ERR: " + err.Error()),
		}, nil
	}

	var res []api.TransactionRes
	for _, trx := range *transactions {
		// Penanganan pointer aman untuk description
		var desc *string
		if trx.Description != "" {
			d := trx.Description
			desc = &d
		}

		res = append(res, api.TransactionRes{
			Id:          trx.ID.String(),
			Amount:      trx.Amount,
			Title:       trx.Title,
			Description: desc,
			Date:        trx.Date,
			Category: api.CategoryRes{
				Id:   trx.Category.ID.String(),
				Name: trx.Category.Name,
				Type: string(trx.Category.Type),
			},
			User: api.UserRes{
				Id:       trx.User.ID.String(),
				Username: trx.User.Username,
				Email:    trx.User.Email,
				UserRole: trx.User.UserRole.String(),
			},
		})
	}

	status := true
	msg := "Success Get All Transactions"
	return api.GetTransactions200JSONResponse{
		Data:    &res,
		Status:  &status,
		Message: &msg,
	}, nil
}

func (c *TransactionAPI) CreateTransaction(ctx context.Context, req api.CreateTransactionRequestObject) (api.CreateTransactionResponseObject, error) {
	userID, _, err := utils.GetUserInfo(ctx)
	if err != nil {

		return api.CreateTransaction400JSONResponse{
			Status:  utils.BoolPtr(false),
			Message: utils.StringPtr("Cannot get context"),
			Errors:  utils.StringPtr("ERR : " + err.Error()),
		}, nil
	}

	inputData := api.CreateTransactionReq(*req.Body)
	newTrx, err := c.Service.Create(ctx, userID, inputData)
	if err != nil {
		return api.CreateTransaction400JSONResponse{
			Status:  utils.BoolPtr(false),
			Message: utils.StringPtr("failed to create transaction: "),
			Errors:  utils.StringPtr("ERR : " + err.Error()),
		}, nil
	}

	var desc *string
	if newTrx.Description != "" {
		d := newTrx.Description
		desc = &d
	}

	res := api.TransactionRes{
		Id:          newTrx.ID.String(),
		Title:       newTrx.Title,
		Amount:      newTrx.Amount,
		Description: desc,
		Date:        newTrx.Date,
		Category: api.CategoryRes{
			Id:   newTrx.Category.ID.String(),
			Name: newTrx.Category.Name,
			Type: string(newTrx.Category.Type),
		},
		User: api.UserRes{
			Id:       newTrx.User.ID.String(),
			Username: newTrx.User.Username,
			Email:    newTrx.User.Email,
			UserRole: newTrx.User.UserRole.String(),
		},
	}

	return api.CreateTransaction201JSONResponse{
		Status:  utils.BoolPtr(true),
		Message: utils.StringPtr("Success Create Transaction"),
		Data:    &res,
	}, nil
}

func (c *TransactionAPI) FindTransactionById(ctx context.Context, request api.FindTransactionByIdRequestObject) (api.FindTransactionByIdResponseObject, error) {
	userID, _, err := utils.GetUserInfo(ctx)
	if err != nil {

		return api.FindTransactionById401JSONResponse{
			Status:  utils.BoolPtr(false),
			Message: utils.StringPtr("Cannot get context"),
			Errors:  utils.StringPtr("ERR : " + err.Error()),
		}, nil
	}

	transactionID, err := uuid.Parse(request.Id)
	if err != nil {
		return api.FindTransactionById400JSONResponse{
			Status:  utils.BoolPtr(false),
			Message: utils.StringPtr("Invalid transaction ID"),
			Errors:  utils.StringPtr("ERR : " + err.Error()),
		}, nil
	}

	trx, err := c.Service.GetTransactionByID(ctx, userID, transactionID)
	if err != nil {
		return api.FindTransactionById400JSONResponse{
			Status:  utils.BoolPtr(false),
			Message: utils.StringPtr("Cannot get transaction"),
			Errors:  utils.StringPtr("ERR : " + err.Error()),
		}, nil
	}

	var desc *string
	if trx.Description != "" {
		d := trx.Description
		desc = &d
	}

	res := api.TransactionRes{
		Id:          trx.ID.String(),
		Amount:      trx.Amount,
		Title:       trx.Title,
		Description: desc,
		Date:        trx.Date,
		Category: api.CategoryRes{
			Id:   trx.Category.ID.String(),
			Name: trx.Category.Name,
			Type: string(trx.Category.Type),
		},
		User: api.UserRes{
			Id:       trx.User.ID.String(),
			Username: trx.User.Username,
			Email:    trx.User.Email,
			UserRole: trx.User.UserRole.String(),
		},
	}

	return api.FindTransactionById200JSONResponse{
		Data:    &res,
		Status:  utils.BoolPtr(true),
		Message: utils.StringPtr("Success Find Transaction"),
	}, nil
}

func (c *TransactionAPI) UpdateTransaction(ctx context.Context, request api.UpdateTransactionRequestObject) (api.UpdateTransactionResponseObject, error) {
	userID, _, err := utils.GetUserInfo(ctx)
	if err != nil {
		return api.UpdateTransaction401JSONResponse{
			Status:  utils.BoolPtr(false),
			Message: utils.StringPtr("Cannot get context"),
			Errors:  utils.StringPtr("ERR : " + err.Error()),
		}, nil
	}

	transactionID, err := uuid.Parse(request.Id)
	if err != nil {
		return api.UpdateTransaction400JSONResponse{
			Status:  utils.BoolPtr(false),
			Message: utils.StringPtr("Invalid transaction ID"),
		}, nil
	}

	inputData := api.UpdateTransactionReq(*request.Body)
	updatedTrx, err := c.Service.UpdateTransaction(ctx, userID, transactionID, inputData)
	if err != nil {
		return api.UpdateTransaction400JSONResponse{
			Status:  utils.BoolPtr(false),
			Message: utils.StringPtr("failed to update transaction"),
			Errors:  utils.StringPtr("ERR : " + err.Error()),
		}, nil
	}

	var desc *string
	if updatedTrx.Description != "" {
		d := updatedTrx.Description
		desc = &d
	}

	res := api.TransactionRes{
		Id:          updatedTrx.ID.String(),
		Amount:      updatedTrx.Amount,
		Title:       updatedTrx.Title,
		Description: desc,
		Date:        updatedTrx.Date,
		Category: api.CategoryRes{
			Id:   updatedTrx.Category.ID.String(),
			Name: updatedTrx.Category.Name,
			Type: string(updatedTrx.Category.Type),
		},
		User: api.UserRes{
			Id:       updatedTrx.User.ID.String(),
			Username: updatedTrx.User.Username,
			Email:    updatedTrx.User.Email,
			UserRole: updatedTrx.User.UserRole.String(),
		},
	}

	return api.UpdateTransaction200JSONResponse{
		Data:    &res,
		Status:  utils.BoolPtr(true),
		Message: utils.StringPtr("Success Update Transaction"),
	}, nil
}

func (c *TransactionAPI) DeleteTransaction(ctx context.Context, request api.DeleteTransactionRequestObject) (api.DeleteTransactionResponseObject, error) {
	userID, _, err := utils.GetUserInfo(ctx)
	if err != nil {
		return api.DeleteTransaction401JSONResponse{
			Status:  utils.BoolPtr(false),
			Message: utils.StringPtr("Cannot get context"),
			Errors:  utils.StringPtr("ERR : " + err.Error()),
		}, nil
	}

	transactionID, err := uuid.Parse(request.Id)
	if err != nil {
		return api.DeleteTransaction400JSONResponse{
			Status:  utils.BoolPtr(false),
			Message: utils.StringPtr("Invalid transaction ID"),
			Errors:  utils.StringPtr("ERR : " + err.Error()),
		}, nil
	}

	// 🛑 TUGAS MUTLAK LU DI SERVICE:
	// Pastikan fungsi ini di Service TIDAK LAGI meminta parameter walletID.
	err = c.Service.SoftDeleteTransaction(ctx, userID, transactionID)
	if err != nil {
		return api.DeleteTransaction400JSONResponse{
			Status:  utils.BoolPtr(false),
			Message: utils.StringPtr("failed to delete transaction"),
			Errors:  utils.StringPtr("ERR : " + err.Error()),
		}, nil
	}

	return api.DeleteTransaction200JSONResponse{
		Status:  utils.BoolPtr(true),
		Message: utils.StringPtr("Success Delete Transaction"),
	}, nil
}
