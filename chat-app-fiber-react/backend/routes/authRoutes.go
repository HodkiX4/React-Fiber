package routes

import (
	"github.com/HodkiX4/react-go-projects/handlers"
	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(router fiber.Router) {
	auth := router.Group("/auth")

	auth.Post("/register", handlers.SignupHandler)
	auth.Post("/login", handlers.LoginHandler)

}
