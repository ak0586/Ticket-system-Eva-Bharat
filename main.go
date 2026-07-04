package main

import (
	"log"
	"os"

	"ticket-system/db"
	"ticket-system/handlers"
	"ticket-system/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	db.Init()

	router := gin.Default()

	router.GET("/health", handlers.HealthCheck)

	auth := router.Group("/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}

	tickets := router.Group("/tickets")
	tickets.Use(middleware.Auth())
	{
		tickets.POST("", handlers.CreateTicket)
		tickets.GET("", handlers.ListTickets)
		tickets.GET("/:id", handlers.GetTicket)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("starting on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
