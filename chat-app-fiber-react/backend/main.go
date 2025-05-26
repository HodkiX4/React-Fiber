package main

import (
	"log"
	"os"

	"github.com/HodkiX4/react-go-projects/config"
	"github.com/HodkiX4/react-go-projects/routes"
	"github.com/HodkiX4/react-go-projects/ws"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/websocket/v2"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	// Initialize a new Fiber app
	app := fiber.New()

	// Middleware for CORS and logging
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, PATCH, DELETE, OPTIONS",
	}))
	app.Use(logger.New())

	// Connect to MongoDB
	config.ConnectDB()

	// Register your regular routes (e.g., API routes)
	routes.RegisterRoutes(app)

	// WebSocket route using the new HandleWebSocket function
	app.Get("/ws", websocket.New(ws.HandleWebSocket))

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	// Start API server (Fiber) in a goroutine
	go func() {
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("Fiber server failed: %v", err)
		}
	}()

	// Block forever so the server can run
	select {}
}
