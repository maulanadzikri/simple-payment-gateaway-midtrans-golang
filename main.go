package main

import (
	"fmt"
	"log"

	"github.com/bagussubagja/backend-payment-gateway-go/api/routes"
	"github.com/bagussubagja/backend-payment-gateway-go/config"
	repository "github.com/bagussubagja/backend-payment-gateway-go/internal/repositories"
	"github.com/bagussubagja/backend-payment-gateway-go/internal/services"
	"github.com/bagussubagja/backend-payment-gateway-go/storage"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	db, err := storage.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("could not connect to db: %v", err)
	}
	fmt.Println("Database connected successfully")

	config.InitRedis()	

	userRepo := repository.NewUserRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	authService := services.NewAuthService(userRepo, cfg)
	userService := services.NewUserService(userRepo)
	midtransService := services.NewMidtransService(cfg)
	paymentService := services.NewPaymentService(transactionRepo, midtransService)

	router := routes.SetupRouter(authService, userService, paymentService, cfg)

	serverAddress := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server is running on port %s", cfg.ServerPort)
	if err := router.Run(serverAddress); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
