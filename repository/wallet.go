package repository

import (
	"cashflow_gin/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WalletRepository interface {
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
	err := r.db.Where("id = ?", walletID).First(&wallet).Error
	return wallet, err
}
