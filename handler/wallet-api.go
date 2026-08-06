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
	Service            services.WalletService
	TransactionService services.TransactionService
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
				Type: string(v.Category.Type),
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
	res := []api.WalletRes{}

	// 1. LOOPING DOMPET UTAMA
	for _, v := range myWallets {

		var trxRes []api.TransactionRes

		// 2. LOOPING TRANSAKSI ANAK
		for _, t := range v.Transactions {

			// A. RAKIT KATEGORI SECARA AMAN (Seluruhnya ada di dalam blok if)
			var catRes api.CategoryRes
			if t.Category.ID != uuid.Nil {
				catRes = api.CategoryRes{
					Id:      t.Category.ID.String(),
					Name:    t.Category.Name,
					Type:    string(t.Category.Type),
					GroupId: utils.UUIDPtrToStringPtr(t.Category.GroupID),
				}
			}

			// B. RAKIT USER SECARA AMAN (Wajib lu tambahin kalau Front-end butuh data siapa yang transaksi)
			// Asumsi nama struct dari OpenAPI lu adalah api.UserRes
			/*
			   var userRes api.UserRes
			   if t.User.ID != uuid.Nil {
			      userRes = api.UserRes{
			          Id:       t.User.ID.String(),
			          Username: t.User.Username,
			          Email:    t.User.Email,
			      }
			   }
			*/

			// C. TEMPELKAN KE WADAH TRANSAKSI
			trxRes = append(trxRes, api.TransactionRes{
				Id:       t.ID.String(),
				Title:    t.Title,
				Amount:   t.Amount,
				Date:     t.Date, // KITA BUKA MUTLAK AGAR TIDAK 0001-01-01
				Category: catRes, // Objek kategori yang sudah aman

				// User: userRes,  // Buka ini kalau struct TransactionRes lu butuh data user
			})
		}

		// 3. TEMPELKAN KE INDUK (WALLET)
		res = append(res, api.WalletRes{
			Id:               v.ID.String(),
			Name:             v.Name,
			UserId:           utils.UUIDPtrToStringPtr(v.UserID),
			GroupId:          utils.UUIDPtrToStringPtr(v.GroupID),
			Balance:          v.Balance,
			TransactionCount: v.TransactionCount,
			Transactions:     trxRes,
		})
	}

	// HAPUS fmt.Println(res)
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

func (c *WalletAPI) GetWalletChartData(ctx context.Context, request api.GetWalletChartDataRequestObject) (api.GetWalletChartDataResponseObject, error) {
	userId, _, err := utils.GetUserInfo(ctx)

	if err != nil {
		return api.GetWalletChartData200JSONResponse{
			Status:  utils.BoolPtr(false),
			Message: utils.StringPtr("Cannot get Context"),
		}, nil
	}

	walletId, err := uuid.Parse(request.Id)
	if err != nil {
		return api.GetWalletChartData200JSONResponse{
			Status:  utils.BoolPtr(false),
			Message: utils.StringPtr("Invalid wallet ID format"),
		}, nil
	}

	// Ambil start_date & end_date dari request.Params (Gunakan logika default fallback 30 hari seperti sebelumnya)
	// Lalu panggil service -> repo
	points, err := c.Service.GetWalletChartData(ctx, userId, walletId, request.Params)
	if err != nil {
		return api.GetWalletChartData200JSONResponse{
			Status:  utils.BoolPtr(false),
			Message: utils.StringPtr("Get Chart Data Failed"),
		}, nil
	}

	var chartRes []api.WalletChartPoint
	for _, p := range points {
		chartRes = append(chartRes, api.WalletChartPoint{
			Date:       utils.StringPtr(p.Date.Format("2006-01-02")), // Ubah ke string YYYY-MM-DD
			Income:     utils.IntPtr(int(p.Income)),
			Expense:    utils.IntPtr(int(p.Expense)),
			Investment: utils.IntPtr(int(p.Investment)),
		})
	}

	return api.GetWalletChartData200JSONResponse{
		Status:  utils.BoolPtr(true),
		Message: utils.StringPtr("Success get wallet chart data"),
		Data:    &chartRes,
	}, nil
}

func (h *WalletAPI) GetWalletTransactions(ctx context.Context, request api.GetWalletTransactionsRequestObject) (api.GetWalletTransactionsResponseObject, error) {
	userId, _, err := utils.GetUserInfo(ctx)
	if err != nil {
		return api.GetWalletTransactions401JSONResponse{
			Status:  utils.BoolPtr(false),
			Message: utils.StringPtr("Cannot get user info"),
			Errors:  utils.StringPtr("ERR : " + err.Error()),
		}, nil
	}

	// 2. Validasi Wallet ID dari URL Path
	walletId, err := uuid.Parse(request.Id)
	if err != nil {
		return api.GetWalletTransactions400JSONResponse{
			Status:  utils.BoolPtr(false),
			Message: utils.StringPtr("Cannot get user id"),
			Errors:  utils.StringPtr("ERR : " + err.Error()),
		}, nil
	}
	page, limit, offset := utils.ValidatePagination(*request.Params.Page, *request.Params.Limit)
	// 3. Panggil Service
	transactions, totalData, err := h.TransactionService.GetTransactionsByWallet(ctx, userId, walletId, request.Params, page, limit, offset)
	if err != nil {
		return api.GetWalletTransactions500JSONResponse{
			Status:  utils.BoolPtr(false),
			Message: utils.StringPtr("Cannot get Transaction"),
			Errors:  utils.StringPtr("ERR : " + err.Error()),
		}, nil
	}
	// 4. Mapping Data Database ke Data OpenAPI
	var dataRes []api.TxWithWallet
	for _, trx := range transactions {
		dataRes = append(dataRes, api.TxWithWallet{
			Id:          trx.ID.String(),
			Amount:      int64(trx.Amount),
			Title:       trx.Title,
			Date:        trx.Date,
			Description: &trx.Description,
			Category: api.CategoryRes{
				Id:   trx.Category.ID.String(),
				Name: trx.Category.Name,
				Type: string(trx.Category.Type),
			},
			Wallet: &api.TxWalletRes{
				Id:   trx.Wallet.ID.String(),
				Name: trx.Wallet.Name,
			},
		})
	}

	totalPages := (int(totalData) + limit - 1) / limit

	// 5. Kembalikan 200 OK
	return api.GetWalletTransactions200JSONResponse{
		Status:  utils.BoolPtr(true),
		Message: utils.StringPtr("Berhasil mengambil riwayat transaksi"),
		Data:    &dataRes,
		Meta: &api.PaginationMeta{
			CurrentPage: utils.IntPtr(page),
			TotalPages:  utils.IntPtr(totalPages),
			TotalItems:  utils.IntPtr(int(totalData)),
		},
	}, nil
}
