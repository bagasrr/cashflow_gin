package config

import (
	"cashflow_gin/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabaseConnection() (*gorm.DB, error) {
	// 1. LANGSUNG PANGGIL AppConfig. Gak perlu import package config lagi.
	host := AppConfig.Database.Host
	user := AppConfig.Database.User
	password := AppConfig.Database.Password
	dbName := AppConfig.Database.Name
	port := AppConfig.Database.Port // Ini tipe datanya int

	// 2. Perhatikan %d untuk port
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta",
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

	// Automigrate
	err = con.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Group{},
		&models.GroupMember{},
		&models.Wallet{},
		&models.Transaction{},
	)
	if err != nil {
		log.Println("Gagal AutoMigrate:", err)
		return nil, err
	}

	return con, nil
}
