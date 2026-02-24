package controllers

import (
	"cashflow_gin/dto/response"
	"cashflow_gin/models"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// --- SHARED HELPERS (Usable by all controllers) ---

func GetUserID(ctx *gin.Context) (uuid.UUID, error) {
	userIDClaim, exists := ctx.Get("user_id")
	if !exists {
		return uuid.Nil, fmt.Errorf("user ID missing in context")
	}
	// Jika middleware kamu sudah simpan dalam bentuk uuid.UUID, cast ke uuid.UUID langsung.
	// Jika masih string, parse dulu:
	return uuid.Parse(fmt.Sprintf("%v", userIDClaim))
}

func GetUserRole(ctx *gin.Context) (models.UserRole, error) {
	roleClaim, exists := ctx.Get("user_role")
	if !exists {
		return 0, fmt.Errorf("user role missing in context")
	}

	// JWT MapClaims usually unmarshals numbers as float64
	if roleFloat, ok := roleClaim.(float64); ok {
		return models.UserRole(int8(roleFloat)), nil
	}

	// In case it's stored as another type (e.g., int, string)
	return 0, fmt.Errorf("invalid role format in context")
}

func GetParamID(ctx *gin.Context, key string) (uuid.UUID, error) {
	idStr := ctx.Param(key)
	if idStr == "" {
		return uuid.Nil, fmt.Errorf("param %s is empty", key)
	}
	return uuid.Parse(idStr)
}

func SendSuccess(ctx *gin.Context, code int, message string, data interface{}) {
	ctx.JSON(code, response.BaseResponse{
		Status:  true,
		Message: message,
		Data:    data,
	})
}

func SendError(ctx *gin.Context, code int, message string, err error) {
	errVal := ""
	if err != nil {
		errVal = err.Error()
	}
	ctx.JSON(code, response.BaseResponse{
		Status:  false,
		Message: message,
		Errors:  errVal,
	})
}
