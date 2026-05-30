package repository

import (
	"cashflow_gin/models"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	CreateSystemCategories(ctx context.Context) error
	GetSystemCategories(ctx context.Context, limit, offset int) (*[]models.Category, int64, error)

	Create(ctx context.Context, category *models.Category) (*models.Category, error)
	CreateMyDefault(ctx context.Context, cat *models.Category) error
	FindAll(ctx context.Context, limit, offset int) (*[]models.Category, int64, error)

	FindByName(ctx context.Context, name string) (*models.Category, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) (*[]models.Category, error)
	FindByGroupID(ctx context.Context, groupID uuid.UUID) (*[]models.Category, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.Category, error)
	Update(ctx context.Context, category *models.Category) (*models.Category, error)
	Delete(ctx context.Context, category *models.Category) error
	// FindByIDAndUserID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*models.Category, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	var category models.Category
	// Cek kategori berdasarkan ID dan pastikan user_id juga cocok (security)
	// Note: Parameter userID bisa lu tambah nanti buat validasi ownership
	err := r.db.WithContext(ctx).First(&category, "id = ?", id).Error
	return &category, err
}

func (r *categoryRepository) CreateMyDefault(ctx context.Context, cat *models.Category) error {
	return r.db.WithContext(ctx).Create(&cat).Error
}

func (r *categoryRepository) CreateSystemCategories(ctx context.Context) error {
	var count int64

	// 1. GERBANG PENGECEKAN (O(1) Query)
	// Cek apakah tabel sudah punya kategori sistem (user_id IS NULL)
	err := r.db.WithContext(ctx).
		Model(&models.Category{}).
		Where("user_id IS NULL").
		Count(&count).Error

	if err != nil {
		return err
	}

	// 2. LOGIKA PENANGKAL
	// Kalau datanya udah ada (count > 0), langsung keluar. Jangan lanjut insert.
	if count > 0 {
		// Lu bisa return nil (menganggap seeding sudah sukses di masa lalu)
		// atau return error kalau lu emang mau ngasih warning keras.
		return nil
	}

	// 3. EKSEKUSI INSERT
	categories := []models.Category{
		{Name: "Salary", Type: "INCOME"},
		{Name: "Freelance", Type: "INCOME"},
		{Name: "Payment Received", Type: "INCOME"},
		{Name: "Gift", Type: "INCOME"},

		{Name: "Groceries", Type: "EXPENSE"},
		{Name: "Food", Type: "EXPENSE"},
		{Name: "Clothing", Type: "EXPENSE"},
		{Name: "Debt", Type: "EXPENSE"},
		{Name: "Subscription", Type: "EXPENSE"},
		{Name: "Utilities", Type: "EXPENSE"},
		{Name: "Transport", Type: "EXPENSE"},
		{Name: "Entertainment", Type: "EXPENSE"},
	}

	// Gunakan pointer nil untuk UserID dan GroupID (karena ini kategori sistem)
	for i := range categories {
		categories[i].UserID = nil
		categories[i].GroupID = nil
	}

	return r.db.WithContext(ctx).Create(&categories).Error
}

// Tanda tangan fungsi WAJIB diubah untuk me-return totalItems (int64)
func (r *categoryRepository) GetSystemCategories(ctx context.Context, limit, offset int) (*[]models.Category, int64, error) {
	var categories []models.Category
	var totalItems int64

	// 1. Definisikan tabel dan filter utama terlebih dahulu (Sebelum dieksekusi!)
	// Gunakan "IS NULL" karena Kategori Sistem/Default tidak punya pemilik.
	query := r.db.WithContext(ctx).Model(&models.Category{}).Where("user_id IS NULL")

	// 2. Eksekusi COUNT untuk mendapatkan total keseluruhan data (wajib untuk paginasi Meta)
	if err := query.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	// 3. Eksekusi FIND dengan Limit dan Offset untuk mengambil data di halaman ini saja
	if err := query.Limit(limit).Offset(offset).Find(&categories).Error; err != nil {
		return nil, 0, err
	}

	// Kembalikan 3 nilai: Data, Total Data, dan Error
	return &categories, totalItems, nil
}

// Tanda tangan fungsi berubah: nambahin return 'int64' buat totalItems
func (r *categoryRepository) FindAll(ctx context.Context, limit int, offset int) (*[]models.Category, int64, error) {
	var categories []models.Category
	var totalItems int64

	// 1. Definisikan tabel base-nya (Model)
	query := r.db.WithContext(ctx).Model(&models.Category{})

	// 2. Eksekusi perhitungan TOTAL KESELURUHAN (tanpa limit/offset)
	if err := query.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	// 3. Eksekusi pencarian data untuk HALAMAN INI (dengan limit/offset)
	if err := query.Limit(limit).Offset(offset).Find(&categories).Error; err != nil {
		return nil, 0, err
	}

	// Kembalikan datanya dan angka totalItems-nya
	return &categories, totalItems, nil
}

func (r *categoryRepository) FindByName(ctx context.Context, name string) (*models.Category, error) {
	var category models.Category
	err := r.db.WithContext(ctx).First(&category, "name = ?", name).Error
	return &category, err
}

func (r *categoryRepository) FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) (*[]models.Category, error) {
	var categories []models.Category
	err := r.db.WithContext(ctx).Where("user_id = ? OR user_id IS NULL", userID).Limit(limit).Offset(offset).Find(&categories).Error
	return &categories, err
}

func (r *categoryRepository) FindByGroupID(ctx context.Context, groupID uuid.UUID) (*[]models.Category, error) {
	var categories []models.Category
	err := r.db.WithContext(ctx).Where("group_id = ?", groupID).Find(&categories).Error
	return &categories, err
}

func (r *categoryRepository) Create(ctx context.Context, category *models.Category) (*models.Category, error) {
	err := r.db.WithContext(ctx).Create(&category).Error
	return category, err
}

func (r *categoryRepository) Update(ctx context.Context, category *models.Category) (*models.Category, error) {
	err := r.db.WithContext(ctx).Save(&category).Error
	return category, err
}

func (r *categoryRepository) Delete(ctx context.Context, category *models.Category) error {
	err := r.db.WithContext(ctx).Delete(&category).Error
	return err
}

// func (r *categoryRepository) FindByIDAndUserID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*models.Category, error) {
// 	var category models.Category
// 	err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&category).Error
// 	return &category, err
// }
