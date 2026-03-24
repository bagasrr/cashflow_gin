package handler

import (
	"cashflow_gin/api"
	"cashflow_gin/dto/request"
	"cashflow_gin/services"
	"context"

	"github.com/google/uuid"
)

type AuthAPI struct {
	Service services.AuthService
}

// Register memenuhi interface dari openapi.yaml
func (a *AuthAPI) Register(ctx context.Context, req api.RegisterRequestObject) (api.RegisterResponseObject, error) {
	// 1. Mapping Eksternal DTO -> Internal DTO
	reqInput := request.CreateUserRequest{
		Username: req.Body.Username,
		Email:    req.Body.Email,
		Password: req.Body.Password,
	}

	// 2. Panggil layer Service
	user, err := a.Service.Register(ctx, reqInput)
	if err != nil {
		errStr := err.Error()
		status := false
		msg := "Internal Server Error"
		return api.Register500JSONResponse{
			Status:  &status,
			Message: &msg,
			Errors:  &errStr,
		}, nil
	}

	// 3. Mapping kembali ke Response YAML
	userID, _ := uuid.Parse(user.ID)
	status := true
	msg := "Register Success"
	return api.Register201JSONResponse{
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

// Login memenuhi interface dari openapi.yaml
func (a *AuthAPI) Login(ctx context.Context, req api.LoginRequestObject) (api.LoginResponseObject, error) {
	reqInput := request.LoginRequest{
		Email:    req.Body.Email,
		Password: req.Body.Password,
	}

	token, err := a.Service.Login(ctx, &reqInput)
	if err != nil {
		errStr := err.Error()
		status := false
		msg := "Login Gagal"
		return api.Login500JSONResponse{
			Status:  &status,
			Message: &msg,
			Errors:  &errStr,
		}, nil
	}

	status := true
	msg := "Login Success"
	return api.Login200JSONResponse{
		Status:  &status,
		Message: &msg,
		Data: &struct {
			Token *string `json:"token,omitempty"`
		}{Token: &token},
	}, nil
}
