package main

import (
	"cashflow_gin/config"
	_ "cashflow_gin/docs"
	"cashflow_gin/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Cashflow API Gin
// @version         1.0
// @description     API untuk manajemen keuangan pribadi (Cashflow).
// @termsOfService  http://swagger.io/terms/

// @contact.name   Bagas Rr
// @contact.url    http://bagasrr.my.id

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @host      localhost:8080
// @BasePath  /api
func main(){
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	db, err := config.NewDatabaseConnection()
	if err != nil {
		log.Fatal("Gagal Konek Database, err")
	}
	


	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	

	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	routes.SetupRoutes(db, r)

	r.Run(os.Getenv("APP_PORT"))
}