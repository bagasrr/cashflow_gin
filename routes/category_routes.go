package routes

import (
	"cashflow_gin/controllers"
	"cashflow_gin/middlewares"

	"github.com/gin-gonic/gin"
)

func CategoryRoutes(r *gin.RouterGroup, controller *controllers.CategoryController) {
	categories := r.Group("/categories")
	categories.Use(middlewares.AuthMiddleware())
	{
		categories.POST("/default-cat-admin-only-wlee", controller.CreateDefaultCategories)
		categories.GET("/", controller.GetAllCategories)

		categories.POST("/mine", controller.CreateMy)
		categories.GET("/mine", controller.GetMine)
		categories.PATCH("/:id/update", controller.UpdateById)
		categories.PATCH("/:id/delete", controller.DeleteById)
	}
}
