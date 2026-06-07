package services

import (
	"cashflow_gin/api"
	"cashflow_gin/models"
	"cashflow_gin/repository"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type CategoryService interface {
	CreateMyDefault(ctx context.Context, userId uuid.UUID, input api.CreateCategoryReq) (*models.Category, error)
	CreateSystemCategories(ctx context.Context) error
	GetAllCategories(ctx context.Context, userRole models.UserRole, limit, page int) (*[]models.Category, int64, error)
	GetSystemCategories(ctx context.Context, limit, page int) (*[]models.Category, int64, error)

	Create(ctx context.Context, userID uuid.UUID, input api.CreateCategoryReq) (*models.Category, error)
	CreateMy(ctx context.Context, userID uuid.UUID, input api.CreateCategoryReq) (*models.Category, error)
	GetMine(ctx context.Context, userID uuid.UUID, page, limit int) (*[]models.Category, error)

	GetById(ctx context.Context, categoryID, reqId uuid.UUID, reqRole models.UserRole) (*models.Category, error)
	UpdateById(ctx context.Context, userID, catId uuid.UUID, input api.UpdateCategoryReq) (*models.Category, error)
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

func (s *categoryService) CreateMyDefault(ctx context.Context, userId uuid.UUID, input api.CreateCategoryReq) (*models.Category, error) {
	category := models.Category{
		Name:   input.Name,
		Type:   input.Type,
		UserID: &userId,
	}
	err := s.repo.CreateMyDefault(ctx, &category)
	if err != nil {
		return nil, err
	}

	res := &models.Category{
		Base: models.Base{
			ID: category.ID,
		},
		UserID:  category.UserID,
		GroupID: category.GroupID,
		Name:    category.Name,
		Type:    category.Type,
	}

	return res, err
}

func (s *categoryService) Create(ctx context.Context, userID uuid.UUID, input api.CreateCategoryReq) (*models.Category, error) {
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
			UserID: nil,
		}
		createdCategory, err = s.repo.Create(ctx, &category)
		if err != nil {
			return nil, err
		}
		// res := &models.Category{
		// 	ID:      createdCategory.ID,
		// 	UserID:  createdCategory.UserID.String(),
		// 	GroupID: createdCategory.GroupID.String(),
		// 	Name:    createdCategory.Name,
		// 	Type:    createdCategory.Type,
		// }
		// return res, nil
	} else {
		category := models.Category{
			UserID: &userID,
			Name:   input.Name,
			Type:   input.Type,
		}

		if input.GroupId != nil {
			id, err := uuid.Parse(*input.GroupId)
			if err != nil {
				return nil, errors.New("invalid group id")
			}
			category.GroupID = &id
		}
		createdCategory, err = s.repo.Create(ctx, &category)
		if err != nil {
			return nil, err
		}

		// res := models.Category{
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
	res := models.Category{
		Base: models.Base{
			ID: createdCategory.ID,
		},
		UserID:  createdCategory.UserID,
		GroupID: createdCategory.GroupID,
		Name:    createdCategory.Name,
		Type:    createdCategory.Type,
	}

	return &res, nil
}

func (s *categoryService) CreateSystemCategories(ctx context.Context) error {
	err := s.repo.CreateSystemCategories(ctx)
	return err
}

// Tanda tangan fungsi berubah: nambahin return 'int64' buat totalItems
func (s *categoryService) GetAllCategories(ctx context.Context, userRole models.UserRole, limit, page int) (*[]models.Category, int64, error) {
	if userRole > models.RoleModerator {
		return nil, 0, errors.New("forbidden: access is denied")
	}

	if limit <= 0 { // Langsung tangkap angka minus juga pakai <=
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit

	// Panggil Repo yang sekarang mengembalikan 3 nilai
	cat, totalItems, err := s.repo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// RETURN LANGSUNG. Jangan pakai looping sampah.
	return cat, totalItems, nil
}

func (s *categoryService) GetSystemCategories(ctx context.Context, limit, page int) (*[]models.Category, int64, error) {
	if limit <= 0 { // Langsung tangkap angka minus juga pakai <=
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit
	cat, totalItems, err := s.repo.GetSystemCategories(ctx, limit, offset)
	return cat, totalItems, err
}

func (s *categoryService) CreateMy(ctx context.Context, userID uuid.UUID, input api.CreateCategoryReq) (*models.Category, error) {
	category := models.Category{
		UserID: &userID,
		Name:   input.Name,
		Type:   input.Type,
	}

	if input.GroupId != nil {
		id, err := uuid.Parse(*input.GroupId)
		if err != nil {
			return nil, errors.New("invalid group id")
		}
		category.GroupID = &id
	}

	createdCategory, err := s.repo.Create(ctx, &category)
	if err != nil {
		return nil, err
	}
	return createdCategory, nil
}

func (s *categoryService) GetMine(ctx context.Context, userID uuid.UUID, page, limit int) (*[]models.Category, error) {
	offset := (page - 1) * limit
	categories, err := s.repo.FindByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (s *categoryService) GetById(ctx context.Context, categoryID, reqId uuid.UUID, reqRole models.UserRole) (*models.Category, error) {

	category, err := s.repo.FindByID(ctx, categoryID)
	if err != nil {
		return nil, errors.New("category not found")
	}
	if *category.UserID != reqId || reqRole < models.RoleUser {
		return nil, errors.New("user is not authorized to create category")
	}

	res := models.Category{
		Base: models.Base{
			ID: category.ID,
		},
		UserID:  category.UserID,
		GroupID: category.GroupID,
		Name:    category.Name,
		Type:    category.Type,
	}

	return &res, nil
}

func (s *categoryService) UpdateById(ctx context.Context, userID, catId uuid.UUID, input api.UpdateCategoryReq) (*models.Category, error) {
	fmt.Println("sampe ke service updatebyid")
	category, err := s.repo.FindByID(ctx, catId)
	if err != nil {
		return nil, errors.New("category not found")
	}

	category.Name = input.Name
	category.Type = input.Type

	if category.GroupID != nil {
		fmt.Println("Berhasil masuk ke validasi groupid != nil")
		isGroupAdmin, err := s.groupRepo.IsGroupAdmin(ctx, *category.GroupID, userID)
		if err != nil {
			return nil, errors.New("failed to check group admin status")
		}
		if !isGroupAdmin {
			return nil, errors.New("unauthorized: user is not an admin of the group")
		}
		fmt.Println("Validasi groupid != nil SELESAI")
	}
	if category.UserID != nil {
		fmt.Println("MASUK VALIDA CAT.USERID != Nil")
		if *category.UserID != userID {
			return nil, errors.New("user is not authorized to update category")
		}
	}
	if input.GroupId != nil {
		groupID, err := uuid.Parse(*input.GroupId)
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

	res := models.Category{
		Base: models.Base{
			ID: updatedCategory.ID,
		},
		UserID:  updatedCategory.UserID,
		GroupID: updatedCategory.GroupID,
		Name:    updatedCategory.Name,
		Type:    updatedCategory.Type,
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

	if cat.UserID != &userID {
		return errors.New("unauthorized: user is not the owner of the category")
	}

	return s.repo.Delete(ctx, cat)
}
