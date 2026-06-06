package repository

import (
	"cashflow_gin/models"
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreateWithWalletUpdate(ctx context.Context, transaction *models.Transaction) (*models.Transaction, error)
	FindAll(ctx context.Context) (*[]models.Transaction, error)
	FindAllByUserID(ctx context.Context, userID uuid.UUID) (*[]models.Transaction, error)
	IsOwner(ctx context.Context, userID uuid.UUID, walletID string) bool
	FindByID(ctx context.Context, transactionID uuid.UUID) (*models.Transaction, error)
	UpdateTransaction(ctx context.Context, transaction *models.Transaction) error
	// UBAH MUTLAK: delta sekarang int64
	UpdateTransactionWithWalletBallance(ctx context.Context, transaction *models.Transaction, delta int64) error
	// UBAH MUTLAK: delta sekarang int64
	SoftDeleteTransaction(ctx context.Context, transactionID uuid.UUID, delta int64, walletID uuid.UUID) error
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

// Tanda tangan WAJIB me-return *models.Transaction agar Handler lu dapet data
func (r *transactionRepository) CreateWithWalletUpdate(ctx context.Context, transaction *models.Transaction) (*models.Transaction, error) {
	// 1. MULAI TRANSAKSI DATABASE (ACID)
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// A. Simpan Transaksi
		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		// B. Potong/Tambah Saldo Dompet dengan aman dari Race Condition
		if err := tx.Model(&models.Wallet{}).
			Where("id = ?", transaction.WalletID).
			Updates(map[string]interface{}{
				"balance":    gorm.Expr("balance + ?", transaction.Amount),
				"updated_at": time.Now(), // Hati-hati, import "time" jika belum
			}).Error; err != nil {
			return err // Jika dompet gagal diupdate, transaksi akan di-Rollback otomatis
		}

		return nil // Commit transaksi
	})

	// 2. CEK STATUS TRANSAKSI
	if err != nil {
		return nil, err
	}

	// 3. RELOAD SECARA MUTLAK (The Enterprise Way)
	// Transaksi sukses, sekarang tarik ulang datanya + relasi untuk dikirim ke Handler
	var reloadedTrx models.Transaction
	err = r.db.WithContext(ctx).
		Preload("Category").
		Preload("User").
		First(&reloadedTrx, "id = ?", transaction.ID).Error

	if err != nil {
		return nil, err
	}

	return &reloadedTrx, nil
}
func (r *transactionRepository) FindAll(ctx context.Context) (*[]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.WithContext(ctx).Preload("Category").Preload("User").Find(&transactions).Error
	return &transactions, err
}

func (r *transactionRepository) FindAllByUserID(ctx context.Context, userID uuid.UUID) (*[]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.WithContext(ctx).Preload("Category").Preload("User").Where("user_id = ?", userID).Find(&transactions).Error
	return &transactions, err // UBAH MUTLAK: Tambahkan & agar jadi pointer
}

func (r *transactionRepository) IsOwner(ctx context.Context, userID uuid.UUID, walletID string) bool {
	var wallet models.Wallet
	err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", walletID, userID).First(&wallet).Error
	return err == nil
}

func (r *transactionRepository) FindByID(ctx context.Context, transactionID uuid.UUID) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.WithContext(ctx).Preload("Category").Preload("User").First(&transaction, "id = ?", transactionID).Error
	return &transaction, err
}

func (r *transactionRepository) UpdateTransaction(ctx context.Context, transaction *models.Transaction) error {
	return r.db.WithContext(ctx).Save(transaction).Error
}

func (r *transactionRepository) UpdateTransactionWithWalletBallance(ctx context.Context, transaction *models.Transaction, delta int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(transaction).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.Wallet{}).
			Where("id = ?", transaction.WalletID).
			Update("balance", gorm.Expr("balance + ?", delta)).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *transactionRepository) SoftDeleteTransaction(ctx context.Context, transactionID uuid.UUID, delta int64, walletID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(transactionID).Error; err != nil {
			return err
		}

		return nil
	})
}
