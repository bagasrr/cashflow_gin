package repository

import (
	"cashflow_gin/dto/request"
	"cashflow_gin/models"
	"context"

	"gorm.io/gorm"
)

type AuthRepository interface {
	Login(ctx context.Context, input *request.LoginRequest) (*models.User, error)
	Register(ctx context.Context, input *request.CreateUserRequest) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	//CreateUserWithWallet(ctx context.Context, user *models.User, wallet *models.Wallet) (*models.User, error)
	CreateUserWithWallet(ctx context.Context, user *models.User) (*models.User, error)
	FindUserForPasswordReset(ctx context.Context, email string) (*models.User, error)
	UpdatePassword(ctx context.Context, user *models.User) error
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) Login(ctx context.Context, input *request.LoginRequest) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, "email = ? ", input.Email).Error

	return &user, err
}

func (r *authRepository) Register(ctx context.Context, input *request.CreateUserRequest) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Create(&user).Error

	return &user, err
}

func (r *authRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error
	return &user, err
}

// Cuma butuh nerima *models.User, karena walletnya nanti diselipin di dalemnya
func (r *authRepository) CreateUserWithWallet(ctx context.Context, user *models.User) (*models.User, error) {
	// 1. HAPUS TANDA & DI DEPAN user, KARENA user SUDAH POINTER
	// 2. GORM otomatis membuka transaksi dan menginsert relasi (Wallets) di dalamnya
	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *authRepository) FindUserForPasswordReset(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	// GORM akan mencari user yang memiliki email DAN username yang persis sama
	err := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&user).Error

	if err != nil {
		// Biarkan error ini mengalir ke Service.
		// Service nanti yang ngecek apakah errornya gorm.ErrRecordNotFound
		return nil, err
	}

	return &user, nil
}

func (r *authRepository) UpdatePassword(ctx context.Context, user *models.User) error {
	err := r.db.WithContext(ctx).
		Model(user).Omit("email", "username", "user_role", "wallets").
		Updates(user).Error
	return err
}
