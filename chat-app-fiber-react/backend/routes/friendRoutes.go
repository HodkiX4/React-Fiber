package routes

import (
	"github.com/HodkiX4/react-go-projects/handlers"
	"github.com/HodkiX4/react-go-projects/middlewares"
	"github.com/gofiber/fiber/v2"
)

func FriendRoutes(router fiber.Router) {
	friend := router.Group("/friend-request", middlewares.AuthMiddleware)

	friend.Post("/", handlers.SendFriendRequest)
	friend.Patch("/accept/:id", handlers.AcceptFriendRequest)
	friend.Patch("/reject/:id", handlers.RejectFriendRequest)
	friend.Get("/pending", handlers.GetPendingFriendRequests)
}
