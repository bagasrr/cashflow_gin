package services

import (
	"cashflow_gin/dto/request"
	"cashflow_gin/dto/response"
	"cashflow_gin/models"
	"cashflow_gin/repository"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(input *request.LoginRequest) (string, error)
	Register(input request.CreateUserRequest) (*response.UserResponse, error)
}

type authService struct {
	repo repository.AuthRepository
}

func NewAuthService(r repository.AuthRepository) AuthService {
	return &authService{repo: r}
}

func (s *authService) Login(input *request.LoginRequest) (string, error) {
	// 1. Cari user berdasarkan email (panggil Repo)
	user, err := s.repo.Login(input)
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

func (s *authService) Register(input request.CreateUserRequest) (*response.UserResponse, error) {
	// cek apakah email atau username udah ada di db?
	_, err := s.repo.FindByEmail(input.Email)
	// kalo udah ada kan err = nil, kembalikan error
	if err == nil {
		return nil, errors.New("email atau username sudah terdaftar")
	}

	// kalo err != nil / belum ada
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("error hashing password")
	}

	defaultRole := models.RoleUser

	user := models.User{
		Username:         input.Username,
		Email:            input.Email,
		Password:         string(hashedPassword),
		UserRole:         defaultRole, // assign a pointer to models.RoleUser
		SubscriptionPlan: "free",
	}

	wallet := &models.Wallet{
		UserID:   &user.ID,
		Name:     fmt.Sprintf("Frist Wallet %s", user.Username),
		Balance:  0,
		Currency: "IDR",
	}

	err = s.repo.CreateUserWithWallet(&user, wallet)
	if err != nil {
		return nil, err
	}

	res := &response.UserResponse{
		ID:       user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
		UserRole: user.UserRole.String(),
		Wallets: []response.WalletResponse{
			{
				ID:      wallet.ID,
				Name:    wallet.Name,
				Balance: wallet.Balance,
			},
		},
	}

	return res, nil
}
