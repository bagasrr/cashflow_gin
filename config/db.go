package config

import (
	"cashflow_gin/models"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabaseConnection() (*gorm.DB, error) {
	// Ambil data dari .env
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		host, user, password, dbName, port)

	con, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := con.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	// Automigrate (opsional, tapi buat dev enak)
	err = con.AutoMigrate(
		&models.User{},
		&models.Category{}, // Jangan lupa Category juga harus ada
		&models.Group{},
		&models.GroupMember{},
		&models.Wallet{},
		&models.Transaction{},
	)
	if err != nil {
		fmt.Println("Gagal AutoMigrate:", err)
		return nil, err
	}

	return con, nil
}
