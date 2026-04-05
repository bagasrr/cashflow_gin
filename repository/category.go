package repository

import (
	"cashflow_gin/models"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *models.Category) (*models.Category, error)
	CreateDefault(ctx context.Context, cat *models.Category) error
	CreateDefaultCategories(ctx context.Context) (*[]models.Category, error)
	FindAll(ctx context.Context, limit, offset int) (*[]models.Category, error)

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

func (r *categoryRepository) CreateDefault(ctx context.Context, cat *models.Category) error {
	return r.db.WithContext(ctx).Create(&cat).Error
}

func (r *categoryRepository) CreateDefaultCategories(ctx context.Context) (*[]models.Category, error) {
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
	return &categories, r.db.WithContext(ctx).Create(&categories).Error
}

func (r *categoryRepository) FindAll(ctx context.Context, limit int, offset int) (*[]models.Category, error) {
	var categories []models.Category
	err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Find(&categories).Error
	return &categories, err
}

func (r *categoryRepository) FindByName(ctx context.Context, name string) (*models.Category, error) {
	var category models.Category
	err := r.db.WithContext(ctx).First(&category, "name = ?", name).Error
	return &category, err
}

func (r *categoryRepository) FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) (*[]models.Category, error) {
	var categories []models.Category
	err := r.db.WithContext(ctx).Where("user_id = ? ", userID).Limit(limit).Offset(offset).Find(&categories).Error
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
