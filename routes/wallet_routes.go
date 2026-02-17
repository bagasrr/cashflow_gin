package routes

import (
	"cashflow_gin/controllers"
	"cashflow_gin/middlewares"

	"github.com/gin-gonic/gin"
)

func WalletRoutes(r *gin.RouterGroup, controller *controllers.WalletController) {
	wallets := r.Group("/wallets")
	wallets.Use(middlewares.AuthMiddleware()) // Middleware dipasang di sini
	{
		wallets.GET("/", controller.GetAllWallets)
		wallets.GET("/:id/detail", controller.GetWalletByID)
	}
}
