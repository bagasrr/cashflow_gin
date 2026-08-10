package utils

import (
	"cashflow_gin/models"
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
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

func ValidatePagination(page, limit int) (validPage int, validLimit int, offset int) {
	// 1. Sanitasi Limit
	if limit <= 0 {
		limit = 10 // Default
	} else if limit > 100 {
		limit = 100 // Hard-cap maksimal biar server gak jebol
	}

	// 2. Sanitasi Page
	if page <= 0 {
		page = 1
	}

	// 3. Kalkulasi Offset
	offset = (page - 1) * limit

	return page, limit, offset
}

// Helper function murni
func AbsInt64(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}

func GetStringFromContext(ctx context.Context, key string) (string, error) {
	ginCtx, ok := ctx.(*gin.Context)
	if !ok {
		return "", errors.New("fatal: context is not a gin.Context")
	}
	val, exists := ginCtx.Get(key)
	if !exists {
		return "", errors.New("fatal: context missing from context")
	}
	valStr, ok := val.(string)
	if !ok {
		return "", errors.New("fatal: invalid value in context")
	}

	return valStr, nil
}

func GetFloatFromContext(ctx context.Context, key string) (float64, error) {
	ginCtx, ok := ctx.(*gin.Context)
	if !ok {
		return 0, errors.New("fatal: context is not a gin.Context")
	}

	val, exists := ginCtx.Get(key)
	if !exists {
		return 0, errors.New("fatal: context missing from context")
	}

	// EKSEKUSI MUTLAK: Langsung tembak ke float64, karena dari Middleware asalnya udah float64!
	f, ok := val.(float64)
	if !ok {
		// Kasih pesan error yang jelas biar kalau meledak lagi lu tahu tipe datanya apa
		return 0, fmt.Errorf("fatal: invalid value in context, expected float64 but got %T", val)
	}

	return f, nil
}

const txKey contextKey = "tx_gorm"

// InjectTx menyuntikkan transaksi ke dalam Context
func InjectTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

// ExtractTx membongkar transaksi dari Context di dalam Repository
func ExtractTx(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(txKey).(*gorm.DB); ok {
		return tx
	}
	return defaultDB
}
