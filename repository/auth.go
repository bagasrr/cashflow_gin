package repository

import (
	"cashflow_gin/dto/request"
	"cashflow_gin/models"

	"gorm.io/gorm"
)

type AuthRepository interface {
	Login(input *request.LoginRequest) (*models.User, error)
	Register(input *request.CreateUserRequest) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	CreateUserWithWallet(user *models.User, wallet *models.Wallet) error
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) Login(input *request.LoginRequest) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "email = ? ", input.Email).Error

	return &user, err
}

func (r *authRepository) Register(input *request.CreateUserRequest) (*models.User, error) {
	var user models.User
	err := r.db.Create(&user).Error

	return &user, err
}

func (r *authRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "email = ?", email).Error
	return &user, err
}

func (r *authRepository) CreateUserWithWallet(user *models.User, wallet *models.Wallet) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
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
