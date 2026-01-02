package config

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var ctx = context.Background()

func InitRedis(){
	// Menghubungkan ke Redis yang baru diinstall
	RedisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,	
	})

	// Tes Koneksi
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil{
		panic(fmt.Sprintf("Could not connect to redis: %v", err))
	}
	fmt.Println("Connected to Redis Successfully")
}