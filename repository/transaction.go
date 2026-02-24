package repository

import (
	"cashflow_gin/models"
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreateWithWalletUpdate(ctx context.Context, transaction *models.Transaction) error
	FindAll(ctx context.Context) ([]models.Transaction, error)
	FindAllByUserID(ctx context.Context, userID uuid.UUID) ([]models.Transaction, error)
	IsOwner(ctx context.Context, userID uuid.UUID, walletID string) bool
	FindByID(ctx context.Context, transactionID uuid.UUID) (*models.Transaction, error)
	UpdateTransaction(ctx context.Context, transaction *models.Transaction) error
	UpdateTransactionWithWalletBallance(ctx context.Context, transaction *models.Transaction, delta float64) error
	SoftDeleteTransaction(ctx context.Context, transactionID uuid.UUID, delta float64, walletID uuid.UUID) error
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

// INI LOGIC PENTING: Transaction Database (ACID)
func (r *transactionRepository) CreateWithWalletUpdate(ctx context.Context, transaction *models.Transaction) error {
	// Mulai DB Transaction
	return r.db.WithContext(ctx).Exec("SELECT pg_sleep(5)").Transaction(func(tx *gorm.DB) error {
		// 1. Create Transaction Record
		if err := tx.Create(transaction).Error; err != nil {
			return err // Rollback otomatis kalau error
		}

		// 2. Update Wallet Balance
		// Logic matematika (tambah/kurang) sudah ditentukan di Service lewat field Amount
		// Kita pakai gorm.Expr biar aman dari race condition
		if err := tx.Model(&models.Wallet{}).
			Where("id = ?", transaction.WalletID).
			Updates(map[string]interface{}{
				"balance":    gorm.Expr("balance + ?", transaction.Amount), // Aman
				"updated_at": time.Now(),
			}).Error; err != nil {
			return err // Rollback otomatis
		}

		return nil // Commit
	})
}

func (r *transactionRepository) FindAll(ctx context.Context) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.WithContext(ctx).Preload("Category").Preload("Wallet").Find(&transactions).Error
	return transactions, err
}

func (r *transactionRepository) FindAllByUserID(ctx context.Context, userID uuid.UUID) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.WithContext(ctx).Preload("Category").Preload("Wallet").Where("user_id = ?", userID).Find(&transactions).Error
	return transactions, err
}

func (r *transactionRepository) IsOwner(ctx context.Context, userID uuid.UUID, walletID string) bool {
	var wallet models.Wallet
	err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", walletID, userID).First(&wallet).Error
	return err == nil
}

func (r *transactionRepository) FindByID(ctx context.Context, transactionID uuid.UUID) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.WithContext(ctx).Preload("Category").Preload("Wallet").First(&transaction, "id = ?", transactionID).Error
	return &transaction, err
}

func (r *transactionRepository) UpdateTransaction(ctx context.Context, transaction *models.Transaction) error {
	return r.db.WithContext(ctx).Save(transaction).Error
}

func (r *transactionRepository) UpdateTransactionWithWalletBallance(ctx context.Context, transaction *models.Transaction, delta float64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Update Transaction Record
		if err := tx.Save(transaction).Error; err != nil {
			return err // Rollback otomatis kalau error
		}

		// 2. Update Wallet Balance
		// Logic matematika (tambah/kurang) sudah ditentukan di Service lewat field Amount
		// Kita pakai gorm.Expr biar aman dari race condition
		if err := tx.Model(&models.Wallet{}).
			Where("id = ?", transaction.WalletID).
			Update("balance", gorm.Expr("balance + ?", delta)).Error; err != nil {
			return err // Rollback otomatis
		}

		return nil // Commit
	})
}

func (r *transactionRepository) SoftDeleteTransaction(ctx context.Context, transactionId uuid.UUID, delta float64, walletID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Soft Delete Transaction Record
		if err := tx.Where("id = ?", transactionId).Delete(&models.Transaction{}).Error; err != nil {
			return err // Rollback otomatis kalau error
		}

		// 2. Update Wallet Balance (kembalikan ke kondisi sebelum transaksi)
		if err := tx.Model(&models.Wallet{}).
			Where("id = ?", walletID). // Asumsi kita punya walletID di transactionID, bisa juga lewat join
			Update("balance", gorm.Expr("balance + ?", delta)).Error; err != nil {
			return err // Rollback otomatis
		}

		return nil // Commit
	})
}
