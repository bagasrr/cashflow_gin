package controllers

import (
	"cashflow_gin/dto/response"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// --- PRIVATE HELPERS (DRY IMPLEMENTATION) ---

func (c *TransactionController) getUserID(ctx *gin.Context) (uuid.UUID, error) {
	userIDClaim, exists := ctx.Get("user_id")
	if !exists {
		return uuid.Nil, fmt.Errorf("user ID missing in context")
	}
	// Jika middleware kamu sudah simpan dalam bentuk uuid.UUID, cast ke uuid.UUID langsung.
	// Jika masih string, parse dulu:
	return uuid.Parse(fmt.Sprintf("%v", userIDClaim))
}

func (c *TransactionController) getParamID(ctx *gin.Context, key string) (uuid.UUID, error) {
	idStr := ctx.Param(key)
	if idStr == "" {
		return uuid.Nil, fmt.Errorf("param %s is empty", key)
	}
	return uuid.Parse(idStr)
}

func (c *TransactionController) sendSuccess(ctx *gin.Context, message string, data interface{}) {
	ctx.JSON(http.StatusOK, response.BaseResponse{
		Status:  true,
		Message: message,
		Data:    data,
	})
}

func (c *TransactionController) sendError(ctx *gin.Context, code int, message string, err error) {
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
