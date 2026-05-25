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
	return api.CreateGroup201JSONResponse{}, nil
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
