package repository

import (
	"cashflow_gin/models"
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WalletRepository interface {
	CreateWallet(ctx context.Context, wallet models.Wallet) (*models.Wallet, error)
	FindAll(ctx context.Context, offset int, limit int) (*[]models.Wallet, error)
	FindByID(ctx context.Context, walletID uuid.UUID) (*models.Wallet, error)
	FindAllMine(ctx context.Context, userID uuid.UUID, limit int, offset int) ([]models.Wallet, int64, error)
	SoftDeleteWallet(ctx context.Context, wallet uuid.UUID) error
	GetWalletByID(ctx context.Context, walletID uuid.UUID) (*models.Wallet, error)
	UpdateWallet(ctx context.Context, wallet *models.Wallet) (*models.Wallet, error)
	GetWalletChartData(ctx context.Context, userID, walletID uuid.UUID, startDate, endDate time.Time) ([]models.WalletChartPoint, error)
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

// HAPUS BINTANG DI RETURN TYPE
func (r *walletRepository) FindAllMine(ctx context.Context, userID uuid.UUID, limit int, offset int) ([]models.Wallet, int64, error) {
	var wallets []models.Wallet
	var totalItems int64

	// 1. EKSEKUSI A: MENGHITUNG TOTAL
	if err := r.db.WithContext(ctx).
		Model(&models.Wallet{}).
		Where("user_id = ?", userID).
		Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	// 2. EKSEKUSI B: MENARIK DATA WALLET
	// Kita hapus Preload di sini karena GORM cacat kalau disuruh Limit per grup
	if err := r.db.WithContext(ctx).
		Model(&models.Wallet{}).
		Select("wallets.*, "+countQuery+" as transaction_count").
		Where("user_id = ?", userID).
		Limit(limit).Offset(offset).
		Order("created_at DESC"). // Tambahkan urutan wallet yang logis!
		Find(&wallets).Error; err != nil {
		return nil, 0, err
	}

	// 3. EKSEKUSI C: CONTROLLED MANUAL INJECTION (Solusi Limit 5 Per Dompet)
	// Karena wallets dibatasi limit (misal 10), loop ini maksimal hanya jalan 10 kali. Sangat aman.
	for i := range wallets {
		var recentTransactions []models.Transaction

		// Tarik 5 transaksi terbaru KHUSUS untuk dompet ini
		r.db.WithContext(ctx).
			Preload("Category").
			Where("wallet_id = ?", wallets[i].ID).
			Order("created_at DESC"). // Atau pakai "date DESC" jika lu pakai kolom Date
			Limit(5).
			Find(&recentTransactions)

		// Suntikkan ke memori array dompet
		wallets[i].Transactions = recentTransactions
	}

	// Kembalikan MURNI (tanpa &)
	return wallets, totalItems, nil
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

func (r *walletRepository) GetWalletChartData(ctx context.Context, userID, walletID uuid.UUID, startDate, endDate time.Time) ([]models.WalletChartPoint, error) {
	var chartPoints []models.WalletChartPoint

	// Query untuk mengelompokkan pemasukan dan pengeluaran per hari
	err := r.db.WithContext(ctx).
		Table("transactions").
		Select(`
            transactions.date::date as tx_date,
            COALESCE(SUM(CASE WHEN categories.type = 'INCOME' THEN transactions.amount ELSE 0 END), 0) as total_income,
            COALESCE(SUM(CASE WHEN categories.type = 'EXPENSE' THEN transactions.amount ELSE 0 END), 0) as total_expense,
            COALESCE(SUM(CASE WHEN categories.type = 'INVESTMENT' THEN transactions.amount ELSE 0 END), 0) as total_investment
        `).
		Joins("JOIN categories ON categories.id = transactions.category_id").
		Where("transactions.user_id = ? AND transactions.wallet_id = ? AND transactions.deleted_at IS NULL", userID, walletID).
		Where("transactions.date >= ? AND transactions.date <= ?", startDate, endDate).
		Group("transactions.date::date").
		Order("transactions.date::date ASC"). // Urutkan dari tanggal tertua ke terbaru untuk grafik
		Scan(&chartPoints).Error

	return chartPoints, err
}
