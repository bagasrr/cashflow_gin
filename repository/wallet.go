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
	FindAllMine(ctx context.Context, userID uuid.UUID, limit, offset int) (*[]models.Wallet, int64, error)
	SoftDeleteWallet(ctx context.Context, wallet uuid.UUID) error
	GetWalletByID(ctx context.Context, walletID uuid.UUID) (*models.Wallet, error)
	UpdateWallet(ctx context.Context, wallet *models.Wallet) (*models.Wallet, error)
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
func (r *walletRepository) FindAllMine(ctx context.Context, userID uuid.UUID, limit int, offset int) (*[]models.Wallet, int64, error) {
	var wallets []models.Wallet
	var totalItems int64
	query := r.db.WithContext(ctx).
		Select("wallets.*, "+countQuery+" as transaction_count").
		Where("user_id = ?", userID)

	if err := query.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Limit(limit).Offset(offset).
		Find(&wallets).Error; err != nil {
		return nil, 0, err
	}

	return &wallets, totalItems, nil
}

func (r *walletRepository) CreateWallet(ctx context.Context, wallet models.Wallet) (*models.Wallet, error) {
	err := r.db.WithContext(ctx).Create(&wallet).Error
	return &wallet, err
}
func (r *walletRepository) SoftDeleteWallet(ctx context.Context, walletId uuid.UUID) error {
	err := r.db.WithContext(ctx).Delete(&models.Wallet{}, "id = ?", walletId).Error
	return err
}

// 1. Tarik data lama
func (r *walletRepository) GetWalletByID(ctx context.Context, walletID uuid.UUID) (*models.Wallet, error) {
	var wallet models.Wallet
	// Gak perlu Preload transaksi, kita cuma mau ngedit nama
	err := r.db.WithContext(ctx).First(&wallet, "id = ?", walletID).Error
	return &wallet, err
}

// 2. Timpa data
func (r *walletRepository) UpdateWallet(ctx context.Context, wallet *models.Wallet) (*models.Wallet, error) {
	err := r.db.WithContext(ctx).Save(wallet).Error
	return wallet, err
}
