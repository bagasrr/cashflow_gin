package routes

import (
	"cashflow_gin/controllers"
	"cashflow_gin/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.RouterGroup, controller *controllers.UserController) {
	users := r.Group("/users")
	users.Use(middlewares.AuthMiddleware())
	{
		users.GET("/", controller.FindAllUser)
		users.GET("/me", controller.GetMyProfile)
	}
}
