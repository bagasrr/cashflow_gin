package services

import (
	"cashflow_gin/dto/request"
	"cashflow_gin/dto/response"
	"cashflow_gin/models"
	"cashflow_gin/repository"
	"context"
	"errors"

	"github.com/google/uuid"
)

type CategoryService interface {
	CreateDefault(ctx context.Context, input *request.CreateCategoryRequest) (*response.CategoryResponse, error)
	CreateDefaultCategories(ctx context.Context) (*[]models.Category, error)
	GetAllCategories(ctx context.Context, userRole models.UserRole, limit, page int) (*[]response.CategoryResponse, error)

	CreateMy(ctx context.Context, userID uuid.UUID, input request.CreateCategoryRequest) (*response.CategoryResponse, error)
	GetMine(ctx context.Context, userID uuid.UUID) (*[]response.CategoryResponse, error)

	UpdateById(ctx context.Context, userID, categoryID uuid.UUID, input request.CreateCategoryRequest) (*response.CategoryResponse, error)
	DeleteById(ctx context.Context, userID, categoryID uuid.UUID) error
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(r repository.CategoryRepository) CategoryService {
	return &categoryService{repo: r}
}

func (s *categoryService) CreateDefault(ctx context.Context, input *request.CreateCategoryRequest) (*response.CategoryResponse, error) {
	category := models.Category{
		Name:   input.Name,
		Type:   input.Type,
		UserID: uuid.Nil,
	}
	err := s.repo.CreateDefault(ctx, &category)
	if err != nil {
		return nil, err
	}

	res := &response.CategoryResponse{
		ID:   category.ID,
		Name: category.Name,
		Type: category.Type,
	}

	return res, err
}

func (s *categoryService) CreateDefaultCategories(ctx context.Context) (*[]models.Category, error) {
	newCategory, err := s.repo.CreateDefaultCategories(ctx)
	return newCategory, err
}

func (s *categoryService) GetAllCategories(ctx context.Context, userRole models.UserRole, limit, page int) (*[]response.CategoryResponse, error) {
	if userRole > models.RoleModerator {
		return nil, errors.New("forbidden: access is denied")
	}

	if limit == 0 {
		limit = 10
	}
	if page == 0 {
		page = 1
	}
	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit

	cat, err := s.repo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	var res []response.CategoryResponse
	for _, category := range *cat {
		r := response.CategoryResponse{
			UserID: category.UserID.String(),
			Name:   category.Name,
			Type:   category.Type,
		}
		if category.GroupID != nil {
			r.GroupID = category.GroupID.String()
		}
		res = append(res, r)
	}

	return &res, err
}

func (s *categoryService) CreateMy(ctx context.Context, userID uuid.UUID, input request.CreateCategoryRequest) (*response.CategoryResponse, error) {
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

	createdCategory, err := s.repo.Create(ctx, &category)
	if err != nil {
		return nil, err
	}

	res := response.CategoryResponse{
		UserID: createdCategory.UserID.String(),
		Name:   createdCategory.Name,
		Type:   createdCategory.Type,
	}
	if createdCategory.GroupID != nil {
		res.GroupID = createdCategory.GroupID.String()
	}

	return &res, nil
}

func (s *categoryService) GetMine(ctx context.Context, userID uuid.UUID) (*[]response.CategoryResponse, error) {
	categories, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var res []response.CategoryResponse
	for _, category := range *categories {
		r := response.CategoryResponse{
			Name: category.Name,
			Type: category.Type,
		}
		if category.GroupID != nil {
			r.GroupID = category.GroupID.String()
		}
		res = append(res, r)
	}

	return &res, nil
}

func (s *categoryService) UpdateById(ctx context.Context, userID, categoryID uuid.UUID, input request.CreateCategoryRequest) (*response.CategoryResponse, error) {
	category, err := s.repo.FindByIDAndUserID(ctx, categoryID, userID)
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

	updatedCategory, err := s.repo.Update(ctx, category)
	if err != nil {
		return nil, err
	}

	res := response.CategoryResponse{
		UserID: updatedCategory.UserID.String(),
		Name:   updatedCategory.Name,
		Type:   updatedCategory.Type,
	}
	if updatedCategory.GroupID != nil {
		res.GroupID = updatedCategory.GroupID.String()
	}

	return &res, nil
}

func (s *categoryService) DeleteById(ctx context.Context, userID, categoryID uuid.UUID) error {
	category, err := s.repo.FindByIDAndUserID(ctx, categoryID, userID)
	if err != nil {
		return errors.New("category not found or unauthorized")
	}

	return s.repo.Delete(ctx, category)
}
