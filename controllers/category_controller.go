package controllers

import (
	"cashflow_gin/dto/request"
	"cashflow_gin/dto/response"
	"cashflow_gin/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
func (c *CategoryController) CreateDefaultCategories(ctx *gin.Context) {
	category, err := c.services.CreateDefaultCategories()
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			response.BaseResponse{
				Status:  false,
				Message: "Failed to create default categories",
				Errors:  err,
				Data:    nil,
			},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.BaseResponse{
			Status:  err == nil,
			Message: "Default categories created successfully",
			Errors:  err,
			Data:    category,
		},
	)
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
	roleClaim, exists := ctx.Get("user_role")
	if !exists {
		ctx.JSON(
			http.StatusUnauthorized,
			response.BaseResponse{
				Status:  false,
				Message: "Unauthorized access",
				Errors:  "Insufficient permissions",
				Data:    nil,
			},
		)
		return
	}

	cat, err := c.services.GetAllCategories(roleClaim.(float64))
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			response.BaseResponse{
				Status:  false,
				Message: "Failed to retrieve categories",
				Errors:  err.Error(),
				Data:    nil,
			},
		)
		return
	}
	ctx.JSON(
		http.StatusOK,
		response.BaseResponse{
			Status:  true,
			Message: "Success retrieve all categories",
			Errors:  nil,
			Data:    cat,
		},
	)
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
	userIDStr, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, response.BaseResponse{
			Status:  false,
			Message: "Unauthorized",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.BaseResponse{
			Status:  false,
			Message: "Invalid user ID",
		})
		return
	}

	var input request.CreateCategoryRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, response.BaseResponse{
			Status:  false,
			Message: "Invalid input",
			Errors:  err.Error(),
		})
		return
	}

	category, err := c.services.CreateMy(userID, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.BaseResponse{
			Status:  false,
			Message: "Failed to create category",
			Errors:  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, response.BaseResponse{
		Status:  true,
		Message: "Category created successfully",
		Data:    category,
	})
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
	userIDStr, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, response.BaseResponse{
			Status:  false,
			Message: "Unauthorized",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.BaseResponse{
			Status:  false,
			Message: "Invalid user ID",
		})
		return
	}

	categories, err := c.services.GetMine(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.BaseResponse{
			Status:  false,
			Message: "Failed to retrieve categories",
			Errors:  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.BaseResponse{
		Status:  true,
		Message: "Success retrieve my categories",
		Data:    categories,
	})
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
	userIDStr, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, response.BaseResponse{
			Status:  false,
			Message: "Unauthorized",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.BaseResponse{
			Status:  false,
			Message: "Invalid user ID",
		})
		return
	}

	categoryIDStr := ctx.Param("id")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.BaseResponse{
			Status:  false,
			Message: "Invalid category ID",
		})
		return
	}

	var input request.CreateCategoryRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, response.BaseResponse{
			Status:  false,
			Message: "Invalid input",
			Errors:  err.Error(),
		})
		return
	}

	category, err := c.services.UpdateById(userID, categoryID, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.BaseResponse{
			Status:  false,
			Message: "Failed to update category",
			Errors:  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.BaseResponse{
		Status:  true,
		Message: "Category updated successfully",
		Data:    category,
	})
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
	userIDStr, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, response.BaseResponse{
			Status:  false,
			Message: "Unauthorized",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.BaseResponse{
			Status:  false,
			Message: "Invalid user ID",
		})
		return
	}

	categoryIDStr := ctx.Param("id")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.BaseResponse{
			Status:  false,
			Message: "Invalid category ID",
		})
		return
	}

	err = c.services.DeleteById(userID, categoryID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.BaseResponse{
			Status:  false,
			Message: "Failed to delete category",
			Errors:  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.BaseResponse{
		Status:  true,
		Message: "Category deleted successfully",
	})
}
