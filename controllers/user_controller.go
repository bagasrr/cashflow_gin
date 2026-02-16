package controllers

import (
	"cashflow_gin/dto/response"
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
	users, err := c.service.FindAllUser()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			response.BaseResponse{
				Status:  false,
				Message: "error",
				Errors:  err,
			})
		return
	}

	ctx.JSON(http.StatusOK, response.BaseResponse{
		Status:  true,
		Message: "success",
		Data:    users,
		Errors:  nil,
	})
}
