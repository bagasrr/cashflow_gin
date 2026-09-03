package handler

import (
	"cashflow_gin/api"
	"cashflow_gin/services"
	"cashflow_gin/utils"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type AuthAPI struct {
	Service     services.AuthService
	RedisClient *redis.Client
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
			Message: utils.StringPtr("Email or Password is wrong"),
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
			SetCookie: utils.StringPtr(cookieString), // Lempar cookie lewat pintu yang benar
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

func (c *AuthAPI) Logout(ctx context.Context, request api.LogoutRequestObject) (api.LogoutResponseObject, error) {
	// 1. Tarik data yang sudah disiapkan oleh Middleware
	// Karena lu pakai oapi-codegen (context standar), lu harus menggunakan utilitas buatan lu
	// untuk mengekstrak value dari Gin Context yang terbungkus.
	tokenString, err := utils.GetStringFromContext(ctx, "token_string")
	if err != nil || tokenString == "" {
		return api.Logout500JSONResponse{
			Status:  utils.BoolPtr(false),
			Message: utils.StringPtr("Internal Server Error"),
			Errors:  utils.StringPtr("Gagal mengekstrak token dari context"),
		}, nil
	}

	expFloat, err := utils.GetFloatFromContext(ctx, "token_exp")
	if err != nil {
		return api.Logout500JSONResponse{
			Status:  utils.BoolPtr(false),
			Message: utils.StringPtr("Internal Server Error"),
			Errors:  utils.StringPtr("Gagal mengekstrak waktu kedaluwarsa"),
		}, nil
	}

	// 2. MATEMATIKA TTL MUTLAK
	// Konversi float64 (standar Unix timestamp di JWT) ke bentuk Time Golang
	expiresAt := time.Unix(int64(expFloat), 0)

	// Hitung sisa umur token: Waktu Kedaluwarsa - Waktu Saat Ini
	ttl := time.Until(expiresAt)

	// 3. EKSEKUSI BLACKLIST KE REDIS
	// Jika TTL masih lebih dari 0 detik (token belum mati secara alami), masukkan ke daftar hitam
	if ttl > 0 {
		// Gunakan c.RedisClient (pastikan lu udah inject Redis ke struct AuthAPI lu)
		redisErr := c.RedisClient.Set(ctx, tokenString, "revoked", ttl).Err()
		if redisErr != nil {
			return api.Logout500JSONResponse{
				Status:  utils.BoolPtr(false),
				Message: utils.StringPtr("Logout Gagal"),
				Errors:  utils.StringPtr("Gagal memblokir token di server"),
			}, nil
		}
	}

	// 4. CLIENT-SIDE DESTRUCTION (Jika lu menggunakan Cookie)
	// Kalau Front-end lu (Next.js/Flutter) murni menggunakan Bearer Header dari LocalStorage,
	// lu nggak perlu ngirim instruksi hapus cookie. Front-end lu yang wajib nge-clear storage-nya sendiri.
	// TAPI, kalau API lu memakai HttpOnly Cookie, buka komentar di bawah ini dan pastikan lu punya
	// utilitas untuk mengakses Gin ResponseWriter dari context.

	destroyCookieString := "token=; Max-Age=0; Path=/; Domain=localhost; HttpOnly; SameSite=Lax"

	// 5. KEMBALIKAN RESPONS SUKSES MUTLAK LEWAT JALUR RESMI
	return api.Logout200JSONResponse{
		Body: api.SuccessBaseRes{ // Sesuaikan dengan nama struct response lu yang bener dari oapi-codegen
			Status:  utils.BoolPtr(true),
			Message: utils.StringPtr("Logout berhasil, sesi telah dihancurkan"),
		},
		Headers: api.Logout200ResponseHeaders{
			SetCookie: utils.StringPtr(destroyCookieString),
		},
	}, nil
}
