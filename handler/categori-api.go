package handler

import (
	"cashflow_gin/api" // Ini manggil hasil generate YAML (sesuaikan path kalau salah)
	"cashflow_gin/dto/request"
	"cashflow_gin/services"
	"context"
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

	// Panggil layer Service lu yang lama
	res, err := c.Service.CreateDefault(ctx, reqInput)
	if err != nil {
		// Kalau service error, return response 500 sesuai kontrak YAML
		return api.CreateCategory500JSONResponse{}, nil
	}

	message := "Create Category Success"
	status := true
	return api.CreateCategory201JSONResponse{
		Data: &api.CategoryRes{
			Id:   &res.ID,
			Name: &res.Name,
			Type: &res.Type,
		},
		Message: &message, // Assign the pointer to a string variable
		Status:  &status,
	}, nil
}

func (c *CategoryAPI) GetDefaultCategories(ctx context.Context, request api.GetDefaultCategoriesRequestObject) (api.GetDefaultCategoriesResponseObject, error) {
	return api.GetDefaultCategories200JSONResponse{}, nil
}

func (c *CategoryAPI) GetCategories(ctx context.Context, request api.GetCategoriesRequestObject) (api.GetCategoriesResponseObject, error) {
	return api.GetCategories200JSONResponse{}, nil
}

func (c *CategoryAPI) CreateDefaultCategories(ctx context.Context, request api.CreateDefaultCategoriesRequestObject) (api.CreateDefaultCategoriesResponseObject, error) {
	return api.CreateDefaultCategories201JSONResponse{}, nil
}

func (c *CategoryAPI) GetMyCategories(ctx context.Context, request api.GetMyCategoriesRequestObject) (api.GetMyCategoriesResponseObject, error) {
	return api.GetMyCategories200JSONResponse{}, nil
}

func (c *CategoryAPI) DeleteCategory(ctx context.Context, request api.DeleteCategoryRequestObject) (api.DeleteCategoryResponseObject, error) {
	return api.DeleteCategory200JSONResponse{}, nil
}

func (c *CategoryAPI) GetCategoryById(ctx context.Context, request api.GetCategoryByIdRequestObject) (api.GetCategoryByIdResponseObject, error) {
	return api.GetCategoryById200JSONResponse{}, nil
}

func (c *CategoryAPI) UpdateCategory(ctx context.Context, request api.UpdateCategoryRequestObject) (api.UpdateCategoryResponseObject, error) {
	return api.UpdateCategory200JSONResponse{}, nil
}
