package routes

import (
	"cashflow_gin/controllers"
	"cashflow_gin/middlewares"

	"github.com/gin-gonic/gin"
)

func GroupRoutes(r *gin.RouterGroup, controller *controllers.GroupController) {
	groups := r.Group("/groups")
	groups.Use(middlewares.AuthMiddleware()) // Middleware dipasang di sini
	{
		groups.GET("/", controller.GetAllGroups)
		groups.POST("/", controller.CreateGroup)
		groups.GET("/:id", controller.GetGroupByID)
		// groups.PATCH("/:id/update", controller.UpdateGroup)
		groups.PATCH("/:id/remove-user", controller.RemoveUserFromGroup)
	}
}
