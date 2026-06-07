package handler

import (
	"cashflow_gin/api" // Ini manggil hasil generate YAML (sesuaikan path kalau salah)
	"cashflow_gin/models"
	"cashflow_gin/services"
	"cashflow_gin/utils"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
)

// 1. Ini struct yang dicari-cari sama server.go lu
type CategoryAPI struct {
	Service services.CategoryService
}

// 2. Ini adalah fungsi WAJIB yang dipaksa oleh interface dari hasil oapi-codegen
func (c *CategoryAPI) CreateCategory(ctx context.Context, req api.CreateCategoryRequestObject) (api.CreateCategoryResponseObject, error) {
	// Mapping input dari DTO otomatis ke DTO manual lu

	userID, err := utils.GetUserID(ctx)
	log.Printf("UserID di Context: %s", userID)

	if err != nil {
		return api.CreateCategory500JSONResponse{
			Status:  utils.BoolPtr(false),
			Message: utils.StringPtr("Cannot get user context"),
			Errors:  utils.StringPtr("ERR : " + err.Error()),
		}, nil
	}

	input := *req.Body
	// Panggil layer Service lu yang lama
	res, err := c.Service.Create(ctx, userID, input)
	if err != nil {
		// DETEKSI ERROR DUPLIKASI DARI DATABASE
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") ||
			strings.Contains(err.Error(), "SQLSTATE 23505") {

			log.Println("❌ ERROR VALIDASI: Nama Kategori Sudah Ada")
			// WAJIB RETURN 400 (Bad Request), bukan 500.
			// Pastikan di openapi.yaml lu udah definisiin response 400 untuk endpoint ini.
			return api.CreateCategory400JSONResponse{
				Status:  utils.BoolPtr(false),
				Message: utils.StringPtr("Nama Kategori Sudah Ada"),
			}, nil
		}
		return api.CreateCategory500JSONResponse{
			Status:  utils.BoolPtr(false),
			Message: utils.StringPtr("Huhuuu, Cannot Create Category"),
			Errors:  utils.StringPtr("ERR : " + err.Error()),
		}, nil
	}

	return api.CreateCategory201JSONResponse{
		Data: &api.CategoryRes{
			Id:      res.ID.String(),
			Name:    res.Name,
			Type:    string(res.Type),
			UserId:  utils.UUIDPtrToStringPtr(res.UserID),
			GroupId: utils.UUIDPtrToStringPtr(res.GroupID),
		},
		Message: utils.StringPtr("Create Category Success"), // Assign the pointer to a string variable
		Status:  utils.BoolPtr(true),
	}, nil
}

func (c *CategoryAPI) GetSystemCategories(ctx context.Context, request api.GetSystemCategoriesRequestObject) (api.GetSystemCategoriesResponseObject, error) {
	limitValue := 10
	pageValue := 1

	if request.Params.Limit != 0 {
		limitValue = request.Params.Limit
	}

	if request.Params.Page != 0 {
		pageValue = request.Params.Page
	}
	_, err := utils.GetUserID(ctx)
	if err != nil {
		return api.GetSystemCategories500JSONResponse{
			Status:  utils.BoolPtr(false),
			Errors:  utils.StringPtr("ERR : " + err.Error()),
			Message: utils.StringPtr("Cannot Parse User COntext"),
		}, nil
	}
	cat, totalItems, err := c.Service.GetSystemCategories(ctx, limitValue, pageValue)

	var res []api.CategoryRes
	for _, v := range *cat {
		res = append(res, api.CategoryRes{
			Id:      v.ID.String(),
			GroupId: utils.UUIDPtrToStringPtr(v.GroupID),
			UserId:  utils.UUIDPtrToStringPtr(v.UserID),
			Name:    v.Name,
			Type:    string(v.Type),
		})
	}

	totalPages := (int(totalItems) + limitValue - 1) / limitValue

	return api.GetSystemCategories200JSONResponse{
		Data:    &res,
		Message: utils.StringPtr("Get Categories Success"),
		Status:  utils.BoolPtr(true),
		Meta: &api.PaginationMeta{
			CurrentPage: utils.IntPtr(pageValue),
			TotalItems:  utils.IntPtr(int(totalItems)),
			TotalPages:  utils.IntPtr(totalPages),
		},
	}, nil
}

func (c *CategoryAPI) GetCategories(ctx context.Context, request api.GetCategoriesRequestObject) (api.GetCategoriesResponseObject, error) {
	// Coba ambil pake string manual buat ngetes
	role, err := utils.GetUserRole(ctx)
	if err != nil {
		log.Println("❌ ERROR AUTH:", err)

		status := false
		msg := "Gagal Auth: " + err.Error()
		return api.GetCategories500JSONResponse{Status: &status, Message: &msg}, nil
	}
	limitValue := 10
	pageValue := 1

	if request.Params.Limit != 0 {
		limitValue = request.Params.Limit
	}

	if request.Params.Page != 0 {
		pageValue = request.Params.Page
	}

	cat, totalItems, err := c.Service.GetAllCategories(ctx, role, limitValue, pageValue)
	if err != nil {
		status := false
		msg := "Gagal Database: " + err.Error()
		return api.GetCategories500JSONResponse{Status: &status, Message: &msg}, nil
	}

	var res []api.CategoryRes
	for _, v := range *cat {
		res = append(res, api.CategoryRes{
			Id:      v.ID.String(),
			GroupId: utils.UUIDPtrToStringPtr(v.GroupID),
			UserId:  utils.UUIDPtrToStringPtr(v.UserID),
			Name:    v.Name,
			Type:    string(v.Type),
		})
	}

	totalPages := (int(totalItems) + limitValue - 1) / limitValue

	log.Println("✅ Get Categories Success")

	// CARA INISIALISASI ANONYMOUS STRUCT POINTER
	return api.GetCategories200JSONResponse{
		Data:    &res,
		Message: utils.StringPtr("Get Categories Success"),
		Status:  utils.BoolPtr(true),
		Meta: &api.PaginationMeta{
			CurrentPage: utils.IntPtr(pageValue),
			TotalItems:  utils.IntPtr(int(totalItems)),
			TotalPages:  utils.IntPtr(totalPages),
		},
	}, nil
}

func (c *CategoryAPI) CreateSystemCategories(ctx context.Context, request api.CreateSystemCategoriesRequestObject) (api.CreateSystemCategoriesResponseObject, error) {
	_, role, err := utils.GetUserInfo(ctx)
	if err != nil {
		return api.CreateSystemCategories401JSONResponse{
			Status:  utils.BoolPtr(false),
			Message: utils.StringPtr("Cannot get context"),
			Errors:  utils.StringPtr("ERR :  " + err.Error()),
		}, nil
	}
	if role != models.RoleAdmin {
		return api.CreateSystemCategories401JSONResponse{
			Status:  utils.BoolPtr(false),
			Message: utils.StringPtr("Access Denied : Only System Admin can run this endpoint"),
		}, nil
	}
	err = c.Service.CreateSystemCategories(ctx)
	if err != nil {
		return api.CreateSystemCategories500JSONResponse{
			Status:  utils.BoolPtr(false),
			Message: utils.StringPtr("Gagal Auth: " + err.Error()),
			Errors:  utils.StringPtr("ERR :  " + err.Error()),
		}, nil
	}
	return api.CreateSystemCategories201JSONResponse{
		Status:  utils.BoolPtr(true),
		Message: utils.StringPtr("Create Default System Categories Success"),
	}, nil
}

func (c *CategoryAPI) GetMyCategories(ctx context.Context, request api.GetMyCategoriesRequestObject) (api.GetMyCategoriesResponseObject, error) {
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		status := false
		msg := "Gagal Auth: " + err.Error()
		return api.GetMyCategories500JSONResponse{
			Status:  &status,
			Message: &msg,
		}, nil
	}

	pageValue := 1
	limitValue := 10

	if request.Params.Limit != 0 {
		limitValue = request.Params.Limit
	}

	if request.Params.Page != 0 {
		pageValue = request.Params.Page
	}
	cat, err := c.Service.GetMine(ctx, userID, pageValue, limitValue)
	if err != nil {
		status := false
		msg := "Gagal mengambil kategori: " + err.Error()
		return api.GetMyCategories500JSONResponse{
			Status:  &status,
			Message: &msg,
		}, nil
	}
	var res []api.CategoryRes
	for _, v := range *cat {

		res = append(res, api.CategoryRes{
			Id:      v.ID.String(),
			UserId:  utils.UUIDPtrToStringPtr(v.UserID),
			GroupId: utils.UUIDPtrToStringPtr(v.GroupID),
			Name:    v.Name,
			Type:    string(v.Type),
		})
	}
	return api.GetMyCategories200JSONResponse{
		Status:  utils.BoolPtr(true),
		Message: utils.StringPtr("Get My Categories Success"),
		Data:    &res,
	}, nil
}

func (c *CategoryAPI) DeleteCategory(ctx context.Context, request api.DeleteCategoryRequestObject) (api.DeleteCategoryResponseObject, error) {
	userId, _, err := utils.GetUserInfo(ctx)
	fmt.Println("Requestor ID in handler : " + userId.String())
	if err != nil {
		return api.DeleteCategory400JSONResponse{
			Message: utils.StringPtr("User Id not Found In The context, Please Login Frist"),
			Status:  utils.BoolPtr(false),
			Errors:  utils.StringPtr("ERR :  " + err.Error()),
		}, nil
	}
	catId, err := uuid.Parse(request.Id)
	if err != nil {
		return api.DeleteCategory400JSONResponse{
			Message: utils.StringPtr("User Id not Found In The context, Please Login Frist"),
			Status:  utils.BoolPtr(false),
			Errors:  utils.StringPtr("ERR :  " + err.Error()),
		}, nil
	}

	err = c.Service.DeleteById(ctx, userId, catId)
	if err != nil {
		return api.DeleteCategory500JSONResponse{
			Message: utils.StringPtr("Delete Category Failed: " + err.Error()),
			Status:  utils.BoolPtr(false),
			Errors:  utils.StringPtr("ERR :  " + err.Error()),
		}, nil
	}
	return api.DeleteCategory200JSONResponse{
		Message: utils.StringPtr("Delete Category Success"),
		Status:  utils.BoolPtr(true),
	}, nil
}

func (c *CategoryAPI) GetCategoryById(ctx context.Context, request api.GetCategoryByIdRequestObject) (api.GetCategoryByIdResponseObject, error) {
	requestorId, requestorRole, err := utils.GetUserInfo(ctx)
	catId, err := uuid.Parse(request.Id)
	if err != nil {
		return api.GetCategoryById400JSONResponse{
			Status:  utils.BoolPtr(false),
			Message: utils.StringPtr("Cant Parse Id"),
		}, nil
	}
	cat, err := c.Service.GetById(ctx, catId, requestorId, requestorRole)
	if err != nil {
		return api.GetCategoryById500JSONResponse{
			Status:  utils.BoolPtr(false),
			Message: utils.StringPtr("Cant Get Category"),
			Errors:  utils.StringPtr("ERR : " + err.Error()),
		}, nil
	}
	res := api.CategoryRes{
		Id:      cat.ID.String(),
		UserId:  utils.UUIDPtrToStringPtr(cat.UserID),
		GroupId: utils.UUIDPtrToStringPtr(cat.GroupID),
		Name:    cat.Name,
		Type:    string(cat.Type),
	}
	return api.GetCategoryById200JSONResponse{
		Data:    &res,
		Status:  utils.BoolPtr(true),
		Message: utils.StringPtr("Get Category Success"),
	}, nil
}

func (c *CategoryAPI) UpdateCategory(ctx context.Context, request api.UpdateCategoryRequestObject) (api.UpdateCategoryResponseObject, error) {
	requestorId, err := utils.GetUserID(ctx)
	if err != nil {
		return api.UpdateCategory400JSONResponse{
			Message: utils.StringPtr("Gagal Auth: " + err.Error()),
			Status:  utils.BoolPtr(false),
		}, nil
	}
	catId, err := uuid.Parse(request.Id)
	if err != nil {
		return api.UpdateCategory400JSONResponse{
			Message: utils.StringPtr("Tidak Bisa Mendapatkan :id " + err.Error()),
			Status:  utils.BoolPtr(false),
		}, nil
	}
	inputBody := *request.Body
	cat, err := c.Service.UpdateById(ctx, requestorId, catId, inputBody)
	if err != nil {
		return api.UpdateCategory500JSONResponse{
			Message: utils.StringPtr("Cannot Update Category"),
			Status:  utils.BoolPtr(false),
			Errors:  utils.StringPtr("ERR : " + err.Error()),
		}, nil
	}
	res := api.CategoryRes{
		Id:      cat.ID.String(),
		GroupId: utils.UUIDPtrToStringPtr(cat.GroupID),
		Name:    cat.Name,
		Type:    string(cat.Type),
		UserId:  utils.UUIDPtrToStringPtr(cat.UserID),
	}
	return api.UpdateCategory200JSONResponse{
		Data:    &res,
		Message: utils.StringPtr("Update Category Success"),
		Status:  utils.BoolPtr(true),
	}, nil
}
