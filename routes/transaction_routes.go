package routes

import (
	"cashflow_gin/controllers"
	"cashflow_gin/middlewares"

	"github.com/gin-gonic/gin"
)

func TransactionRoutes(r *gin.RouterGroup, controller *controllers.TransactionController) {
	transactions := r.Group("/transactions")
	transactions.Use(middlewares.AuthMiddleware()) // Middleware dipasang di sini
	{
		transactions.POST("/", controller.Create)
		transactions.GET("/", controller.FindAll)
		transactions.GET("/:id/detail", controller.GetTransactionByID)
		transactions.PATCH("/:id/update", controller.UpdateTransaction)
		transactions.PATCH("/:id/wallet/:walletid/soft-delete", controller.SoftDeleteTransaction)
	}
}
