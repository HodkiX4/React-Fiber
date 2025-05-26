package middlewares

import (
	"fmt"
	"log"
	"os"

	"github.com/HodkiX4/react-go-projects/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuthMiddleware is used to protect routes that require authentication.
// It checks for a valid JWT token in the cookies and verifies it.

func AuthMiddleware(c *fiber.Ctx) error {
	// Get token from cookies
	tokenString := c.Cookies("token")

	// If no token is provided, return Unauthorized
	if tokenString == "" {
		return utils.RespondWithError(c, fiber.StatusUnauthorized, "Missing or invalid token")
	}

	// Parse the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// Validate the algorithm used in the token (should be HMAC)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret key used for verification
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	// If parsing fails, return Unauthorized
	if err != nil {
		log.Printf("Error parsing token: %v", err)
		// Check if the error is due to expiration
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return utils.RespondWithError(c, fiber.StatusUnauthorized, "Token expired")
			}
		}
		return utils.RespondWithError(c, fiber.StatusUnauthorized, "Invalid token")
	}

	// Check if the token is valid and extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Extract user ID from claims
		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			return utils.RespondWithError(c, fiber.StatusUnauthorized, "Invalid token claims")
		}

		// Convert user ID to primitive.ObjectID
		userID, err := primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			return utils.RespondWithError(c, fiber.StatusUnauthorized, "Invalid user ID in token")
		}

		// Store user ID in Fiber's context locals to access in handlers
		c.Locals("userID", userID)

		// Token is valid, continue to the next middleware/handler
		return c.Next()
	}

	// If token is invalid, respond with Unauthorized
	return utils.RespondWithError(c, fiber.StatusUnauthorized, "Invalid token")
}
