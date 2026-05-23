package utils

import (
	"cashflow_gin/models"
	"context"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Bikin kunci unik supaya aman dari bentrok dengan package lain
type contextKey string

const (
	UserIDKey   contextKey = "user_id"
	UserRoleKey contextKey = "user_role"
)

func GetUserID(ctx context.Context) (uuid.UUID, error) {
	// 1. Bongkar penyamaran context standar jadi Gin Context
	ginCtx, ok := ctx.(*gin.Context)
	if !ok {
		return uuid.Nil, fmt.Errorf("fatal: context is not a gin.Context")
	}

	// 2. Ambil dari brankas Gin
	val, exists := ginCtx.Get("user_id")
	if !exists {
		return uuid.Nil, fmt.Errorf("user ID missing in context (Unauthorized)")
	}

	idStr, ok := val.(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid user ID format in context")
	}

	return uuid.Parse(idStr)
}

func GetUserRole(ctx context.Context) (models.UserRole, error) {
	// 1. Bongkar penyamaran
	ginCtx, ok := ctx.(*gin.Context)
	if !ok {
		return 0, fmt.Errorf("fatal: context is not a gin.Context")
	}

	// 2. Ambil dari brankas Gin
	val, exists := ginCtx.Get("user_role")
	if !exists {
		return 0, fmt.Errorf("user role missing in context")
	}

	// val pasti string karena di middleware lu udah ada: userRole = fmt.Sprintf("%v", rawRole)
	roleStr, ok := val.(string)
	if !ok {
		return 0, fmt.Errorf("invalid user role format in context")
	}

	roleInt, err := strconv.Atoi(roleStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse role number: %s", roleStr)
	}

	return models.UserRole(roleInt), nil
}

func GetUserInfo(ctx context.Context) (uuid.UUID, models.UserRole, error) {
	ginCtx, ok := ctx.(*gin.Context)
	if !ok {
		return uuid.Nil, 0, fmt.Errorf("fatal: context is not a gin.Context")
	}
	userId, exists := ginCtx.Get("user_id")
	if !exists {
		return uuid.Nil, 0, fmt.Errorf("user ID missing in context (Unauthorized)")
	}
	userRole, exists := ginCtx.Get("user_role")
	if !exists {
		return uuid.Nil, 0, fmt.Errorf("user role missing in context")
	}

	roleStr, ok := userRole.(string)
	if !ok {
		return uuid.Nil, 0, fmt.Errorf("invalid user role format in context")
	}

	idStr, ok := userId.(string)
	if !ok {
		return uuid.Nil, 0, fmt.Errorf("invalid user ID format in context")
	}

	roleInt, err := strconv.Atoi(roleStr)
	if err != nil {
		return uuid.Nil, 0, fmt.Errorf("failed to parse role number: %s", roleStr)
	}

	idUser, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, 0, fmt.Errorf("failed to parse user ID format in context")
	}
	return idUser, models.UserRole(roleInt), nil
}

func SafeStringDereference(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

func BoolPtr(b bool) *bool {
	return &b
}

func StringPtr(s string) *string {
	return &s
}

func UUIDPtrToStringPtr(u *uuid.UUID) *string {
	if u == nil {
		return nil
	}
	s := u.String()
	return &s
}

// Taruh di utils/helper.go
func IntPtr(i int) *int {
	return &i
}
