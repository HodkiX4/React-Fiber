package routes

import (
	"github.com/HodkiX4/react-go-projects/handlers"
	"github.com/HodkiX4/react-go-projects/middlewares"
	"github.com/gofiber/fiber/v2"
)

func ChatRoutes(route fiber.Router) {
	chat := route.Group("/chat", middlewares.AuthMiddleware)

	chat.Post("/message", handlers.SendMessageHandler)
	chat.Delete("/message/:messageId", handlers.DeleteMessageHandler)
	chat.Get("/message/:roomId", handlers.GetMessagesHandler)

	chat.Post("/rooms", handlers.CreateChatRoomHandler)
	chat.Delete("/rooms/:roomId", handlers.DeleteChatRoomHandler)
	chat.Get("/rooms", handlers.GetChatRoomsHandler)
}
