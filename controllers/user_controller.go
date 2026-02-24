package controllers

import (
	"cashflow_gin/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service services.UserService
}

func NewUserController(s services.UserService) *UserController {
	return &UserController{service: s}
}

// FindAllUser godoc
// @Summary      Find All User
// @Description  Mendapatkan daftar semua user. Admin Only can access this endpoint.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Success      200 {object} response.BaseResponse{data=response.UserResponse}
// @Failure		 500 {object} response.BaseResponse
// @Security 	 BearerAuth
// @Router       /users/ [get]
func (c *UserController) FindAllUser(ctx *gin.Context) {
	users, err := c.service.FindAllUser(ctx.Request.Context())
	if err != nil {
		SendError(ctx, http.StatusInternalServerError, "error", err)
		return
	}

	SendSuccess(ctx, http.StatusOK, "success", users)
}

func (c *UserController) GetMyProfile(ctx *gin.Context) {
	// 1 baris untuk ambil UserID
	userID, err := GetUserID(ctx)
	if err != nil {
		SendError(ctx, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	user, err := c.service.GetMyProfile(ctx.Request.Context(), userID)
	if err != nil {
		SendError(ctx, http.StatusInternalServerError, "Failed to retrieve profile", err)
		return
	}

	SendSuccess(ctx, http.StatusOK, "Profile retrieved successfully", user)
}
