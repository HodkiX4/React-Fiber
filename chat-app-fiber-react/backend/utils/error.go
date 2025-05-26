package utils

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

// Function used to handle errors on the backend
func RespondWithError(c *fiber.Ctx, status int, message string) error {
	// Log the error with the status code
	log.Printf("Error: %s, Status: %d", message, status)

	// Return error response to the client
	if c != nil {
		return c.Status(status).JSON(fiber.Map{
			"error": message,
		})
	}

	// If `c` is nil (non-request context), return a fiber error
	return fiber.NewError(status, message)
}
