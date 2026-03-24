package handler

import (
	"cashflow_gin/api"
	"cashflow_gin/services"
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

	// Looping untuk mapping array ke format YAML
	var responseData []api.UserRes
	for _, user := range *users {
		id, _ := uuid.Parse(user.ID) // Hati-hati, buat variabel lokal di dalam loop
		username := user.Username
		email := user.Email
		role := user.UserRole
		responseData = append(responseData, api.UserRes{
			Id:       &id,
			Username: &username,
			Email:    &email,
			UserRole: &role,
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
	// THE FIX: Cara ekstrak UserID di Strict Server (Asumsi JWT Middleware menyisipkan data ke Context)
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

	status := true
	msg := "Success"
	return api.GetMyProfile200JSONResponse{
		Status:  &status,
		Message: &msg,
		Data: &api.UserRes{
			Id:       &userID,
			Username: &user.Username,
			Email:    &user.Email,
			UserRole: &user.UserRole,
		},
	}, nil
}

func (u *UserAPI) FindUserById(ctx context.Context, request api.FindUserByIdRequestObject) (api.FindUserByIdResponseObject, error) {
	return api.FindUserById200JSONResponse{}, nil
}

func (u *UserAPI) UpdateUser(ctx context.Context, request api.UpdateUserRequestObject) (api.UpdateUserResponseObject, error) {
	return api.UpdateUser200JSONResponse{}, nil
}
