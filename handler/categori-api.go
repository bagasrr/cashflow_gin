package handler

import (
	"cashflow_gin/api" // Ini manggil hasil generate YAML (sesuaikan path kalau salah)
	"cashflow_gin/dto/request"
	"cashflow_gin/services"
	"cashflow_gin/utils"
	"context"
	"log"
	"strings"
)

// 1. Ini struct yang dicari-cari sama server.go lu
type CategoryAPI struct {
	Service services.CategoryService
}

// 2. Ini adalah fungsi WAJIB yang dipaksa oleh interface dari hasil oapi-codegen
func (c *CategoryAPI) CreateCategory(ctx context.Context, req api.CreateCategoryRequestObject) (api.CreateCategoryResponseObject, error) {
	// Mapping input dari DTO otomatis ke DTO manual lu
	reqInput := &request.CreateCategoryRequest{
		Name: req.Body.Name, // Hati-hati, hasil generate YAML biasanya berupa pointer
		Type: req.Body.Type,
	}
	if req.Body.GroupId != nil {
		reqInput.GroupID = *req.Body.GroupId
	}

	userID, err := utils.GetUserID(ctx)
	log.Printf("UserID di Context: %s", userID)

	if err != nil {
		status := false
		msg := "Gagal Auth: " + err.Error()
		return api.CreateCategory500JSONResponse{
			Status:  &status,
			Message: &msg,
		}, nil
	}

	// Panggil layer Service lu yang lama
	res, err := c.Service.Create(ctx, userID, *reqInput)
	if err != nil {
		// DETEKSI ERROR DUPLIKASI DARI DATABASE
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") ||
			strings.Contains(err.Error(), "SQLSTATE 23505") {

			log.Println("❌ ERROR VALIDASI: Nama Kategori Sudah Ada")
			status := false
			msg := "Kategori dengan nama tersebut sudah ada."
			// WAJIB RETURN 400 (Bad Request), bukan 500.
			// Pastikan di openapi.yaml lu udah definisiin response 400 untuk endpoint ini.
			return api.CreateCategory400JSONResponse{
				Status:  &status,
				Message: &msg,
			}, nil
		}

		// Kalau error lain (misal koneksi DB mati), baru return 500
		status := false
		msg := "Gagal membuat kategori: " + err.Error()
		return api.CreateCategory500JSONResponse{
			Status:  &status,
			Message: &msg,
		}, nil
	}

	message := "Create Category Success"
	status := true
	return api.CreateCategory201JSONResponse{
		Data: &api.CategoryRes{
			Id:   res.ID,
			Name: res.Name,
			Type: res.Type,
		},
		Message: &message, // Assign the pointer to a string variable
		Status:  &status,
	}, nil
}

func (c *CategoryAPI) GetDefaultCategories(ctx context.Context, request api.GetDefaultCategoriesRequestObject) (api.GetDefaultCategoriesResponseObject, error) {
	return api.GetDefaultCategories200JSONResponse{}, nil
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

	if request.Params.Limit != nil {
		limitValue = *request.Params.Limit
	}

	if request.Params.Page != nil {
		pageValue = *request.Params.Page
	}
	cat, err := c.Service.GetAllCategories(ctx, role, limitValue, pageValue)
	if err != nil {
		status := false
		msg := "Gagal Database: " + err.Error()
		return api.GetCategories500JSONResponse{Status: &status, Message: &msg}, nil
	}
	var res []api.CategoryRes
	for _, v := range *cat {
		res = append(res, api.CategoryRes{
			Id:      v.ID,
			UserId:  &v.UserID,
			GroupId: &v.GroupID,
			Name:    v.Name,
			Type:    v.Type,
		})
	}
	log.Println("✅ Get Categories Success")
	return api.GetCategories200JSONResponse{
		Data: &res,
	}, nil
}

func (c *CategoryAPI) CreateDefaultCategories(ctx context.Context, request api.CreateDefaultCategoriesRequestObject) (api.CreateDefaultCategoriesResponseObject, error) {
	return api.CreateDefaultCategories201JSONResponse{}, nil
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

	if request.Params.Limit != nil {
		limitValue = *request.Params.Limit
	}

	if request.Params.Page != nil {
		pageValue = *request.Params.Page
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
		uid := v.UserID.String()
		guid := v.GroupID.String()
		res = append(res, api.CategoryRes{
			Id:      v.ID.String(),
			UserId:  &uid,
			GroupId: &guid,
			Name:    v.Name,
			Type:    v.Type,
		})
	}
	return api.GetMyCategories200JSONResponse{
		Status:  utils.BoolPtr(true),
		Message: utils.StringPtr("Get My Categories Success"),
		Data:    &res,
	}, nil
}

func (c *CategoryAPI) DeleteCategory(ctx context.Context, request api.DeleteCategoryRequestObject) (api.DeleteCategoryResponseObject, error) {
	return api.DeleteCategory200JSONResponse{}, nil
}

func (c *CategoryAPI) GetCategoryById(ctx context.Context, request api.GetCategoryByIdRequestObject) (api.GetCategoryByIdResponseObject, error) {
	return api.GetCategoryById200JSONResponse{}, nil
}

func (c *CategoryAPI) UpdateCategory(ctx context.Context, request api.UpdateCategoryRequestObject) (api.UpdateCategoryResponseObject, error) {
	return api.UpdateCategory201JSONResponse{}, nil
}
