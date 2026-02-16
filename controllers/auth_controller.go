package controllers

import (
	"cashflow_gin/dto/request"
	"cashflow_gin/dto/response"
	"cashflow_gin/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service services.AuthService
}

func NewAuthController(s services.AuthService) *AuthController {
	return &AuthController{service: s}
}

// Register godoc
// @Summary      Register User
// @Description  Membuat user baru sekaligus wallet default.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body request.CreateUserRequest true "request body"
// @Success      201 {object} response.BaseResponse{data=response.UserResponse}
// @Failure      500 {object} response.BaseResponse
// @Router       /auth/register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var input request.CreateUserRequest

	// 1. Validasi Input JSON
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			response.BaseResponse{
				Status:  false,
				Message: "Input Tidak Valid",
				Errors:  err.Error(),
			},
		)
		return
	}

	// 2. Panggil Service
	user, err := c.service.Register(input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.BaseResponse{
			Status:  false,
			Message: "Internal Server Error",
			Errors:  err.Error(),
		})
		return
	}

	// 3. Kirim Response
	ctx.JSON(http.StatusCreated, response.BaseResponse{
		Status:  true,
		Message: "success",
		Data:    user,
	})
}

// Login godoc
// @Summary      User Login
// @Description  Autentikasi user dan mendapatkan token JWT.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body request.LoginRequest true "request body"
// @Success      200 {object} response.BaseResponse{data=string}
// @Failure      400 {object} response.BaseResponse
// @Failure      500 {object} response.BaseResponse
// @Router       /auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var input request.LoginRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, response.BaseResponse{
			Status:  false,
			Message: "Input tidak Valid",
			Errors:  err.Error(),
		})
		return
	}

	token, err := c.service.Login(&input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.BaseResponse{
			Status:  false,
			Message: "Error",
			Errors:  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.BaseResponse{
		Status:  true,
		Message: "Login Success",
		Data: struct {
			Token string `json:"token"`
		}{
			Token: token,
		},
	})
}
