package routes

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

type LogPayload struct {
	Message   string `json:"message"`
	TimeStamp string `json:"timeStamp"`
}

func LoggerRoutes(route fiber.Router) {
	route.Post("/log-error", func(c *fiber.Ctx) error {
		var payload LogPayload

		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid payload",
			})
		}

		// Print or store log as needed (e.g., write to file, DB, etc.)
		log.Printf("[ERROR LOG] %s - %s\n", payload.TimeStamp, payload.Message)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Error logged successfully",
		})
	})
}
