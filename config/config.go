package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Bikin hirarki struct sesuai bentuk YAML
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port int
	Mode string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

// Global variable penampung config
var AppConfig Config

func LoadConfig() {
	// 1. Baca .env dulu buat tau kita di environment mana (local/staging/prod)
	_ = godotenv.Load() 
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "local" // Default kalau .env gak ada
	}

	// 2. Kasih tau Viper file YAML mana yang harus dibaca
	viper.SetConfigName("config-" + env) // contoh: akan mencari file 'config-local'
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".") // Cari di folder root project

	// 3. Baca filenya
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Fatal error config file: %v", err)
	}

	// 4. Masukkan isi YAML ke dalam Struct AppConfig lu
	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("Gagal unmarshal config: %v", err)
	}

	log.Printf("✅ Config berhasil diload untuk environment: %s", env)
}