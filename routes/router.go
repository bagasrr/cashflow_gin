package routes

import (
	"cashflow_gin/controllers"
	"cashflow_gin/repository"
	"cashflow_gin/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB, r *gin.Engine) {
	// 1. INIT REPOSITORIES (Layer Paling Bawah)
	// Cukup sekali init per repo!
	userRepo := repository.NewUserRepository(db)
	transRepo := repository.NewTransactionRepository(db)
	catRepo := repository.NewCategoryRepository(db) // <--- Dipake bareng-bareng
	authRepo := repository.NewAuthRepository(db)
	groupRepo := repository.NewGroupRepository(db)   // <--- Repo baru untuk Group
	walletRepo := repository.NewWalletRepository(db) // <--- Repo baru untuk Wallet

	// 2. INIT SERVICES (Layer Tengah)
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(authRepo)
	catService := services.NewCategoryService(catRepo)
	groupService := services.NewGroupService(groupRepo)

	// Perhatikan ini: TransactionService butuh catRepo & userRepo juga
	// Karena kita udah init di atas, tinggal masukin variabelnya.
	transService := services.NewTransactionService(transRepo, catRepo, userRepo, groupRepo, walletRepo)
	walletService := services.NewWalletService(walletRepo, groupRepo) // Service untuk Wallet, kalau nanti butuh logic khusus selain repo langsung bisa ditambahin di sini

	// 3. INIT CONTROLLERS (Layer Atas)
	userController := controllers.NewUserController(userService)
	authController := controllers.NewAuthController(authService)
	catController := controllers.NewCategoryController(catService)
	transController := controllers.NewTransactionController(transService)
	groupController := controllers.NewGroupController(groupService)
	walletController := controllers.NewWalletController(walletService) // Controller untuk Wallet

	// 4. ROUTING GROUP (Panggil file-file routes yang udah dipisah)
	api := r.Group("/api")
	{
		// Lempar Controller yang udah jadi ke masing-masing file route
		AuthRoutes(api, authController)
		UserRoutes(api, userController)
		CategoryRoutes(api, catController)
		TransactionRoutes(api, transController)
		GroupRoutes(api, groupController)
		WalletRoutes(api, walletController)
	}
}
