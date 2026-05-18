package handler

import (
	"cashflow_gin/api"
	"cashflow_gin/dto/request"
	"cashflow_gin/services"
	"cashflow_gin/utils"
	"context"
	"fmt"
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

	fmt.Println("user di handler : ", user)

	// 3. Mapping kembali ke Response YAML
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

	return api.Register201JSONResponse{
		Status:  utils.BoolPtr(true),
		Message: utils.StringPtr("Create user successfully"),
		Data: &api.UserRes{
			Id:       user.ID.String(),
			Username: user.Username,
			Email:    user.Email,
			UserRole: user.UserRole.String(),
			Wallets:  &wallets,
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

func (a *AuthAPI) ForgotPassword(ctx context.Context, req api.ForgotPasswordRequestObject) (api.ForgotPasswordResponseObject, error) {
	err := a.Service.ForgotPassword(ctx, req.Body.Email, req.Body.Password)
	if err != nil {
		errStr := err.Error()
		status := false
		msg := "Forgot Password Error"
		return api.ForgotPassword400JSONResponse{
			Status:  &status,
			Message: &msg,
			Errors:  &errStr,
		}
	}
	var res api.UserRes
	res.Id = user.ID.String()
	res.Username = user.Username
	res.Email = user.Email
	res.UserRole = user.UserRole.String()

	return res, nil
}
