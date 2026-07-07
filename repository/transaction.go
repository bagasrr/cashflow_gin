package repository

import (
	"cashflow_gin/models"
	"context"
	"fmt"
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
	GetTransactionsByWallet(ctx context.Context, userID, walletID uuid.UUID, startDate, endDate time.Time, limit, offset int, search, sortBy, sortOrder string) ([]models.Transaction, int64, error)
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

		// 1. KEMBALIKAN SALDO DOMPET (Reverse Balance)
		// Lu ngelewatin 'delta' dari Service.
		// Kalau di Service delta lu adalah negatif (misal: -100000),
		// maka gorm.Expr("balance + ?", delta) akan menjadi balance + (-100000) alias ngurangin saldo.
		if err := tx.Model(&models.Wallet{}).
			Where("id = ?", walletID).
			Updates(map[string]interface{}{
				"balance":    gorm.Expr("balance + ?", delta),
				"updated_at": time.Now(),
			}).Error; err != nil {
			return err
		}

		// 2. SOFT DELETE TRANSAKSI
		// Kasih tau GORM tabel apa yang mau dituju dengan mengirim pointer &models.Transaction{}
		if err := tx.Delete(&models.Transaction{}, "id = ?", transactionID).Error; err != nil {
			return err
		}

		return nil // Commit transaksi jika kedua operasi di atas sukses
	})
}
func (r *transactionRepository) GetTransactionsByWallet(ctx context.Context, userID, walletID uuid.UUID, startDate, endDate time.Time, limit, offset int, search, sortBy, sortOrder string) ([]models.Transaction, int64, error) {
	var transactions []models.Transaction
	var totalItems int64

	// 1. Fondasi Query (MUTLAK: Gunakan .Model agar GORM tahu tabel mana yang mau dihitung)
	query := r.db.WithContext(ctx).Model(&models.Transaction{}).
		Where("user_id = ? AND wallet_id = ? AND deleted_at IS NULL", userID, walletID).
		Where("date >= ? AND date <= ?", startDate, endDate)

	// 2. Tumpuk Lego Search (SEBELUM DIHITUNG!)
	if search != "" {
		query = query.Where("title ILIKE ?", "%"+search+"%")
	}

	// 3. EKSEKUSI COUNT MUTLAK
	// Di titik ini, query sudah berisi filter tanggal + search (jika ada).
	// Tapi belum tercemar oleh Limit dan Offset.
	if err := query.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	// 4. Tumpuk Lego Sorting
	// (Sorting ditaruh setelah Count karena mengurutkan data cuma buang-buang CPU database saat proses menghitung)
	orderClause := fmt.Sprintf("%s %s", sortBy, sortOrder)
	query = query.Order(orderClause)

	if sortBy != "date" {
		query = query.Order("created_at DESC")
	}

	// 5. Eksekusi Akhir (Pemotongan Data)
	err := query.Limit(limit).
		Offset(offset).
		Preload("Category").
		Find(&transactions).Error

	return transactions, totalItems, err
}
