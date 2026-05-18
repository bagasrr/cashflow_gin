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

	Create(ctx context.Context, userID uuid.UUID, input request.CreateCategoryRequest) (*response.CategoryResponse, error)
	CreateMy(ctx context.Context, userID uuid.UUID, input request.CreateCategoryRequest) (*response.CategoryResponse, error)
	GetMine(ctx context.Context, userID uuid.UUID, page, limit int) (*[]models.Category, error)

	GetById(ctx context.Context, categoryID uuid.UUID) (*response.CategoryResponse, error)
	UpdateById(ctx context.Context, userID, categoryID uuid.UUID, input request.CreateCategoryRequest) (*response.CategoryResponse, error)
	DeleteById(ctx context.Context, userID, categoryID uuid.UUID) error
}

type categoryService struct {
	repo      repository.CategoryRepository
	groupRepo repository.GroupRepository
	userRepo  repository.UserRepository
}

func NewCategoryService(r repository.CategoryRepository, gr repository.GroupRepository, ur repository.UserRepository) CategoryService {
	return &categoryService{repo: r, groupRepo: gr, userRepo: ur}
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
		ID:      category.ID.String(),
		UserID:  category.UserID.String(),
		GroupID: category.GroupID.String(),
		Name:    category.Name,
		Type:    category.Type,
	}

	return res, err
}

func (s *categoryService) Create(ctx context.Context, userID uuid.UUID, input request.CreateCategoryRequest) (*response.CategoryResponse, error) {
	// Create
	// Cek User Role Admin atau bukan
	// Jika admin, auto buat category default ke repo
	// Jika bukan admin,
	// Apakah ada group id?
	// Jika ada, cek apakah user itu admin group tersebut?
	// Jika iya, buat category dengan group id tersebut
	// Jika tidak, return error unauthorized
	// Jika tidak ada, buat category dengan user id tersebut

	isRoleAdmin, err := s.userRepo.IsRoleAdmin(ctx, userID)
	if err != nil {
		return nil, err
	}
	var createdCategory *models.Category
	if isRoleAdmin {
		category := models.Category{
			Name:   input.Name,
			Type:   input.Type,
			UserID: uuid.Nil,
		}
		createdCategory, err = s.repo.Create(ctx, &category)
		if err != nil {
			return nil, err
		}
		// res := &response.CategoryResponse{
		// 	ID:      createdCategory.ID,
		// 	UserID:  createdCategory.UserID.String(),
		// 	GroupID: createdCategory.GroupID.String(),
		// 	Name:    createdCategory.Name,
		// 	Type:    createdCategory.Type,
		// }
		// return res, nil
	} else {
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
		createdCategory, err = s.repo.Create(ctx, &category)
		if err != nil {
			return nil, err
		}

		// res := response.CategoryResponse{
		// 	ID:     createdCategory.ID,
		// 	UserID: createdCategory.UserID.String(),
		// 	Name:   createdCategory.Name,
		// 	Type:   createdCategory.Type,
		// }
		// if createdCategory.GroupID != nil {
		// 	res.GroupID = createdCategory.GroupID.String()
		// }

		// return &res, nil
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
			ID:     category.ID.String(),
			UserID: category.UserID.String(),
			// GroupID: category.GroupID.String(),
			Name: category.Name,
			Type: category.Type,
		}
		if category.GroupID != nil {
			r.GroupID = category.GroupID.String()
		} else {
			r.GroupID = "nil"
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

func (s *categoryService) GetMine(ctx context.Context, userID uuid.UUID, page, limit int) (*[]models.Category, error) {
	offset := (page - 1) * limit
	categories, err := s.repo.FindByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (s *categoryService) GetById(ctx context.Context, categoryID uuid.UUID) (*response.CategoryResponse, error) {
	category, err := s.repo.FindByID(ctx, categoryID)
	if err != nil {
		return nil, errors.New("category not found")
	}

	res := response.CategoryResponse{
		ID:     category.ID.String(),
		UserID: category.UserID.String(),
		Name:   category.Name,
		Type:   category.Type,
	}
	if category.GroupID != nil {
		res.GroupID = category.GroupID.String()
	}

	return &res, nil
}

func (s *categoryService) UpdateById(ctx context.Context, userID, categoryID uuid.UUID, input request.CreateCategoryRequest) (*response.CategoryResponse, error) {
	// category, err := s.repo.FindByIDAndUserID(ctx, categoryID, userID)
	// if err != nil {
	// 	return nil, errors.New("category not found or unauthorized")
	// }

	// Flow Update by id:
	// 1. Cari Category berdasarkan id
	// 2. Apakah itu categori milik group?
	//    - Jika iya, apakah user adalah admin group tersebut?
	//   		- Jika ya, boleh mendelete category tersebut
	//   		- Jika tidak, return error unauthorized
	//    - Jika tidak, apakah category tersebut milik user tersebut?
	//   		- Jika ya, boleh mendelete category tersebut
	//   		- Jika tidak, return error unauthorized
	// 3. Jika salah satu kondisi di atas terpenuhi, lakukan update pada category tersebut.

	category, err := s.repo.FindByID(ctx, categoryID)
	if err != nil {
		return nil, errors.New("category not found")
	}

	category.Name = input.Name
	category.Type = input.Type

	if category.GroupID != nil {
		isGroupAdmin, err := s.groupRepo.IsGroupAdmin(ctx, *category.GroupID, userID)
		if err != nil {
			return nil, errors.New("failed to check group admin status")
		}
		if !isGroupAdmin {
			return nil, errors.New("unauthorized: user is not an admin of the group")
		}
	}

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
		ID:     updatedCategory.ID.String(),
		UserID: updatedCategory.UserID.String(),
		// GroupID: updatedCategory.GroupID.String(),
		Name: updatedCategory.Name,
		Type: updatedCategory.Type,
	}
	if updatedCategory.GroupID != nil {
		res.GroupID = updatedCategory.GroupID.String()
	}

	return &res, nil
}

func (s *categoryService) DeleteById(ctx context.Context, userID, categoryID uuid.UUID) error {
	// Flow Delete by id:
	// 1. Cari Category berdasarkan id
	// 2. Apakah itu categori milik group?
	//    - Jika iya, apakah user adalah admin group tersebut?
	//   	- Jika ya, boleh mendelete category tersebut
	//   	- Jika tidak, return error unauthorized
	//    - Jika tidak, apakah category tersebut milik user tersebut?
	//   		- Jika ya, boleh mendelete category tersebut
	//   		- Jika tidak, return error unauthorized
	// 3. Jika salah satu kondisi di atas terpenuhi, lakukan soft delete pada category tersebut.

	cat, err := s.repo.FindByID(ctx, categoryID)
	if err != nil {
		return errors.New("category not found")
	}

	if cat.GroupID != nil {
		isGroupAdmin, err := s.groupRepo.IsGroupAdmin(ctx, *cat.GroupID, userID)
		if err != nil {
			return errors.New("failed to check group admin status")
		}
		if !isGroupAdmin {
			return errors.New("unauthorized: user is not an admin of the group")
		}
	}

	if cat.UserID != userID {
		return errors.New("unauthorized: user is not the owner of the category")
	}

	return s.repo.Delete(ctx, cat)
}
