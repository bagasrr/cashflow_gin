package services

import (
	"cashflow_gin/models"
	"cashflow_gin/repository"
	"errors"
)

type CategoryService interface {
	CreateDefaultCategories() ([]models.Category, error)
	GetAllCategories(userRole float64) ([]models.Category, error)
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(r repository.CategoryRepository) CategoryService {
	return &categoryService{repo: r}
}

func (s *categoryService) CreateDefaultCategories() ([]models.Category, error) {
	newCategory, err := s.repo.CreateDefaultCategories()
	return newCategory, err
}

func (s *categoryService) GetAllCategories(userRole float64) ([]models.Category, error) {
	if userRole > 2 {
		return nil, errors.New("forbidden: access is denied")
	}

	return s.repo.FindAll()
}
