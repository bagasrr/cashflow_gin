package repository

import (
	"cashflow_gin/dto/request"
	"cashflow_gin/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface{
	
	FindByEmailOrUsername(email, username string) (*models.User, error)
	FindAllUser() ([]models.User, error)
	FindMyProfile(id uuid.UUID) (*models.User, error)
	Login(input *request.LoginRequest) (*models.User, error)
}

type userRepository struct{
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
	err := r.db.Preload("Wallets", func(db *gorm.DB) *gorm.DB{
		return db.Limit(10).Order("updated_at DESC")
	}).Preload("Wallets.Transaction", func(db *gorm.DB) *gorm.DB{
		return db.Limit(10).Order("updated_at DESC")
	}).Preload("Wallets.Transaction.Category").Find(&users).Error
	return users, err
}

func (r *userRepository) FindMyProfile(id uuid.UUID) (*models.User, error) {
	var user *models.User
	err := r.db.First(&user, "id = ?", id).Error
	return user, err
}

func (r *userRepository) Login(input *request.LoginRequest) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "email = ?", input.Email).Error
	
	return &user, err
}