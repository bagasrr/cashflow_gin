package repository

import (
	"cashflow_gin/models"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WalletRepository interface {
	CreateWallet(ctx context.Context, wallet models.Wallet) (*models.Wallet, error)
	FindAll(ctx context.Context, offset int, limit int) (*[]models.Wallet, error)
	FindByID(ctx context.Context, walletID uuid.UUID) (*models.Wallet, error)
	FindAllMine(ctx context.Context, userID uuid.UUID, limit, offset int) (*[]models.Wallet, error)
}

type walletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) WalletRepository {
	return &walletRepository{db: db}
}

const countQuery = "(SELECT COUNT(*) FROM transactions WHERE transactions.wallet_id = wallets.id AND transactions.deleted_at IS NULL)"

// Ubah return type menjadi pointer (*models.Wallet)
func (r *walletRepository) FindByID(ctx context.Context, walletID uuid.UUID) (*models.Wallet, error) {
	var wallet models.Wallet

	err := r.db.WithContext(ctx).
		Select("wallets.*, "+countQuery+" as transaction_count").
		// Cukup tulis anak-anaknya, GORM otomatis ambil bapaknya (Transactions)
		Preload("Transactions.Category").
		Preload("Transactions.User").
		// Pindahkan kondisi Where langsung ke dalam First biar lebih rapi
		First(&wallet, "id = ?", walletID).Error
	if err != nil {
		return nil, err // Return nil kalau error/tidak ketemu, bukan struct kosong
	}

	return &wallet, nil // Return alamat memorinya
}

func (r *walletRepository) FindAll(ctx context.Context, offset int, limit int) (*[]models.Wallet, error) {
	var wallets []models.Wallet

	// Susun logic-nya: Panggil DB -> Set Limit -> Eksekusi (Find)
	err := r.db.WithContext(ctx).
		Select("wallets.*, " + countQuery + " as transaction_count").
		Limit(limit).
		Offset(offset).
		Find(&wallets).Error

	return &wallets, err
}

// repository/wallet.go
func (r *walletRepository) FindAllMine(ctx context.Context, userID uuid.UUID, limit int, offset int) (*[]models.Wallet, error) {
	var wallets []models.Wallet

	err := r.db.WithContext(ctx).
		Select("wallets.*, "+countQuery+" as transaction_count").
		Where("user_id = ?", userID).
		Limit(limit).
		Offset(offset).
		Find(&wallets).Error

	return &wallets, err
}

func (r *walletRepository) CreateWallet(ctx context.Context, wallet models.Wallet) (*models.Wallet, error) {
	err := r.db.WithContext(ctx).Create(&wallet).Error
	return &wallet, err
}
