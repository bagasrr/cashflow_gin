package handler

import (
	"cashflow_gin/api"
	"cashflow_gin/models"
	"cashflow_gin/services"
	"cashflow_gin/utils"
	"context"
)

type GroupAPI struct {
	Service services.GroupService
}

func (c *GroupAPI) GetGroups(ctx context.Context, request api.GetGroupsRequestObject) (api.GetGroupsResponseObject, error) {
	_, role, err := utils.GetUserInfo(ctx)
	if err != nil {
		return api.GetGroups401JSONResponse{
			Message: utils.StringPtr("Failed to get user info : " + err.Error()),
			Status:  utils.BoolPtr(false),
		}, err
	}
	if role != models.RoleAdmin {
		return api.GetGroups401JSONResponse{
			Message: utils.StringPtr("You do not have access to this resource"),
			Status:  utils.BoolPtr(false),
		}, err
	}
	limit, page, offset := utils.ValidatePagination(request.Params.Limit, request.Params.Page)
	groups, totalItem, err := c.Service.GetAllGroups(ctx, limit, offset)
	if err != nil {
		return api.GetGroups500JSONResponse{
			Message: utils.StringPtr("Failed to get groups : " + err.Error()),
			Status:  utils.BoolPtr(false),
		}, err
	}
	var res []api.GroupBaseRes

	// LOOPING PARENT (Group)
	for _, group := range *groups {

		// 1. SIAPKAN SLICE ANAK DI DALAM SINI (Scope Terisolasi)
		var walletRes []api.WalletRes
		var membersRes []api.GroupMembersRes

		// 2. LOOPING ANAK PERTAMA (Wallets)
		// Pastikan 'group.Wallets' sesuai dengan nama field relasi di Model GORM lu
		for _, w := range group.Wallet {
			walletRes = append(walletRes, api.WalletRes{
				Id:               w.ID.String(),
				Name:             w.Name,
				Balance:          w.Balance, // Pastikan tipe datanya match (int64)
				TransactionCount: w.TransactionCount,
				// Gunakan helper pointer kalau GroupId di Wallet opsional
				GroupId: utils.UUIDPtrToStringPtr(w.GroupID),
				// Kosongkan dulu list transaksi untuk di wallet ini biar ringan
				Transactions: []api.TransactionRes{},
			})
		}

		// 3. LOOPING ANAK KEDUA (Members)
		// Pastikan 'group.Members' sesuai dengan nama field relasi di Model GORM lu
		for _, m := range group.Members {
			membersRes = append(membersRes, api.GroupMembersRes{
				Id:     m.ID.String(), // ID relasi pivot
				UserId: m.UserID.String(),
				Role:   m.MembersRole.String(),
			})
		}

		// 4. BUNGKUS DAN GABUNGKAN KE PARENT
		res = append(res, api.GroupBaseRes{
			Id:   group.ID.String(),
			Name: group.Name,
			// Karena description itu opsional di YAML, bungkus pake pointer
			Description: utils.StringPtr(group.Description),
			Wallet:      walletRes,
			Members:     membersRes,
		})
	}

	totalPages := (int(totalItem) + limit - 1) / limit
	// 5. KEMBALIKAN RESPONSE 200 (Bukan 201, GET itu 200 Success)
	return api.GetGroups201JSONResponse{
		Data:    &res,
		Status:  utils.BoolPtr(true),
		Message: utils.StringPtr("Get Group Success"),
		Meta: &api.PaginationMeta{
			CurrentPage: utils.IntPtr(page),
			TotalPages:  utils.IntPtr(totalPages),
			TotalItems:  utils.IntPtr(int(totalItem)),
		},
	}, nil
}

func (c *GroupAPI) GetMyGroups(ctx context.Context, request api.GetMyGroupsRequestObject) (api.GetMyGroupsResponseObject, error) {
	userId, _, err := utils.GetUserInfo(ctx)
	if err != nil {
		return api.GetMyGroups401JSONResponse{
			Message: utils.StringPtr("Failed to get user info : " + err.Error()),
			Status:  utils.BoolPtr(false),
		}, err
	}

	page, limit, offset := utils.ValidatePagination(request.Params.Page, request.Params.Limit)
	myGroups, totalData, err := c.Service.GetMyGroups(ctx, page, offset, userId)

	if err != nil {
		return api.GetMyGroups500JSONResponse{
			Message: utils.StringPtr("Failed to get groups : " + err.Error()),
			Status:  utils.BoolPtr(false),
		}, err
	}
	var res []api.GroupBaseRes
	for _, group := range *myGroups {
		res = append(res, api.GroupBaseRes{
			Id:          group.ID.String(),
			Name:        group.Name,
			Description: utils.StringPtr(group.Description),
			Wallet:      []api.WalletRes{},
			Members:     []api.GroupMembersRes{},
		})
	}
	totalPages := (int(totalData) + limit - 1) / limit

	return api.GetMyGroups201JSONResponse{
		Message: utils.StringPtr("Get Group Success"),
		Data:    &res,
		Status:  utils.BoolPtr(true),
		Meta: &api.PaginationMeta{
			CurrentPage: utils.IntPtr(limit),
			TotalPages:  utils.IntPtr(totalPages),
			TotalItems:  utils.IntPtr(int(totalData)),
		},
	}, nil
}

func (c *GroupAPI) CreateGroup(ctx context.Context, request api.CreateGroupRequestObject) (api.CreateGroupResponseObject, error) {
	userId, _, err := utils.GetUserInfo(ctx)
	if err != nil {
		return api.CreateGroup400JSONResponse{
			Errors:  utils.StringPtr("Err Message : " + err.Error()),
			Message: utils.StringPtr("Failed Parsing User Info | User Info Not Found"),
			Status:  utils.BoolPtr(false),
		}, err
	}
	createdGroup, err := c.Service.CreateGroup(ctx, userId, *request.Body)
	if err != nil {
		return api.CreateGroup500JSONResponse{
			Errors:  utils.StringPtr("Err Message : " + err.Error()),
			Message: utils.StringPtr("Failed to create Group"),
			Status:  utils.BoolPtr(false),
		}, nil // MUTLAK: Return nil di sini, bukan err!
	}

	var res api.GroupBaseRes

	// (Asumsi: Jika di YAML lu field ini tidak ada di daftar 'required',
	// lu harus membungkusnya pakai utils.StringPtr)
	res.Id = createdGroup.ID.String()
	res.Name = createdGroup.Name
	res.Description = utils.StringPtr(createdGroup.Description)

	// 2. MAPPING WALLET (Jembatan dari Array DB ke Single Object API)
	// Validasi panjang array untuk menghindari Panic jika DB gagal insert Wallet
	if len(createdGroup.Wallet) > 0 {
		w := createdGroup.Wallet[0] // Ambil index pertama secara mutlak
		res.Wallet = api.WalletRes{
			Id:               w.ID.String(),
			Name:             w.Name,
			Balance:          w.Balance,
			TransactionCount: w.TransactionCount,
			GroupId:          utils.UUIDPtrToStringPtr(w.GroupID),
			Transactions:     []api.TransactionRes{}, // Dompet baru pastinya nol transaksi
		}
	}

	// 3. MAPPING MEMBERS
	var membersRes []api.GroupMembersRes
	for _, m := range createdGroup.Members {
		membersRes = append(membersRes, api.GroupMembersRes{
			Id:       m.ID.String(),
			UserId:   m.UserID.String(),
			Username: m.User.Username,
			Role:     m.MembersRole.String(),
		})
	}
	res.Members = membersRes

	// 4. KEMBALIKAN KE FRONTEND
	// Tergantung oapi-codegen lu, kalau GroupBaseRes bukan wrapper utama,
	// sesuaikan dengan struct 201 lu.
	return api.CreateGroup201JSONResponse(res), nil

}

func (c *GroupAPI) DeleteGroup(ctx context.Context, request api.DeleteGroupRequestObject) (api.DeleteGroupResponseObject, error) {
	return api.DeleteGroup200JSONResponse{}, nil
}

func (c *GroupAPI) GetGroupById(ctx context.Context, request api.GetGroupByIdRequestObject) (api.GetGroupByIdResponseObject, error) {
	return api.GetGroupById200JSONResponse{}, nil
}

func (c *GroupAPI) UpdateGroup(ctx context.Context, request api.UpdateGroupRequestObject) (api.UpdateGroupResponseObject, error) {
	return api.UpdateGroup201JSONResponse{}, nil
}
