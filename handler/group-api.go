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
	page := request.Params.Page
	limit := request.Params.Limit
	groups, err := c.Service.GetAllGroups(ctx, limit, page)
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

	// 5. KEMBALIKAN RESPONSE 200 (Bukan 201, GET itu 200 Success)
	return api.GetGroups201JSONResponse{
		Data:    &res,
		Status:  utils.BoolPtr(true),
		Message: utils.StringPtr("Get Group Success"),
	}, nil
}

func (c *GroupAPI) GetMyGroups(ctx context.Context, request api.GetMyGroupsRequestObject) (api.GetMyGroupsResponseObject, error) {

	return api.GetMyGroups201JSONResponse{}, nil
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
