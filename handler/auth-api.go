package handler

import (
	"cashflow_gin/api"
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
	reqInput := api.RegisterReq{
		Username: req.Body.Username,
		Email:    req.Body.Email,
		Password: req.Body.Password,
		Nickname: req.Body.Nickname,
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

func (a *AuthAPI) Login(ctx context.Context, req api.LoginRequestObject) (api.LoginResponseObject, error) {
	reqInput := api.LoginReq{
		Email:    req.Body.Email,
		Password: req.Body.Password,
	}
	fmt.Println(req.Body.Email)
	fmt.Println(req.Body.Password)

	token, err := a.Service.Login(ctx, reqInput)
	if err != nil {

		return api.Login500JSONResponse{
			Status:  utils.BoolPtr(false),
			Message: utils.StringPtr("Can't Login"),
			Errors:  utils.StringPtr("ERR: " + err.Error()),
		}, nil
	}
	fmt.Println("token di handler : ", token)

	// RAKIT COOKIE STRING SECARA MANUAL
	// Format mutlak: Nama=Value; Max-Age=...; Path=/; Domain=...; HttpOnly
	cookieString := fmt.Sprintf("token=%s; Max-Age=86400; Path=/; Domain=localhost; HttpOnly; SameSite=Lax", token)
	// Catatan: Saat lu deploy ke server1.bagasrr.my.id, ganti Domain=localhost jadi Domain=bagasrr.my.id dan tambahkan tulisan ; Secure di ujungnya.

	return api.Login200JSONResponse{
		// GAK ADA LAGI TOKEN DI BODY JSON
		Body: api.SuccessBaseRes{
			Status:  utils.BoolPtr(true),
			Message: utils.StringPtr("Login successfully"),
		},
		Headers: api.Login200ResponseHeaders{
			SetCookie: cookieString, // Lempar cookie lewat pintu yang benar
		},
	}, nil
}
func (a *AuthAPI) ForgotPassword(ctx context.Context, req api.ForgotPasswordRequestObject) (api.ForgotPasswordResponseObject, error) {
	err := a.Service.ForgotPassword(ctx, req.Body.Email, req.Body.Password)
	if err != nil {
		return api.ForgotPassword400JSONResponse{
			Status:  utils.BoolPtr(false),
			Message: utils.StringPtr(err.Error()),
		}, err
	}

	return api.ForgotPassword200JSONResponse{
		Status:  utils.BoolPtr(true),
		Message: utils.StringPtr("Forgot password successfully"),
	}, nil
}
