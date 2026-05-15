package services

import (
	"cashflow_gin/dto/request"
	"cashflow_gin/models"
	"cashflow_gin/repository"
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, input *request.LoginRequest) (string, error)
	Register(ctx context.Context, input request.CreateUserRequest) (*models.User, error)
}

type authService struct {
	repo repository.AuthRepository
}

func NewAuthService(r repository.AuthRepository) AuthService {
	return &authService{repo: r}
}

func (s *authService) Login(ctx context.Context, input *request.LoginRequest) (string, error) {
	// 1. Cari user berdasarkan email (panggil Repo)
	user, err := s.repo.Login(ctx, input)
	if err != nil {
		return "", errors.New("email atau password salah") // Jangan kasih tau email gak ada (security)
	}

	// 2. Bandingkan Password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return "", errors.New("email atau password salah")
	}

	// 3. Generate JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   user.ID.String(),
		"user_role": user.UserRole,
		"exp":       time.Now().Add(time.Hour * 24).Unix(), // Token berlaku 24 jam
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return tokenString, err
}

func (s *authService) Register(ctx context.Context, input request.CreateUserRequest) (*models.User, error) {
	// 1. VALIDASI DATABASE YANG SOLID
	_, err := s.repo.FindByEmail(ctx, input.Email)
	if err == nil {
		// User benar-benar ditemukan
		return nil, errors.New("email atau username sudah terdaftar")
	}
	// CEK: Apakah errornya karena DB mati, atau murni karena user belum ada?
	// Asumsi lu pake GORM, pastikan import "gorm.io/gorm"
	if err != nil && err.Error() != "record not found" {
		// Kalau errornya BUKAN "record not found", berarti DB lu lagi bermasalah
		return nil, fmt.Errorf("database error: %v", err)
	}

	// 2. HASHING PASSWORD
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("error hashing password")
	}

	defaultRole := models.RoleUser

	// 3. RAKIT USER SEKALIGUS DOMPETNYA (GORM ASSOCIATION)
	// Ingat, field relasi di model User lu namanya Wallets (bentuknya Slice/Array)
	user := models.User{
		Username:         input.Username,
		Email:            input.Email,
		Password:         string(hashedPassword),
		UserRole:         defaultRole,
		SubscriptionPlan: "free",

		// Kita langsung tempelkan dompet pertamanya di sini
		Wallets: []models.Wallet{
			{
				Name:     fmt.Sprintf("First Wallet %s", input.Username), // Typo fixed
				Balance:  0,
				Currency: "IDR",
			},
		},
	}

	// 4. LEMPAR KE REPO (Satu pemanggilan saja)
	// Ingat, kita butuh mengirim alamat memori dari struct value 'user' di atas
	createdUser, err := s.repo.CreateUserWithWallet(ctx, &user)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}
