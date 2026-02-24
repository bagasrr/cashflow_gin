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
	CreateUserWithWallet(ctx context.Context, user *models.User, wallet *models.Wallet) error
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

func (r *authRepository) CreateUserWithWallet(ctx context.Context, user *models.User, wallet *models.Wallet) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		wallet.UserID = &user.ID
		if err := tx.Create(wallet).Error; err != nil {
			return err
		}
		return nil
	})
}
