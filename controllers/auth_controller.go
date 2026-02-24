package controllers

import (
	"cashflow_gin/dto/request"
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
		SendError(ctx, http.StatusBadRequest, "Input Tidak Valid", err)
		return
	}

	// 2. Panggil Service
	user, err := c.service.Register(ctx.Request.Context(), input)
	if err != nil {
		SendError(ctx, http.StatusInternalServerError, "Internal Server Error", err)
		return
	}

	// 3. Kirim Response
	SendSuccess(ctx, http.StatusCreated, "success", user)
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
		SendError(ctx, http.StatusBadRequest, "Input tidak Valid", err)
		return
	}

	token, err := c.service.Login(ctx.Request.Context(), &input)
	if err != nil {
		SendError(ctx, http.StatusInternalServerError, "Error", err)
		return
	}

	SendSuccess(ctx, http.StatusOK, "Login Success", struct {
		Token string `json:"token"`
	}{
		Token: token,
	})
}
