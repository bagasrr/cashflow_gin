package server

import (
	"cashflow_gin/api"
	"cashflow_gin/config"
	"cashflow_gin/handler"
	"cashflow_gin/middlewares"
	"cashflow_gin/repository"
	"cashflow_gin/services"
	"fmt"
	"log"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

// 1. MASTER API: Kunci untuk melunasi kontrak StrictServerInterface
type MasterAPI struct {
	*handler.AuthAPI
	*handler.CategoryAPI
	*handler.TransactionAPI
	*handler.UserAPI
	*handler.WalletAPI
	*handler.GroupAPI
	*handler.DashboardAPI
}

type standardError struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Errors  string `json:"errors"`
}

func Run() {
	config.LoadConfig()

	db, err := config.NewDatabaseConnection()
	if err != nil {
		log.Fatal("Gagal Konek Database: ", err)
	}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // URL Next.js lu
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	// ------------------------------------
	// ZONA SETUP DESIGN-FIRST (MASA DEPAN)
	// ------------------------------------

	// A. Inisialisasi Repositories (Layer Bawah)
	authRepo := repository.NewAuthRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	userRepo := repository.NewUserRepository(db)
	walletRepo := repository.NewWalletRepository(db)
	groupRepo := repository.NewGroupRepository(db) // Buka kalau group udah ada
	dashboardRepo := repository.NewDashboardRepository(db)

	// B. Inisialisasi Services (Layer Tengah)
	authService := services.NewAuthService(authRepo)
	categoryService := services.NewCategoryService(categoryRepo, groupRepo, userRepo)
	userService := services.NewUserService(userRepo)
	// Hati-hati: TransactionService lu biasanya butuh banyak repo
	transactionService := services.NewTransactionService(transactionRepo, categoryRepo, userRepo, groupRepo, walletRepo)
	walletService := services.NewWalletService(walletRepo, nil)
	groupService := services.NewGroupService(groupRepo)
	dashboardService := services.NewDashboardService(dashboardRepo)

	// C. Inisialisasi Handlers (Gerbang Luar)
	authAPI := &handler.AuthAPI{Service: authService}
	categoryAPI := &handler.CategoryAPI{Service: categoryService}
	transactionAPI := &handler.TransactionAPI{Service: transactionService}
	userAPI := &handler.UserAPI{Service: userService}
	walletAPI := &handler.WalletAPI{Service: walletService}
	groupAPI := &handler.GroupAPI{Service: groupService}
	dashboardAPI := &handler.DashboardAPI{Service: dashboardService}

	// D. Gabungkan ke Master API
	masterHandler := &MasterAPI{
		AuthAPI:        authAPI,
		CategoryAPI:    categoryAPI,
		TransactionAPI: transactionAPI,
		UserAPI:        userAPI,
		WalletAPI:      walletAPI,
		GroupAPI:       groupAPI,
		DashboardAPI:   dashboardAPI,
	}

	r.Use(middlewares.AuthMiddleware())
	// E. Daftarkan ke Gin Router
	strictHandler := api.NewStrictHandler(masterHandler, nil)
	// PENTING: Jangan lupa kasih BaseURL "/api" kalau di yaml lu path-nya mulai dari "/api"
	// Tapi kalau di yaml lu udah nulis "/api/categories", gak usah pake WithBaseURL di sini.

	api.RegisterHandlersWithOptions(r, strictHandler, api.GinServerOptions{
		BaseURL: "/api",
		ErrorHandler: func(c *gin.Context, err error, statusCode int) {
			c.JSON(statusCode, standardError{
				Status:  false,
				Message: "Invalid Request Format",
				Errors:  err.Error(),
			})
		},
	})

	// ------------------------------------
	// ZONA SETUP LAMA (PERINGATAN: BOM WAKTU)
	// ------------------------------------
	// Gw biarin ini nyala, TAPI lu harus segera mematikan rute-rute di router.go
	// yang udah lu pindahin ke OpenAPI, biar Gin gak Panic karena rute bentrok.
	// routes.SetupRoutes(db, r)

	// ------------------------------------
	// ZONA SWAGGER DOCS
	// ------------------------------------
	r.StaticFile("/openapi.yaml", "./api/openapi.yaml")
	r.GET("/docs", func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		c.String(200, swaggerHTML)
	})

	// ------------------------------------
	// ZONA JALANKAN SERVER
	// ------------------------------------
	port := config.AppConfig.Server.Port
	if port == 0 {
		port = 8080
	}
	portStr := fmt.Sprintf(":%d", port)

	log.Printf("🚀 Server jalan di port %s", portStr)
	r.Run(portStr)
}

var swaggerHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>Cashflow API Docs</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui.css" />
</head>
<body style="margin:0;">
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-bundle.js" crossorigin></script>
  <script>
    window.onload = () => {
      window.ui = SwaggerUIBundle({
        url: '/openapi.yaml',
        dom_id: '#swagger-ui',
      });
    };
  </script>
</body>
</html>`
