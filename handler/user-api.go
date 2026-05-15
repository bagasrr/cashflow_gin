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

type UserAPI struct {
	Service services.UserService
}

func (u *UserAPI) FindAllUsers(ctx context.Context, req api.FindAllUsersRequestObject) (api.FindAllUsersResponseObject, error) {
	users, err := u.Service.FindAllUser(ctx)
	if err != nil {
		errStr := err.Error()
		status := false
		msg := "Error retrieving users"
		return api.FindAllUsers500JSONResponse{
			Status:  &status,
			Message: &msg,
			Errors:  &errStr,
		}, nil
	}

	var responseData []api.UserRes
	for _, user := range *users {
		responseData = append(responseData, api.UserRes{
			Id:       user.ID.String(), // UBAH MUTLAK: Harus pakai .String()
			Username: user.Username,
			Email:    user.Email,
			UserRole: user.UserRole.String(), // UBAH MUTLAK: Harus pakai .String()
		})
	}

	status := true
	msg := "Success"
	return api.FindAllUsers200JSONResponse{
		Status:  &status,
		Message: &msg,
		Data:    &responseData,
	}, nil
}

func (u *UserAPI) GetMyProfile(ctx context.Context, req api.GetMyProfileRequestObject) (api.GetMyProfileResponseObject, error) {
	userIDClaim := ctx.Value("user_id")
	if userIDClaim == nil {
		status := false
		msg := "Unauthorized"
		errStr := "User ID missing in context"
		return api.GetMyProfile401JSONResponse{
			Status:  &status,
			Message: &msg,
			Errors:  &errStr,
		}, nil
	}

	userID, err := uuid.Parse(fmt.Sprintf("%v", userIDClaim))
	if err != nil {
		status := false
		msg := "Unauthorized"
		errStr := "Invalid UUID format"
		return api.GetMyProfile401JSONResponse{Status: &status, Message: &msg, Errors: &errStr}, nil
	}

	user, err := u.Service.GetMyProfile(ctx, userID)
	if err != nil {
		errStr := err.Error()
		status := false
		msg := "Failed to retrieve profile"
		return api.GetMyProfile500JSONResponse{Status: &status, Message: &msg, Errors: &errStr}, nil
	}

	var wallets []api.WalletRes
	for _, wallet := range user.Wallets {
		wallets = append(wallets, api.WalletRes{
			Id:               wallet.ID.String(),
			Balance:          wallet.Balance,
			GroupId:          utils.UUIDPtrToStringPtr(wallet.GroupID),
			Name:             wallet.Name,
			TransactionCount: wallet.TransactionCount,
		})
	}

	res := api.UserRes{
		Email:    user.Email,
		Username: user.Username,
		Id:       user.ID.String(),
		UserRole: user.UserRole.String(),
		Wallets:  wallets,
	}

	return api.GetMyProfile200JSONResponse{
		Status:  utils.BoolPtr(true),
		Message: utils.StringPtr("Success get user profile"),
		Data:    &res,
	}, nil
}

func (u *UserAPI) FindUserById(ctx context.Context, request api.FindUserByIdRequestObject) (api.FindUserByIdResponseObject, error) {
	requesterRole, err := utils.GetUserRole(ctx)
	if err != nil {
		status := false
		msg := "Failed to get requester role: " + err.Error()
		return api.FindUserById500JSONResponse{
			Status:  &status,
			Message: &msg,
		}, nil
	}

	targetID, err := uuid.Parse(request.Id)
	if err != nil {
		status := false
		msg := "Failed to parse user ID: " + err.Error()
		return api.FindUserById500JSONResponse{
			Status:  &status,
			Message: &msg,
		}, nil
	}

	target, err := u.Service.FindUserByID(ctx, targetID, requesterRole)
	if err != nil {
		status := false
		msg := "Failed to retrieve user: " + err.Error()
		return api.FindUserById500JSONResponse{
			Status:  &status,
			Message: &msg,
		}, nil
	}

	status := true
	msg := "Success"
	return api.FindUserById200JSONResponse{
		Status:  &status,
		Message: &msg,
		Data: &api.UserRes{
			Id:       target.ID.String(), // UBAH MUTLAK: Harus pakai .String()
			Username: target.Username,
			Email:    target.Email,
			UserRole: target.UserRole.String(), // UBAH MUTLAK: Harus pakai .String()
		},
	}, nil
}

func (u *UserAPI) UpdateUser(ctx context.Context, request api.UpdateUserRequestObject) (api.UpdateUserResponseObject, error) {
	requestor, err := utils.GetUserID(ctx)
	if err != nil {
		return api.UpdateUser500JSONResponse{
			Status:  utils.BoolPtr(false),
			Message: utils.StringPtr(err.Error()),
		}, nil
	}

	targetID, err := uuid.Parse(request.Id)
	if err != nil {
		status := false
		msg := "Failed to parse target user ID: " + err.Error()
		return api.UpdateUser500JSONResponse{
			Status:  &status,
			Message: &msg,
		}, nil
	}

	targetUser := models.User{
		Base: models.Base{
			ID: targetID,
		},
		Username: request.Body.Username,
		Email:    request.Body.Email,
		Password: request.Body.Password,
	}

	res, err := u.Service.UpdateUser(ctx, requestor, targetUser)
	if err != nil {
		status := false
		msg := "Failed to update user: " + err.Error()
		return api.UpdateUser500JSONResponse{
			Status:  &status,
			Message: &msg,
		}, nil
	}

	status := true
	msg := "Success"
	return api.UpdateUser201JSONResponse{
		Status:  &status,
		Message: &msg,
		Data: &api.UserRes{
			Id:       res.ID.String(),
			Username: res.Username,
			Email:    res.Email,
			UserRole: res.UserRole.String(),
		},
	}, nil
}
