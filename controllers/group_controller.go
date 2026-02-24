package controllers

import (
	"cashflow_gin/dto/request"
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
	groupID, err := GetParamID(ctx, "id")
	if err != nil {
		SendError(ctx, http.StatusBadRequest, "Invalid group ID", err)
		return
	}

	group, err := c.services.GetGroupByID(ctx.Request.Context(), groupID)
	if err != nil {
		SendError(ctx, http.StatusInternalServerError, "Failed to retrieve group", err)
		return
	}

	SendSuccess(ctx, http.StatusOK, "Group retrieved successfully", group)
}

func (c *GroupController) GetAllGroups(ctx *gin.Context) {
	groups, err := c.services.GetAllGroups(ctx.Request.Context())
	if err != nil {
		SendError(ctx, http.StatusInternalServerError, "Failed to retrieve groups", err)
		return
	}

	SendSuccess(ctx, http.StatusOK, "Groups retrieved successfully", groups)
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
		SendError(ctx, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	userID, err := GetUserID(ctx)
	if err != nil {
		SendError(ctx, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	newGroup, err := c.services.CreateGroup(ctx.Request.Context(), userID, req)
	if err != nil {
		SendError(ctx, http.StatusInternalServerError, "Failed to create group", err)
		return
	}

	SendSuccess(ctx, http.StatusOK, "Group created successfully", newGroup)
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
	groupID, err := GetParamID(ctx, "id")
	if err != nil {
		SendError(ctx, http.StatusBadRequest, "Invalid group ID", err)
		return
	}

	var input removeUser
	if err := ctx.ShouldBindJSON(&input); err != nil {
		SendError(ctx, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	userUUID, err := uuid.Parse(input.UserID)
	if err != nil {
		SendError(ctx, http.StatusBadRequest, "Invalid User ID format", err)
		return
	}

	err = c.services.RemoveUserFromGroup(ctx.Request.Context(), groupID, userUUID)
	if err != nil {
		SendError(ctx, http.StatusInternalServerError, "Failed to remove user from group", err)
		return
	}

	SendSuccess(ctx, http.StatusOK, "User removed from group successfully", nil)
}

// DeleteGroup godoc
// @Summary      Delete Group
// @Description  Menghapus grup berdasarkan ID.
// @Tags         Groups
// @Accept       json
// @Produce      json
// @Param        id path string true "ID Grup"
// @Success      200 {object} response.BaseResponse
// @Failure      400 {object} response.BaseResponse
// @Failure      500 {object} response.BaseResponse
// @Security 	 BearerAuth
// @Router       /groups/{id} [delete]
func (c *GroupController) DeleteGroup(ctx *gin.Context) {
	groupID, err := GetParamID(ctx, "id")
	if err != nil {
		SendError(ctx, http.StatusBadRequest, "Invalid group ID", err)
		return
	}

	err = c.services.DeleteGroup(ctx.Request.Context(), groupID)
	if err != nil {
		SendError(ctx, http.StatusInternalServerError, "Failed to delete group", err)
		return
	}

	SendSuccess(ctx, http.StatusOK, "Group deleted successfully", nil)
}
