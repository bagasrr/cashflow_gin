package controllers

import (
	"cashflow_gin/dto/response"
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
