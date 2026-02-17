package repository

import (
	"cashflow_gin/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WalletRepository interface {
	FindAll() (*[]models.Wallet, error)
	FindByID(walletID uuid.UUID) (models.Wallet, error)
}

type walletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) WalletRepository {
	return &walletRepository{db: db}
}

func (r *walletRepository) FindByID(walletID uuid.UUID) (models.Wallet, error) {
	var wallet models.Wallet
	err := r.db.Where("id = ?", walletID).Preload("Transactions").Preload("Transactions.Category").Preload("Transactions.User").First(&wallet).Error
	return wallet, err
}

func (r *walletRepository) FindAll() (*[]models.Wallet, error) {
	var wallets []models.Wallet
	err := r.db.Find(&wallets, func(db *gorm.DB) *gorm.DB {
		return db.Limit(10)
	}).Error
	return &wallets, err
}
