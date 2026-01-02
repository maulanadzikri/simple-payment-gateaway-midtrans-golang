package routes

import (
	"time"

	"github.com/bagussubagja/backend-payment-gateway-go/api/handler"
	"github.com/bagussubagja/backend-payment-gateway-go/api/middleware"
	"github.com/bagussubagja/backend-payment-gateway-go/config"
	"github.com/bagussubagja/backend-payment-gateway-go/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(authSvc services.AuthService, userSvc services.UserService, paymentSvc services.PaymentService, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000", 
			// "https://525a355aecd6.ngrok-free.app",	
		},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type","Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders: []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	}))

	entryHandler := handler.NewEntryHandler()
	authHandler := handler.NewAuthHandler(authSvc)
	userHandler := handler.NewUserHandler(userSvc)
	paymentHandler := handler.NewPaymentHandler(paymentSvc, userSvc)

	r.GET("/", entryHandler.GetEntry)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	apiV1 := r.Group("/api/v1")
	{
		auth := apiV1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", middleware.AuthMiddleware(authSvc), authHandler.Logout)
		}

		apiV1.POST("/payments/notification", paymentHandler.HandleNotification)
	}

	authorized := apiV1.Group("/")
	authorized.Use(middleware.AuthMiddleware(authSvc))
	{
		authorized.GET("/profile", userHandler.GetProfile)
		payments := authorized.Group("/payments")
		{
			payments.POST("/create", paymentHandler.CreatePayment)
			payments.GET("/status/:orderID", paymentHandler.GetStatus)
			payments.GET("/history", paymentHandler.GetHistory)
			payments.POST("/qris", paymentHandler.CreateQrisPayment)
		}

	}

	return r
}
