package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	// Ganti import config ini sesuai path project Cashflow lu
	"cashflow_gin/config"
	"cashflow_gin/dto/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

func AuthMiddleware(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		publicRoutes := []string{
			"/docs", // Biar Swagger UI tetep bisa dibuka
			"/api/auth/login",
			"/api/auth/register",
			"/api/auth/forgot-password",
			"/openapi.yaml", // Biar file dokumentasinya tetep bisa dibaca
		}

		// 2. CEK APAKAH USER MENUJU JALUR VIP
		currentPath := c.Request.URL.Path
		for _, route := range publicRoutes {

			// Kalau URL awalan-nya cocok dengan whitelist, langsung loloskan
			if strings.HasPrefix(currentPath, route) {
				c.Next()
				return // Hentikan pengecekan token di bawahnya
			}
		}
		// PASANG RADAR INI SEMENTARA BUAT DEBUGGING
		fmt.Println("\n=== DEBUG MIDDLEWARE ===")
		fmt.Println("Path yang ditembak Postman :", currentPath)
		//fmt.Println("Apakah ada di Whitelist?   :", publicRoutes[2])
		fmt.Println("========================\n")

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.BaseResponse{
				Status:  false,
				Message: "Unauthorized",
				Errors:  "Missing or invalid Authorization header",
			})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 1. TANGKAP ERROR-NYA, JANGAN DIBUANG
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Pastikan metode signing-nya sesuai
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			// 2. AMBIL DARI CONFIG MEMORI, BUKAN OS.GETENV
			// Asumsi lu punya JWTSecret di struct config lu
			return []byte(config.AppConfig.JWTSecret), nil
		})
		fmt.Println("\nToken Valid : " + tokenString)

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.BaseResponse{
				Status:  false,
				Message: "Unauthorized",
				Errors:  "Token is expired or invalid",
			})
			return
		}

		ctx := c.Request.Context()
		redisErr := rdb.Get(ctx, tokenString).Err()
		if redisErr != nil && redisErr != redis.Nil {
			log.Printf("CRITICAL: Redis Middleware Error: %v", redisErr)
			c.AbortWithStatusJSON(http.StatusInternalServerError, response.BaseResponse{
				Status:  false,
				Message: "Internal Server Error",
				Errors:  "Sistem otentikasi sedang mengalami gangguan",
			})
			return
		}

		// Jika tidak ada error sama sekali (berarti data token KETEMU di daftar hitam)
		if redisErr == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.BaseResponse{
				Status:  false,
				Message: "Unauthorized",
				Errors:  "Akses ditolak: Token sudah tidak berlaku (Revoked)",
			})
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.BaseResponse{
				Status:  false,
				Message: "Unauthorized",
				Errors:  "Invalid token structure",
			})
			return
		}
		expClaim, expOk := claims["exp"].(float64)
		if !expOk {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.BaseResponse{
				Status: false, Message: "Unauthorized", Errors: "Token tidak memiliki expiration time",
			})
			return
		}

		// 3. SAFE TYPE ASSERTION (ANTI-PANIC)
		userID, idOk := claims["user_id"].(string)
		// userRole, roleOk := claims["user_role"].(string)
		rawRole := claims["user_role"]
		roleOk := false
		var userRole string

		// Normalisasi tipe data: Apapun bentuknya (float64, int, string), jadikan string
		if rawRole != nil {
			userRole = fmt.Sprintf("%v", rawRole)
			log.Printf("DEBUG: Extracted user_role claim: %v (type %T) -> normalized to string: %s", rawRole, rawRole, userRole)
			roleOk = true
		}

		// Validasi akhir
		if !idOk || !roleOk {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.BaseResponse{
				Status:  false,
				Message: "Unauthorized",
				Errors:  "Token payload is missing required claims or format is invalid",
			})
			return
		}
		c.Set("token_string", tokenString)
		c.Set("token_exp", expClaim)
		c.Set("user_id", userID)
		c.Set("user_role", userRole)

		c.Next()
	}
}
