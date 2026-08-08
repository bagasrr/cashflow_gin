package repository

import (
	"cashflow_gin/models"
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreateWithWalletUpdate(ctx context.Context, transaction *models.Transaction, impact int64) (*models.Transaction, error)
	FindAll(ctx context.Context) (*[]models.Transaction, error)
	FindAllByUserID(ctx context.Context, userID uuid.UUID) (*[]models.Transaction, error)
	IsOwner(ctx context.Context, userID uuid.UUID, walletID string) bool
	FindByID(ctx context.Context, transactionID uuid.UUID) (*models.Transaction, error)
	UpdateTransaction(ctx context.Context, transaction *models.Transaction, delta int64) error
	SoftDeleteTransaction(ctx context.Context, transactionID uuid.UUID, delta int64, walletID uuid.UUID) error
	GetTransactionsByWallet(ctx context.Context, userID, walletID uuid.UUID, startDate, endDate time.Time, limit, offset int, search, sortBy, sortOrder string) ([]models.Transaction, int64, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

// WAJIB: Tambahkan parameter "impact int64"
func (r *transactionRepository) CreateWithWalletUpdate(ctx context.Context, transaction *models.Transaction, impact int64) (*models.Transaction, error) {

	// 1. MULAI TRANSAKSI DATABASE (ACID)
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// A. Simpan Transaksi (MUTLAK: Omit relasi agar GORM tidak menciptakan data siluman)
		if err := tx.Omit("Category", "User").Create(transaction).Error; err != nil {
			return err
		}

		// B. Update Saldo Dompet (Gunakan variabel IMPACT, bukan Amount)
		if impact != 0 {
			if err := tx.Model(&models.Wallet{}).
				Where("id = ?", transaction.WalletID).
				Updates(map[string]interface{}{
					"balance":    gorm.Expr("balance + ?", impact), // impact sudah membawa nilai minus jika EXPENSE
					"updated_at": time.Now(),
				}).Error; err != nil {
				return err
			}
		}

		return nil
	})

	// 2. CEK STATUS TRANSAKSI
	if err != nil {
		return nil, err
	}

	// 3. RELOAD SECARA MUTLAK
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
	err := r.db.WithContext(ctx).Preload("Category").Preload("User").Preload("Wallet").First(&transaction, "id = ?", transactionID).Error
	return &transaction, err
}

// Namanya dipersingkat, tapi logic-nya meng-handle semua skenario (dengan atau tanpa perubahan saldo)
func (r *transactionRepository) UpdateTransaction(ctx context.Context, transaction *models.Transaction, delta int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// 1. OMIT MUTLAK (Fokus ubah tabel transaksi aja)
		if err := tx.Omit("Category", "User").Save(transaction).Error; err != nil {
			return err
		}

		// 2. UPDATE SALDO DOMPET (Hanya kalau ada perubahan nominal/tipe)
		if delta != 0 {
			if err := tx.Model(&models.Wallet{}).
				Where("id = ?", transaction.WalletID).
				Update("balance", gorm.Expr("balance + ?", delta)).Error; err != nil {
				return err
			}
		}

		// 3. REFRESH ABSOLUT (Bantai memori basi)
		var refreshed models.Transaction
		if err := tx.Preload("Category").Preload("User").First(&refreshed, transaction.ID).Error; err != nil {
			return err
		}

		*transaction = refreshed
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

	// 1. FONDASI MUTLAK DENGAN JOIN
	// Paksa GORM untuk menggabungkan tabel sejak awal agar Postgres bisa membaca kolom categories
	query := r.db.WithContext(ctx).Model(&models.Transaction{}).
		Joins("JOIN categories ON categories.id = transactions.category_id").
		Where("transactions.user_id = ? AND transactions.wallet_id = ? AND transactions.deleted_at IS NULL", userID, walletID).
		Where("transactions.date >= ? AND transactions.date <= ?", startDate, endDate)

	// 2. TUMPUK LEGO SEARCH (Gunakan prefix tabel!)
	if search != "" {
		query = query.Where("transactions.title ILIKE ?", "%"+search+"%")
	}

	// 3. EKSEKUSI COUNT
	if err := query.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	// 4. SANITASI MUTLAK & TRANSLASI SORTING (Mencegah SQL Injection)
	// Kita translasi bahasa API Frontend menjadi bahasa Database yang akurat
	var orderColumn string
	switch sortBy {
	case "categories.type": // Jika Frontend ngirim "category"
		orderColumn = "categories.type" // Kita urutkan berdasarkan tipe kategori
	case "categories.name": // Jika Frontend ngirim "category"
		orderColumn = "categories.name" // Kita urutkan berdasarkan tipe kategori
	case "amount":
		orderColumn = "transactions.amount"
	case "title":
		orderColumn = "transactions.title"
	case "description":
		orderColumn = "transactions.description"
	case "updated_at":
		orderColumn = "transactions.updated_at"
	default:
		orderColumn = "transactions.date" // Fallback aman
	}

	// Validasi Sort Order secara hardcode agar tidak bisa disusupi perintah SQL
	if sortOrder != "asc" && sortOrder != "ASC" {
		sortOrder = "DESC"
	}

	// Terapkan sorting yang sudah 100% aman
	query = query.Order(orderColumn + " " + sortOrder)

	// Tie-breaker order
	if orderColumn != "transactions.date" {
		query = query.Order("transactions.created_at DESC")
	}

	// 5. EKSEKUSI AKHIR
	err := query.Limit(limit).
		Offset(offset).
		Preload("Category"). // Preload tetap wajib dipanggil agar struct Golang lu terisi data utuh
		Preload("Wallet").
		Find(&transactions).Error

	return transactions, totalItems, err
}
