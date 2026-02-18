package repository

import (
	"cashflow_gin/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	FindByID(id uuid.UUID) (*models.Category, error)
	CreateDefaultCategories() (*[]models.Category, error)
	FindAll() (*[]models.Category, error)
	FindByName(name string) (*models.Category, error)
	FindByUserID(userID uuid.UUID) (*[]models.Category, error)
	FindByGroupID(groupID uuid.UUID) (*[]models.Category, error)
	Create(category *models.Category) (*models.Category, error)
	Update(category *models.Category) (*models.Category, error)
	Delete(category *models.Category) error
	FindByIDAndUserID(id uuid.UUID, userID uuid.UUID) (*models.Category, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) FindByID(id uuid.UUID) (*models.Category, error) {
	var category models.Category
	// Cek kategori berdasarkan ID dan pastikan user_id juga cocok (security)
	// Note: Parameter userID bisa lu tambah nanti buat validasi ownership
	err := r.db.First(&category, "id = ?", id).Error
	return &category, err
}

func (r *categoryRepository) CreateDefaultCategories() (*[]models.Category, error) {
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
	return &categories, r.db.Create(&categories).Error
}

func (r *categoryRepository) FindAll() (*[]models.Category, error) {
	var categories []models.Category
	err := r.db.Find(&categories).Error
	return &categories, err
}

func (r *categoryRepository) FindByName(name string) (*models.Category, error) {
	var category models.Category
	err := r.db.First(&category, "name = ?", name).Error
	return &category, err
}

func (r *categoryRepository) FindByUserID(userID uuid.UUID) (*[]models.Category, error) {
	var categories []models.Category
	err := r.db.Where("user_id = ?", userID).Find(&categories).Error
	return &categories, err
}

func (r *categoryRepository) FindByGroupID(groupID uuid.UUID) (*[]models.Category, error) {
	var categories []models.Category
	err := r.db.Where("group_id = ?", groupID).Find(&categories).Error
	return &categories, err
}

func (r *categoryRepository) Create(category *models.Category) (*models.Category, error) {
	err := r.db.Create(&category).Error
	return category, err
}

func (r *categoryRepository) Update(category *models.Category) (*models.Category, error) {
	err := r.db.Save(&category).Error
	return category, err
}

func (r *categoryRepository) Delete(category *models.Category) error {
	err := r.db.Delete(&category).Error
	return err
}

func (r *categoryRepository) FindByIDAndUserID(id uuid.UUID, userID uuid.UUID) (*models.Category, error) {
	var category models.Category
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&category).Error
	return &category, err
}
