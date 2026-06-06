package repository

import (
	"cashflow_gin/dto/request"
	"cashflow_gin/models"
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmailOrUsername(ctx context.Context, email, username string) (*models.User, error)
	FindAllUser(ctx context.Context) (*[]models.User, error)
	FindMyProfile(ctx context.Context, id uuid.UUID) (*models.User, error)
	Login(ctx context.Context, input *request.LoginRequest) (*models.User, error)

	IsRoleAdmin(ctx context.Context, userID uuid.UUID) (bool, error)
	FindUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	UpdateUser(ctx context.Context, user models.User) (*models.User, error)
	UpdateUserByAdmin(ctx context.Context, user models.User) (*models.User, error)
	ChangeRole(ctx context.Context, user models.User) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByEmailOrUsername(ctx context.Context, email, username string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ? OR username = ?", email, username).First(&user).Error
	return &user, err
}

func (r *userRepository) FindAllUser(ctx context.Context) (*[]models.User, error) {
	var users []models.User
	walletSelectQuery := `
        wallets.*, 
        (
            SELECT COUNT(*) 
            FROM transactions 
            WHERE transactions.wallet_id = wallets.id
        ) as transaction_count
    `
	err := r.db.WithContext(ctx).Preload("Wallets", func(db *gorm.DB) *gorm.DB {
		return db.Select(walletSelectQuery)
	}).Find(&users).Error
	return &users, err
}

func (r *userRepository) FindMyProfile(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	walletSelectQuery := `
        wallets.*, 
        (
            SELECT COUNT(*) 
            FROM transactions 
            WHERE transactions.wallet_id = wallets.id
        ) as transaction_count
    `
	err := r.db.WithContext(ctx).
		// 1. Preload dengan Custom Query
		Preload("Wallets", func(db *gorm.DB) *gorm.DB {
			return db.Select(walletSelectQuery)
		}).
		// 2. Ambil User-nya (Query Utama bersih aja)
		First(&user, "id = ?", id).Error

	fmt.Println("result di repo : ", user)
	fmt.Println("error di repo : ", err)
	return &user, err
}

func (r *userRepository) Login(ctx context.Context, input *request.LoginRequest) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, "email = ?", input.Email).Error

	return &user, err
}

func (r *userRepository) IsRoleAdmin(ctx context.Context, userID uuid.UUID) (bool, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, "id = ?", userID).Error
	if err != nil {
		return false, err
	}
	return user.UserRole == models.RoleAdmin, nil
}

func (r *userRepository) FindUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	var user models.User
	walletSelectQuery := `
        wallets.*, 
        (
            SELECT COUNT(*) 
            FROM transactions 
            WHERE transactions.wallet_id = wallets.id
        ) as transaction_count
    `
	err := r.db.WithContext(ctx).
		Preload("Wallets", func(db *gorm.DB) *gorm.DB {
			return db.Select(walletSelectQuery)
		}).
		First(&user, "id = ?", userID).Error
	return &user, err
}

func (r *userRepository) UpdateUser(ctx context.Context, user models.User) (*models.User, error) {
	//err := r.db.WithContext(ctx).Save(&user).Error
	err := r.db.WithContext(ctx).Model(user).Omit("id", "password").Save(&user).Error
	return &user, err
}

func (r *userRepository) UpdateUserByAdmin(ctx context.Context, user models.User) (*models.User, error) {
	err := r.db.WithContext(ctx).Model(user).Omit("id").Save(&user).Error
	return &user, err
}

// untuk ganti role dengan cepat ke admin, maupun ke user
func (r *userRepository) ChangeRole(ctx context.Context, user models.User) (*models.User, error) {
	err := r.db.WithContext(ctx).Model(user).Omit("id").Save(&user).Error
	return &user, err
}
