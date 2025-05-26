package routes

import (
	"github.com/HodkiX4/react-go-projects/handlers"
	"github.com/HodkiX4/react-go-projects/middlewares"
	"github.com/gofiber/fiber/v2"
)

func UserRoutes(router fiber.Router) {
	user := router.Group("/user", middlewares.AuthMiddleware)

	user.Get("/me", handlers.MeHandler)
	user.Get("/all", handlers.GetUsersHandler)
}
