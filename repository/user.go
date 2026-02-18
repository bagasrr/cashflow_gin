package repository

import (
	"cashflow_gin/dto/request"
	"cashflow_gin/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmailOrUsername(email, username string) (*models.User, error)
	FindAllUser() ([]models.User, error)
	FindMyProfile(id uuid.UUID) (*models.User, error)
	Login(input *request.LoginRequest) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByEmailOrUsername(email, username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ? OR username = ?", email, username).First(&user).Error
	return &user, err
}

func (r *userRepository) FindAllUser() ([]models.User, error) {
	var users []models.User
	subQuery := `(
        SELECT COUNT(*) 
        FROM transactions t
        JOIN wallets w ON t.wallet_id = w.id
        WHERE w.user_id = users.id
    )`
	err := r.db.Select("users.*, " + subQuery + " as transaction_count").Preload("Wallets").Find(&users).Error
	return users, err
}

func (r *userRepository) FindMyProfile(id uuid.UUID) (*models.User, error) {
	var user *models.User
	walletSelectQuery := `
        wallets.*, 
        (
            SELECT COUNT(*) 
            FROM transactions 
            WHERE transactions.wallet_id = wallets.id
        ) as transaction_count
    `
	err := r.db.
		// 1. Preload dengan Custom Query
		Preload("Wallets", func(db *gorm.DB) *gorm.DB {
			return db.Select(walletSelectQuery)
		}).
		// 2. Ambil User-nya (Query Utama bersih aja)
		First(&user, "id = ?", id).Error
	return user, err
}

func (r *userRepository) Login(input *request.LoginRequest) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "email = ?", input.Email).Error

	return &user, err
}
