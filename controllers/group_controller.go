package controllers

import (
	"cashflow_gin/dto/request"
	"cashflow_gin/dto/response"
	"cashflow_gin/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GroupController struct {
	services services.GroupService
}

type removeUser struct {
	UserID string `json:"user_id"`
}

func NewGroupController(service services.GroupService) *GroupController {
	return &GroupController{
		services: service,
	}
}

// GetGroupByID godoc
// @Summary      Get Group By ID
// @Description  Mendapatkan detail grup berdasarkan ID, termasuk anggota dan informasi dompet.
// @Tags         Groups
// @Accept       json
// @Produce      json
// @Param        id path string true "ID Grup"
// @Success      200 {object} response.BaseResponse{data=response.GroupResponse}
// @Failure      400 {object} response.BaseResponse
// @Failure      500 {object} response.BaseResponse
// @Security 	 BearerAuth
// @Router       /groups/{id} [get]
func (c *GroupController) GetGroupByID(ctx *gin.Context) {
	groupID := ctx.Param("id")
	if groupID == "" {
		ctx.JSON(http.StatusBadRequest, response.BaseResponse{
			Status:  false,
			Message: "Group ID is required",
			Errors:  nil,
			Data:    nil,
		})
		return
	}

	group, err := c.services.GetGroupByID(groupID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.BaseResponse{
			Status:  false,
			Message: "Failed to retrieve group",
			Errors:  err.Error(),
			Data:    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, response.BaseResponse{
		Status:  true,
		Message: "Group retrieved successfully",
		Data:    group,
	})
}

func (c *GroupController) GetAllGroups(ctx *gin.Context) {
	groups, err := c.services.GetAllGroups()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.BaseResponse{
			Status:  false,
			Message: "Failed to retrieve groups",
			Errors:  err.Error(),
			Data:    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, response.BaseResponse{
		Status:  true,
		Message: "Groups retrieved successfully",
		Data:    groups,
	})
}

// CreateGroup godoc
// @Summary      Create Group
// @Description  Membuat grup baru dengan anggota yang ditentukan.
// @Tags         Groups
// @Accept       json
// @Produce      json
// @Param        request body request.CreateGroupRequest true "request body"
// @Success      200 {object} response.BaseResponse{data=response.GroupResponse}
// @Failure      400 {object} response.BaseResponse
// @Failure      500 {object} response.BaseResponse
// @Security 	 BearerAuth
// @Router       /groups [post]
func (c *GroupController) CreateGroup(ctx *gin.Context) {
	var req request.CreateGroupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.BaseResponse{
			Status:  false,
			Message: "Invalid request body",
			Errors:  err.Error(),
			Data:    nil,
		})
		return
	}

	userIDClaim, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, response.BaseResponse{
			Status:  false,
			Message: "Unauthorized",
			Errors:  "User ID not found in token",
			Data:    nil,
		})
		return
	}

	userIDStr, ok := userIDClaim.(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, response.BaseResponse{
			Status:  false,
			Message: "Internal Server Error",
			Errors:  "User ID claim is not a string",
			Data:    nil,
		})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.BaseResponse{
			Status:  false,
			Message: "Internal Server Error",
			Errors:  "Invalid User ID format",
			Data:    nil,
		})
		return
	}

	newGroup, err := c.services.CreateGroup(userID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.BaseResponse{
			Status:  false,
			Message: "Failed to create group",
			Errors:  err.Error(),
			Data:    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, response.BaseResponse{
		Status:  true,
		Message: "Group created successfully",
		Data:    newGroup,
	})
}

// RemoveUserFromGroup godoc
// @Summary      Remove User From Group
// @Description  Menghapus pengguna dari grup.
// @Tags         Groups
// @Accept       json
// @Produce      json
// @Param        id path string true "ID Grup"
// @Success      200 {object} response.BaseResponse
// @Failure      400 {object} response.BaseResponse
// @Failure      500 {object} response.BaseResponse
// @Security 	 BearerAuth
// @Router       /groups/{id}/remove-user [patch]
func (c *GroupController) RemoveUserFromGroup(ctx *gin.Context) {
	groupID := ctx.Param("id")
	var input removeUser
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, response.BaseResponse{
			Status:  false,
			Message: "Invalid request body",
			Errors:  err.Error(),
			Data:    nil,
		})
		return
	}

	if groupID == "" || input.UserID == "" {
		ctx.JSON(http.StatusBadRequest, response.BaseResponse{
			Status:  false,
			Message: "Group ID and User ID are required",
			Errors:  nil,
			Data:    nil,
		})
		return
	}

	err := c.services.RemoveUserFromGroup(groupID, input.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.BaseResponse{
			Status:  false,
			Message: "Failed to remove user from group",
			Errors:  err.Error(),
			Data:    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, response.BaseResponse{
		Status:  true,
		Message: "User removed from group successfully",
		Data:    nil,
	})
}
