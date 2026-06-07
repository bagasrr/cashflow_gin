package handler

import (
	"cashflow_gin/api"
	"cashflow_gin/models"
	"cashflow_gin/services"
	"cashflow_gin/utils"
	"context"

	"github.com/google/uuid"
)

type GroupAPI struct {
	Service services.GroupService
}

func (c *GroupAPI) GetGroups(ctx context.Context, request api.GetGroupsRequestObject) (api.GetGroupsResponseObject, error) {
	_, role, err := utils.GetUserInfo(ctx)
	if err != nil {
		return api.GetGroups401JSONResponse{
			Errors:  utils.StringPtr("ERR MESSAGE : " + err.Error()),
			Message: utils.StringPtr("Failed to get user info : " + err.Error()),
			Status:  utils.BoolPtr(false),
		}, nil
	}
	if role != models.RoleAdmin {
		return api.GetGroups401JSONResponse{
			Message: utils.StringPtr("You do not have access to this resource"),
			Status:  utils.BoolPtr(false),
		}, nil
	}
	page, limit, offset := utils.ValidatePagination(request.Params.Page, request.Params.Limit)
	groups, totalItem, err := c.Service.GetAllGroups(ctx, limit, offset)
	if err != nil {
		return api.GetGroups500JSONResponse{
			Errors:  utils.StringPtr("ERR MESSAGE : " + err.Error()),
			Message: utils.StringPtr("Failed to get groups : " + err.Error()),
			Status:  utils.BoolPtr(false),
		}, nil
	}
	res := []api.GroupRes{}

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
				UserId:  utils.UUIDPtrToStringPtr(w.UserID),
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
		res = append(res, api.GroupRes{
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
	return api.GetGroups200JSONResponse{
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
			Errors:  utils.StringPtr("ERR MESSAGE : " + err.Error()),
			Message: utils.StringPtr("Failed to get user info : " + err.Error()),
			Status:  utils.BoolPtr(false),
		}, nil
	}

	page, limit, offset := utils.ValidatePagination(request.Params.Page, request.Params.Limit)
	myGroups, totalData, err := c.Service.GetMyGroups(ctx, page, offset, userId)

	if err != nil {
		return api.GetMyGroups500JSONResponse{
			Errors:  utils.StringPtr("ERR MESSAGE : " + err.Error()),
			Message: utils.StringPtr("Failed to get groups : " + err.Error()),
			Status:  utils.BoolPtr(false),
		}, nil
	}
	var res []api.GroupRes
	for _, group := range *myGroups {
		res = append(res, api.GroupRes{
			Id:          group.ID.String(),
			Name:        group.Name,
			Description: utils.StringPtr(group.Description),
			Wallet:      []api.WalletRes{},
			Members:     []api.GroupMembersRes{},
		})
	}
	totalPages := (int(totalData) + limit - 1) / limit

	return api.GetMyGroups200JSONResponse{
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
		}, nil
	}
	createdGroup, err := c.Service.CreateGroup(ctx, userId, request.Body)
	if err != nil {
		return api.CreateGroup500JSONResponse{
			Errors:  utils.StringPtr("Err Message : " + err.Error()),
			Message: utils.StringPtr("Failed to create Group"),
			Status:  utils.BoolPtr(false),
		}, nil
	}

	var res api.GroupRes
	res.Id = createdGroup.ID.String()
	res.Name = createdGroup.Name
	res.Description = utils.StringPtr(createdGroup.Description)
	var walletRes []api.WalletRes
	for _, w := range createdGroup.Wallet {
		walletRes = append(walletRes, api.WalletRes{
			Id:               w.ID.String(),
			Name:             w.Name,
			Balance:          w.Balance,
			TransactionCount: w.TransactionCount,
			GroupId:          utils.UUIDPtrToStringPtr(w.GroupID),
			Transactions:     []api.TransactionRes{},
		})
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

	return api.CreateGroup201JSONResponse{
		Message: utils.StringPtr("Create Group Success"),
		Data:    &res,
		Status:  utils.BoolPtr(true),
	}, nil

}

func (c *GroupAPI) DeleteGroup(ctx context.Context, request api.DeleteGroupRequestObject) (api.DeleteGroupResponseObject, error) {
	userId, _, err := utils.GetUserInfo(ctx)
	if err != nil {
		return api.DeleteGroup400JSONResponse{
			Errors:  utils.StringPtr("Err Message : " + err.Error()),
			Message: utils.StringPtr("Failed Parsing User Info | User Info Not Found"),
			Status:  utils.BoolPtr(false),
		}, nil
	}
	groupId, err := uuid.Parse(request.Id)
	if err != nil {
		return api.DeleteGroup400JSONResponse{
			Errors:  utils.StringPtr("Err Message : " + err.Error()),
			Message: utils.StringPtr("Failed to get group id : " + err.Error()),
			Status:  utils.BoolPtr(false),
		}, nil
	}
	delErr := c.Service.DeleteGroup(ctx, userId, groupId)
	if delErr != nil {
		return api.DeleteGroup500JSONResponse{
			Errors:  utils.StringPtr("ERR MESSAGE : " + delErr.Error()),
			Message: utils.StringPtr("Failed to delete Group : " + delErr.Error()),
			Status:  utils.BoolPtr(false),
		}, nil
	}
	return api.DeleteGroup200JSONResponse{
		Message: utils.StringPtr("Delete Group Success"),
		Status:  utils.BoolPtr(true),
	}, nil
}

func (c *GroupAPI) GetGroupById(ctx context.Context, request api.GetGroupByIdRequestObject) (api.GetGroupByIdResponseObject, error) {
	groupId, err := uuid.Parse(request.Id)
	if err != nil {
		return api.GetGroupById400JSONResponse{
			Errors:  utils.StringPtr("Err Message : " + err.Error()),
			Message: utils.StringPtr("Failed to get group id : " + err.Error()),
			Status:  utils.BoolPtr(false),
		}, nil
	}
	group, err := c.Service.GetGroupByID(ctx, groupId)
	if err != nil {
		return api.GetGroupById500JSONResponse{
			Errors:  utils.StringPtr("Err Message : " + err.Error()),
			Message: utils.StringPtr("Failed to get group : " + err.Error()),
			Status:  utils.BoolPtr(false),
		}, nil
	}
	var res api.GroupRes
	res.Id = group.ID.String()
	res.Name = group.Name
	res.Description = utils.StringPtr(group.Description)
	var walletRes []api.WalletRes
	for _, w := range group.Wallet {
		walletRes = append(walletRes, api.WalletRes{
			Id:               w.ID.String(),
			Name:             w.Name,
			Balance:          w.Balance,
			TransactionCount: w.TransactionCount,
			GroupId:          utils.UUIDPtrToStringPtr(w.GroupID),
			Transactions:     []api.TransactionRes{},
		})
	}
	return api.GetGroupById200JSONResponse{
		Message: utils.StringPtr("Get Group Success"),
		Status:  utils.BoolPtr(true),
		Data:    &res,
	}, nil
}

func (c *GroupAPI) UpdateGroup(ctx context.Context, request api.UpdateGroupRequestObject) (api.UpdateGroupResponseObject, error) {
	// 1. Tarik User ID dari JWT Token
	userId, _, err := utils.GetUserInfo(ctx)
	if err != nil {
		return api.UpdateGroup401JSONResponse{
			Message: utils.StringPtr("Err Message : " + err.Error()),
			Status:  utils.BoolPtr(false),
		}, nil
	}

	// 2. Validasi Group ID
	groupId, err := uuid.Parse(request.Id)
	if err != nil {
		return api.UpdateGroup400JSONResponse{
			Message: utils.StringPtr("Err Message : " + err.Error()),
			Status:  utils.BoolPtr(false),
		}, nil
	}

	// 3. Lempar TIGA parameter ke Service (ctx, userID, groupID)
	updatedGroup, err := c.Service.UpdateGroup(ctx, userId, groupId, request.Body)
	if err != nil {
		return api.UpdateGroup500JSONResponse{
			Errors:  utils.StringPtr("Err Message : " + err.Error()),
			Message: utils.StringPtr("Failed to update Group"),
			Status:  utils.BoolPtr(false),
		}, nil
	}

	// 4. MAPPING RESPONSE
	var res api.GroupRes
	res.Id = updatedGroup.ID.String()
	res.Name = updatedGroup.Name
	res.Description = utils.StringPtr(updatedGroup.Description)

	// Mapping Wallet (Pakai pola multi-wallet yang aman)
	var walletRes []api.WalletRes
	for _, w := range updatedGroup.Wallet {
		walletRes = append(walletRes, api.WalletRes{
			Id:               w.ID.String(),
			Name:             w.Name,
			Balance:          w.Balance,
			TransactionCount: w.TransactionCount,
			GroupId:          utils.UUIDPtrToStringPtr(w.GroupID),
			Transactions:     []api.TransactionRes{},
		})
	}
	if walletRes == nil {
		walletRes = []api.WalletRes{}
	}
	res.Wallet = walletRes

	// Jangan lupa mapping Members juga kalau YAML lu minta (res.Members = ...)

	// KEMBALIKAN STATUS 200 OK (Ganti di YAML lu kalau masih 201)
	return api.UpdateGroup200JSONResponse{
		Message: utils.StringPtr("Update Group Success"),
		Status:  utils.BoolPtr(true),
		Data:    &res,
	}, nil
}
