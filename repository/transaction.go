package repository

import (
	"cashflow_gin/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreateWithWalletUpdate(transaction *models.Transaction) error
	FindAll() ([]models.Transaction, error)
	IsOwner(userID uuid.UUID, walletID string) bool
	FindByID(transactionID uuid.UUID) (*models.Transaction, error)
	UpdateTransaction(transaction *models.Transaction) error
	UpdateTransactionWithWalletBallance(transaction *models.Transaction, delta float64) error
	SoftDeleteTransaction(transactionID uuid.UUID, delta float64, walletID uuid.UUID) error
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

// INI LOGIC PENTING: Transaction Database (ACID)
func (r *transactionRepository) CreateWithWalletUpdate(transaction *models.Transaction) error {
	// Mulai DB Transaction
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Create Transaction Record
		if err := tx.Create(transaction).Error; err != nil {
			return err // Rollback otomatis kalau error
		}

		// 2. Update Wallet Balance
		// Logic matematika (tambah/kurang) sudah ditentukan di Service lewat field Amount
		// Kita pakai gorm.Expr biar aman dari race condition
		if err := tx.Model(&models.Wallet{}).
			Where("id = ?", transaction.WalletID).
			Update("balance", gorm.Expr("balance + ?", transaction.Amount)).Error; err != nil {
			return err // Rollback otomatis
		}

		return nil // Commit
	})
}

func (r *transactionRepository) FindAll() ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Preload("Category").Preload("Wallet").Find(&transactions).Error
	return transactions, err
}

func (r *transactionRepository) IsOwner(userID uuid.UUID, walletID string) bool {
	var wallet models.Wallet
	err := r.db.Where("id = ? AND user_id = ?", walletID, userID).First(&wallet).Error
	return err == nil
}

func (r *transactionRepository) FindByID(transactionID uuid.UUID) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Preload("Category").Preload("Wallet").First(&transaction, "id = ?", transactionID).Error
	return &transaction, err
}

func (r *transactionRepository) UpdateTransaction(transaction *models.Transaction) error {
	return r.db.Save(transaction).Error
}

func (r *transactionRepository) UpdateTransactionWithWalletBallance(transaction *models.Transaction, delta float64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
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

func (r *transactionRepository) SoftDeleteTransaction(transactionId uuid.UUID, delta float64, walletID uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
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
