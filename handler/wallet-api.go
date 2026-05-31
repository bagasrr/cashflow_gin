package handler

import (
	"cashflow_gin/api"
	"cashflow_gin/models"
	"cashflow_gin/services"
	"cashflow_gin/utils"
	"context"

	"github.com/google/uuid"
)

type WalletAPI struct {
	Service services.WalletService
}

func (c *WalletAPI) CreatePersonalWallet(ctx context.Context, request api.CreatePersonalWalletRequestObject) (api.CreatePersonalWalletResponseObject, error) {
	userId, _, err := utils.GetUserInfo(ctx)
	if err != nil {
		return api.CreatePersonalWallet400JSONResponse{
			Message: utils.StringPtr("Cannot get Context"),
			Status:  utils.BoolPtr(false),
			Errors:  utils.StringPtr("Err : " + err.Error()),
		}, nil
	}
	var wallet models.Wallet
	wallet.Name = request.Body.Name
	wallet.UserID = &userId
	wallet.GroupID = nil
	wallet.Balance = 0
	wallet.Currency = "IDR"

	createWallet, err := c.Service.CreatePersonalWallet(ctx, wallet)
	if err != nil {
		return api.CreatePersonalWallet500JSONResponse{
			Message: utils.StringPtr("Cannot Create Wallet"),
			Status:  utils.BoolPtr(false),
			Errors:  utils.StringPtr("Err : " + err.Error()),
		}, nil
	}
	var res api.WalletRes
	res.Id = createWallet.ID.String()
	res.Name = createWallet.Name
	res.GroupId = nil
	res.Balance = createWallet.Balance
	return api.CreatePersonalWallet201JSONResponse{
		Status:  utils.BoolPtr(true),
		Message: utils.StringPtr("Create Wallet Successfully"),
		Data:    &res,
	}, nil
}

func (c *WalletAPI) CreateGroupWallet(ctx context.Context, request api.CreatePersonalWalletRequestObject) (api.CreatePersonalWalletResponseObject, error) {
	groupId, err := uuid.Parse(*request.Body.GroupId)
	if err != nil {
		return api.CreatePersonalWallet400JSONResponse{
			Message: utils.StringPtr("Cannot get Context"),
			Status:  utils.BoolPtr(false),
			Errors:  utils.StringPtr("Err : " + err.Error()),
		}, nil
	}

	var wallet models.Wallet
	wallet.Name = request.Body.Name
	wallet.UserID = nil
	wallet.GroupID = &groupId
	wallet.Balance = 0
	wallet.Currency = "IDR"

	createWallet, err := c.Service.CreatePersonalWallet(ctx, wallet)
	if err != nil {
		return api.CreatePersonalWallet500JSONResponse{
			Message: utils.StringPtr("Cannot Create Wallet"),
			Status:  utils.BoolPtr(false),
			Errors:  utils.StringPtr("Err : " + err.Error()),
		}, nil
	}
	var res api.WalletRes
	res.Id = createWallet.ID.String()
	res.Name = createWallet.Name
	res.GroupId = nil
	res.Balance = createWallet.Balance
	return api.CreatePersonalWallet201JSONResponse{
		Status:  utils.BoolPtr(true),
		Message: utils.StringPtr("Create Wallet Successfully"),
		Data:    &res,
	}, nil
}
func (c *WalletAPI) DeleteWallet(ctx context.Context, request api.DeleteWalletRequestObject) (api.DeleteWalletResponseObject, error) {
	userId, _, err := utils.GetUserInfo(ctx)
	if err != nil {
		return api.DeleteWallet400JSONResponse{
			Status:  utils.BoolPtr(false),
			Message: utils.StringPtr("Cannot get Context"),
			Errors:  utils.StringPtr("Err : " + err.Error()),
		}, nil
	}
	walletId, err := uuid.Parse(request.Id)
	if err != nil {
		return api.DeleteWallet400JSONResponse{
			Status:  utils.BoolPtr(false),
			Message: utils.StringPtr("Cannot get Params"),
			Errors:  utils.StringPtr("Err : " + err.Error()),
		}, nil
	}
	err = c.Service.DeleteWallet(ctx, walletId, userId)
	if err != nil {
		return api.DeleteWallet500JSONResponse{
			Status:  utils.BoolPtr(false),
			Message: utils.StringPtr("Delete Failed"),
			Errors:  utils.StringPtr("Err : " + err.Error()),
		}, nil
	}
	return api.DeleteWallet200JSONResponse{
		Status:  utils.BoolPtr(true),
		Message: utils.StringPtr("Delete Successfully"),
	}, nil
}

func (c *WalletAPI) GetWalletById(ctx context.Context, request api.GetWalletByIdRequestObject) (api.GetWalletByIdResponseObject, error) {
	return api.GetWalletById200JSONResponse{}, nil
}

func (c *WalletAPI) GetMyWallets(ctx context.Context, request api.GetMyWalletsRequestObject) (api.GetMyWalletsResponseObject, error) {
	return api.GetMyWallets200JSONResponse{}, nil
}

func (c *WalletAPI) UpdateWallet(ctx context.Context, request api.UpdateWalletRequestObject) (api.UpdateWalletResponseObject, error) {
	return api.UpdateWallet201JSONResponse{}, nil
}
