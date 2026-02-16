package routes

import (
	"cashflow_gin/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.RouterGroup, controller *controllers.UserController) {
	users := r.Group("/users")
	{
		users.GET("/", controller.FindAllUser)
	}
}
