package config

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func InitRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Sesuai dengan Docker lu
		Password: "",               // Kosongin kalau local
		DB:       0,                // Database default Redis
	})

	// Uji koneksi mutlak
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Bencana: Gagal konek ke Redis! %v", err)
	}

	return rdb
}
