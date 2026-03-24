package handler

import (
	"cashflow_gin/api"
	"cashflow_gin/services"
	"context"
)

type WalletAPI struct {
	Service services.WalletService
}

func (c *WalletAPI) CreateWallet(ctx context.Context, request api.CreateWalletRequestObject) (api.CreateWalletResponseObject, error) {
	return api.CreateWallet201JSONResponse{}, nil
}

func (c *WalletAPI) DeleteWallet(ctx context.Context, request api.DeleteWalletRequestObject) (api.DeleteWalletResponseObject, error) {
	return api.DeleteWallet200JSONResponse{}, nil
}

func (c *WalletAPI) GetWalletById(ctx context.Context, request api.GetWalletByIdRequestObject) (api.GetWalletByIdResponseObject, error) {
	return api.GetWalletById200JSONResponse{}, nil
}

func (c *WalletAPI) GetMyWallets(ctx context.Context, request api.GetMyWalletsRequestObject) (api.GetMyWalletsResponseObject, error) {
	return api.GetMyWallets200JSONResponse{}, nil
}

func (c *WalletAPI) UpdateWallet(ctx context.Context, request api.UpdateWalletRequestObject) (api.UpdateWalletResponseObject, error) {
	return api.UpdateWallet200JSONResponse{}, nil
}
