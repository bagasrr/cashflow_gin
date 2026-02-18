package services

import (
	"cashflow_gin/dto/request"
	"cashflow_gin/dto/response"
	"cashflow_gin/models"
	"cashflow_gin/repository"
	"errors"

	"github.com/google/uuid"
)

type CategoryService interface {
	CreateDefaultCategories() (*[]models.Category, error)
	GetAllCategories(userRole float64) (*[]models.Category, error)

	CreateMy(userID uuid.UUID, input request.CreateCategoryRequest) (*response.CategoryResponse, error)
	GetMine(userID uuid.UUID) (*[]response.CategoryResponse, error)

	UpdateById(userID, categoryID uuid.UUID, input request.CreateCategoryRequest) (*response.CategoryResponse, error)
	DeleteById(userID, categoryID uuid.UUID) error
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(r repository.CategoryRepository) CategoryService {
	return &categoryService{repo: r}
}

func (s *categoryService) CreateDefaultCategories() (*[]models.Category, error) {
	newCategory, err := s.repo.CreateDefaultCategories()
	return newCategory, err
}

func (s *categoryService) GetAllCategories(userRole float64) (*[]models.Category, error) {
	if userRole > 2 {
		return nil, errors.New("forbidden: access is denied")
	}

	return s.repo.FindAll()
}

func (s *categoryService) CreateMy(userID uuid.UUID, input request.CreateCategoryRequest) (*response.CategoryResponse, error) {
	category := models.Category{
		UserID: userID,
		Name:   input.Name,
		Type:   input.Type,
	}

	if input.GroupID != "" {
		id, err := uuid.Parse(input.GroupID)
		if err != nil {
			return nil, errors.New("invalid group id")
		}
		category.GroupID = &id
	}

	createdCategory, err := s.repo.Create(&category)
	if err != nil {
		return nil, err
	}

	res := response.CategoryResponse{
		ID:     createdCategory.ID.String(),
		UserID: createdCategory.UserID.String(),
		Name:   createdCategory.Name,
		Type:   createdCategory.Type,
	}
	if createdCategory.GroupID != nil {
		res.GroupID = createdCategory.GroupID.String()
	}

	return &res, nil
}

func (s *categoryService) GetMine(userID uuid.UUID) (*[]response.CategoryResponse, error) {
	categories, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	var res []response.CategoryResponse
	for _, category := range *categories {
		r := response.CategoryResponse{
			ID:     category.ID.String(),
			UserID: category.UserID.String(),
			Name:   category.Name,
			Type:   category.Type,
		}
		if category.GroupID != nil {
			r.GroupID = category.GroupID.String()
		}
		res = append(res, r)
	}

	return &res, nil
}

func (s *categoryService) UpdateById(userID, categoryID uuid.UUID, input request.CreateCategoryRequest) (*response.CategoryResponse, error) {
	category, err := s.repo.FindByIDAndUserID(categoryID, userID)
	if err != nil {
		return nil, errors.New("category not found or unauthorized")
	}

	category.Name = input.Name
	category.Type = input.Type

	if input.GroupID != "" {
		groupID, err := uuid.Parse(input.GroupID)
		if err != nil {
			return nil, errors.New("invalid group id")
		}
		category.GroupID = &groupID
	} else {
		category.GroupID = nil
	}

	updatedCategory, err := s.repo.Update(category)
	if err != nil {
		return nil, err
	}

	res := response.CategoryResponse{
		ID:     updatedCategory.ID.String(),
		UserID: updatedCategory.UserID.String(),
		Name:   updatedCategory.Name,
		Type:   updatedCategory.Type,
	}
	if updatedCategory.GroupID != nil {
		res.GroupID = updatedCategory.GroupID.String()
	}

	return &res, nil
}

func (s *categoryService) DeleteById(userID, categoryID uuid.UUID) error {
	category, err := s.repo.FindByIDAndUserID(categoryID, userID)
	if err != nil {
		return errors.New("category not found or unauthorized")
	}

	return s.repo.Delete(category)
}
