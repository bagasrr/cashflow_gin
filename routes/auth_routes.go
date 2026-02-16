package routes

import (
	"cashflow_gin/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.RouterGroup, controller *controllers.AuthController) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", controller.Register)
		auth.POST("/login", controller.Login)
	}
}
