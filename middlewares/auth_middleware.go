package middlewares

import (
	"cashflow_gin/dto/response"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if !strings.Contains(authHeader, "Bearer") {
            c.AbortWithStatusJSON(http.StatusUnauthorized, response.BaseResponse{
                Status:  false,
                Message: "Unauthorized: No Bearer token | Please login first",
                Errors:  "Missing or invalid Authorization header",
                Data:    nil,
            })
            return
        }

        tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
        
        // token, err := jwt.Parse(tokenString, func(t *jwt.SimpleClaims) (interface{}, error) {
		// 	return []byte(os.Getenv("JWT_SECRET")), nil
        // })

		token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    		return []byte(os.Getenv("JWT_SECRET")), nil
		})

        if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
            // Simpan UserID ke context biar bisa dipake Controller
            c.Set("user_id", claims["user_id"])
            c.Set("user_role", claims["user_role"])
            c.Next()
        } else {
            c.AbortWithStatusJSON(http.StatusUnauthorized, response.BaseResponse{
                Status:  false,
                Message: "Unauthorized: Invalid token | Please login again",
                Errors:  "Token Expired or Invalid",
                Data:    nil,
            })
            return
        }
    }
}