package handler

import (
	"cashflow_gin/api"
	"cashflow_gin/models"
	"cashflow_gin/services"
	"cashflow_gin/utils"
	"context"
	"fmt"

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
	res.GroupId = utils.UUIDPtrToStringPtr(createWallet.GroupID)
	res.UserId = utils.UUIDPtrToStringPtr(createWallet.UserID)
	res.Balance = createWallet.Balance
	return api.CreatePersonalWallet201JSONResponse{
		Status:  utils.BoolPtr(true),
		Message: utils.StringPtr("Create Wallet Successfully"),
		Data:    &res,
	}, nil
}

func (c *WalletAPI) CreateGroupWallet(ctx context.Context, request api.CreateGroupWalletRequestObject) (api.CreateGroupWalletResponseObject, error) {
	groupId, err := uuid.Parse(*request.Body.GroupId)
	if err != nil {
		return api.CreateGroupWallet400JSONResponse{
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

	createWallet, err := c.Service.CreateGroupWallet(ctx, wallet)
	if err != nil {
		return api.CreateGroupWallet500JSONResponse{
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
	return api.CreateGroupWallet201JSONResponse{
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
	userId, _, err := utils.GetUserInfo(ctx)
	walletId, err := uuid.Parse(request.Id)
	if err != nil {
		return api.GetWalletById400JSONResponse{
			Status:  utils.BoolPtr(false),
			Message: utils.StringPtr("Cannot parse Params"),
			Errors:  utils.StringPtr("Err : " + err.Error()),
		}, nil
	}
	wallet, err := c.Service.GetWalletByID(ctx, userId, walletId)
	if err != nil {
		return api.GetWalletById500JSONResponse{
			Status:  utils.BoolPtr(false),
			Message: utils.StringPtr("Get Wallet Failed"),
			Errors:  utils.StringPtr("Err : " + err.Error()),
		}, nil
	}
	var walletTransac []api.TransactionRes
	for _, v := range wallet.Transactions {
		walletTransac = append(walletTransac, api.TransactionRes{
			Id:          v.ID.String(),
			Title:       v.Title,
			Amount:      v.Amount,
			Description: &v.Description,
			Category: api.CategoryRes{
				Id:   v.Category.ID.String(),
				Name: v.Category.Name,
				Type: v.Category.Type,
			},
			User: api.UserRes{
				Id:       v.User.ID.String(),
				Email:    v.User.Email,
				Username: v.User.Username,
				UserRole: v.User.UserRole.String(),
			},
		})
	}
	var res api.WalletRes
	res.Id = wallet.ID.String()
	res.Name = wallet.Name
	res.GroupId = utils.UUIDPtrToStringPtr(wallet.GroupID)
	res.UserId = utils.UUIDPtrToStringPtr(wallet.UserID)
	res.Balance = wallet.Balance
	res.TransactionCount = wallet.TransactionCount
	res.Transactions = walletTransac

	return api.GetWalletById200JSONResponse{
		Message: utils.StringPtr("Get Wallet Successfully"),
		Status:  utils.BoolPtr(true),
		Data:    &res,
	}, nil
}

func (c *WalletAPI) GetMyWallets(ctx context.Context, request api.GetMyWalletsRequestObject) (api.GetMyWalletsResponseObject, error) {
	userId, _, err := utils.GetUserInfo(ctx)
	if err != nil {
		return api.GetMyWallets401JSONResponse{
			Message: utils.StringPtr("Cannot get Context"),
			Status:  utils.BoolPtr(false),
			Errors:  utils.StringPtr("Err : " + err.Error()),
		}, nil
	}
	page, limit, offset := utils.ValidatePagination(request.Params.Page, request.Params.Limit)
	myWallets, totalItems, err := c.Service.GetMine(ctx, userId, limit, offset)
	if err != nil {
		return api.GetMyWallets500JSONResponse{
			Message: utils.StringPtr("Get Wallets Failed"),
			Status:  utils.BoolPtr(false),
			Errors:  utils.StringPtr("Err : " + err.Error()),
		}, nil
	}
	fmt.Println(myWallets)
	res := []api.WalletRes{}

	// 2. LOOPING & MAPPING
	for _, v := range *myWallets {
		res = append(res, api.WalletRes{
			Id:               v.ID.String(),
			Name:             v.Name,
			UserId:           utils.UUIDPtrToStringPtr(v.UserID),
			GroupId:          utils.UUIDPtrToStringPtr(v.GroupID),
			Balance:          v.Balance,
			TransactionCount: v.TransactionCount,

			// 3. PROTEKSI BERSARANG: Cegah properti ini meledak menjadi 'null' di Frontend
			Transactions: []api.TransactionRes{},
		})
	}
	fmt.Println(res)
	totalPages := (int(totalItems) + limit - 1) / limit
	return api.GetMyWallets200JSONResponse{
		Message: utils.StringPtr("Get Wallets Successfully"),
		Status:  utils.BoolPtr(true),
		Data:    &res,
		Meta: &api.PaginationMeta{
			CurrentPage: utils.IntPtr(page),
			TotalPages:  utils.IntPtr(totalPages),
			TotalItems:  utils.IntPtr(int(totalItems)),
		},
	}, nil
}

func (c *WalletAPI) UpdateWallet(ctx context.Context, request api.UpdateWalletRequestObject) (api.UpdateWalletResponseObject, error) {
	// 1. Ambil KTP User (ID)
	userId, _, err := utils.GetUserInfo(ctx)
	if err != nil {
		return api.UpdateWallet401JSONResponse{
			Status:  utils.BoolPtr(false),
			Message: utils.StringPtr("Unauthorized"),
			Errors:  utils.StringPtr("Err : " + err.Error()),
		}, nil
	}

	// 2. Parsing Wallet ID dari URL (Ubah string ke UUID)
	walletId, err := uuid.Parse(request.Id)
	if err != nil {
		return api.UpdateWallet400JSONResponse{
			Status:  utils.BoolPtr(false),
			Message: utils.StringPtr("Invalid wallet ID format"),
			Errors:  utils.StringPtr("Err : " + err.Error()),
		}, nil
	}

	// 3. Lempar ke Service (Asumsi field di JSON body lu namanya 'Name')
	updatedWallet, err := c.Service.UpdateWalletName(ctx, userId, walletId, request.Body.Name)
	if err != nil {
		return api.UpdateWallet500JSONResponse{
			Status:  utils.BoolPtr(false),
			Message: utils.StringPtr("Failed to update wallet"),
			Errors:  utils.StringPtr("Err : " + err.Error()),
		}, nil
	}

	// 4. MAPPING RESPONSE
	var res api.WalletRes
	res.Id = updatedWallet.ID.String()
	res.Name = updatedWallet.Name
	res.Balance = updatedWallet.Balance
	res.TransactionCount = updatedWallet.TransactionCount

	// Gunakan helper/manual check biar gak kena Nil Pointer Panic
	if updatedWallet.GroupID != nil {
		res.GroupId = utils.StringPtr(updatedWallet.GroupID.String())
	} else {
		res.GroupId = nil
	}

	// Set transaksi kosong karena kita gak nge-preload transaksi pas update
	res.Transactions = []api.TransactionRes{}

	// RETURN 200 OK
	return api.UpdateWallet200JSONResponse{
		Message: utils.StringPtr("Wallet updated successfully"),
		Status:  utils.BoolPtr(true),
		Data:    &res,
	}, nil
}
