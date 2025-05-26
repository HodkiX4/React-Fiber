package routes

import "github.com/gofiber/fiber/v2"

// RegisterRoutes sets up all route groups under /api
func RegisterRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Mount each route group
	AuthRoutes(api)
	UserRoutes(api)
	ChatRoutes(api)
	FriendRoutes(api)
	LoggerRoutes(api)
	// ChatRoutes(api) // (add later)
}
