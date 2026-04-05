package controllers

import (
	"cashflow_gin/dto/request"
	"cashflow_gin/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryController struct {
	services services.CategoryService
}

func NewCategoryController(s services.CategoryService) *CategoryController {
	return &CategoryController{
		services: s,
	}
}

// CreateDefaultCategories godoc
// @Summary      Create Default Categories
// @Description  Membuat kategori default untuk pengguna baru.
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Success      200 {object} response.BaseResponse{data=[]response.CategoryResponse}
// @Failure      500 {object} response.BaseResponse
// @Security 	 BearerAuth
// @Router       /categories/default [post]
func (c *CategoryController) CreateDefault(ctx *gin.Context) {
	var input request.CreateCategoryRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		SendError(ctx, http.StatusBadRequest, "Invalid input", err)
		return
	}

	category, err := c.services.CreateDefault(ctx.Request.Context(), &input)
	if err != nil {
		SendError(ctx, http.StatusInternalServerError, "Failed to create default categories", err)
		return
	}

	SendSuccess(ctx, http.StatusOK, "Default categories created successfully", category)
}

// CreateDefaultCategories godoc
// @Summary      Create Default Categories
// @Description  Membuat kategori default untuk pengguna baru.
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Success      200 {object} response.BaseResponse{data=[]response.CategoryResponse}
// @Failure      500 {object} response.BaseResponse
// @Security 	 BearerAuth
// @Router       /categories/default [post]
func (c *CategoryController) CreateDefaultCategories(ctx *gin.Context) {
	category, err := c.services.CreateDefaultCategories(ctx.Request.Context())
	if err != nil {
		SendError(ctx, http.StatusInternalServerError, "Failed to create default categories", err)
		return
	}

	SendSuccess(ctx, http.StatusOK, "Default categories created successfully", category)
}

// Get All Categories godoc
// @Summary      Get All Categories
// @Description  Mendapatkan semua kategori yang tersedia.
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Success      200 {object} response.BaseResponse{data=[]response.CategoryResponse}
// @Failure      500 {object} response.BaseResponse
// @Security 	 BearerAuth
// @Router       /categories [get]
func (c *CategoryController) GetAllCategories(ctx *gin.Context) {
	role, err := GetUserRole(ctx)
	if err != nil {
		SendError(ctx, http.StatusUnauthorized, "Unauthorized access", err)
		return
	}

	cat, err := c.services.GetAllCategories(ctx.Request.Context(), role, 100, 1) // Default limit 100, page 1
	if err != nil {
		SendError(ctx, http.StatusInternalServerError, "Failed to retrieve categories", err)
		return
	}
	SendSuccess(ctx, http.StatusOK, "Success retrieve all categories", cat)
}

// CreateMy godoc
// @Summary      Create My Category
// @Description  Membuat kategori baru untuk pengguna saat ini.
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param        request body request.CreateCategoryRequest true "Create Category Request"
// @Success      201 {object} response.BaseResponse
// @Failure      400 {object} response.BaseResponse
// @Failure      401 {object} response.BaseResponse
// @Security 	 BearerAuth
// @Router       /categories/mine [post]
func (c *CategoryController) CreateMy(ctx *gin.Context) {
	userID, err := GetUserID(ctx)
	if err != nil {
		SendError(ctx, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	var input request.CreateCategoryRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		SendError(ctx, http.StatusBadRequest, "Invalid input", err)
		return
	}

	category, err := c.services.CreateMy(ctx.Request.Context(), userID, input)
	if err != nil {
		SendError(ctx, http.StatusInternalServerError, "Failed to create category", err)
		return
	}

	SendSuccess(ctx, http.StatusCreated, "Category created successfully", category)
}

// GetMine godoc
// @Summary      Get My Categories
// @Description  Mendapatkan semua kategori milik pengguna saat ini.
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Success      200 {object} response.BaseResponse
// @Failure      401 {object} response.BaseResponse
// @Security 	 BearerAuth
// @Router       /categories/mine [get]
func (c *CategoryController) GetMine(ctx *gin.Context) {
	userID, err := GetUserID(ctx)
	if err != nil {
		SendError(ctx, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	categories, err := c.services.GetMine(ctx.Request.Context(), userID)
	if err != nil {
		SendError(ctx, http.StatusInternalServerError, "Failed to retrieve categories", err)
		return
	}

	SendSuccess(ctx, http.StatusOK, "Success retrieve my categories", categories)
}

// UpdateById godoc
// @Summary      Update Category
// @Description  Memperbarui kategori berdasarkan ID.
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param        id path string true "Category ID"
// @Param        request body request.CreateCategoryRequest true "Update Category Request"
// @Success      200 {object} response.BaseResponse
// @Failure      400 {object} response.BaseResponse
// @Failure      404 {object} response.BaseResponse
// @Security 	 BearerAuth
// @Router       /categories/{id} [put]
func (c *CategoryController) UpdateById(ctx *gin.Context) {
	userID, err := GetUserID(ctx)
	if err != nil {
		SendError(ctx, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	categoryID, err := GetParamID(ctx, "id")
	if err != nil {
		SendError(ctx, http.StatusBadRequest, "Invalid category ID", err)
		return
	}

	var input request.CreateCategoryRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		SendError(ctx, http.StatusBadRequest, "Invalid input", err)
		return
	}

	category, err := c.services.UpdateById(ctx.Request.Context(), userID, categoryID, input)
	if err != nil {
		SendError(ctx, http.StatusInternalServerError, "Failed to update category", err)
		return
	}

	SendSuccess(ctx, http.StatusOK, "Category updated successfully", category)
}

// DeleteById godoc
// @Summary      Delete Category
// @Description  Menghapus kategori berdasarkan ID (Soft Delete).
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param        id path string true "Category ID"
// @Success      200 {object} response.BaseResponse
// @Failure      400 {object} response.BaseResponse
// @Failure      404 {object} response.BaseResponse
// @Security 	 BearerAuth
// @Router       /categories/{id} [patch]
func (c *CategoryController) DeleteById(ctx *gin.Context) {
	userID, err := GetUserID(ctx)
	if err != nil {
		SendError(ctx, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	categoryID, err := GetParamID(ctx, "id")
	if err != nil {
		SendError(ctx, http.StatusBadRequest, "Invalid category ID", err)
		return
	}

	err = c.services.DeleteById(ctx.Request.Context(), userID, categoryID)
	if err != nil {
		SendError(ctx, http.StatusInternalServerError, "Failed to delete category", err)
		return
	}

	SendSuccess(ctx, http.StatusOK, "Category deleted successfully", nil)
}
