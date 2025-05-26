package handlers

import (
	"time"

	"github.com/HodkiX4/react-go-projects/config"
	"github.com/HodkiX4/react-go-projects/models"
	"github.com/HodkiX4/react-go-projects/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MeHandler handles fetching the profile of the currently authenticated user

func MeHandler(c *fiber.Ctx) error {
	// Get the userID from the context (set by the auth middleware)
	userID, ok := c.Locals("userID").(primitive.ObjectID)
	if !ok {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	// Fetch user data from the database
	var user models.User
	err := config.GetUserCollection().FindOne(c.Context(), bson.M{"_id": userID}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return utils.RespondWithError(c, fiber.StatusNotFound, "User not found")
		}
		return utils.RespondWithError(c, fiber.StatusInternalServerError, "Failed to fetch user data")
	}

	// Respond with the user profile data (exclude sensitive info like password)
	user.Password = "" // Exclude password from response
	user.Salt = ""     // Exclude salt from response

	// In the MeHandler, provide default values if fields are empty
	if user.AvatarURL == "" {
		user.AvatarURL = "https://default-avatar-url.com/default.png" // Default avatar
	}
	if user.Friends == nil {
		user.Friends = []primitive.ObjectID{} // Empty array if no friends
	}

	// Set user status to "online"
	user.Status = models.StatusOnline

	// Update in DB
	_, updateErr := config.GetUserCollection().UpdateOne(
		c.Context(),
		bson.M{"_id": userID},
		bson.M{"$set": bson.M{"status": models.StatusOnline}},
	)
	if updateErr != nil {
		return utils.RespondWithError(c, fiber.StatusInternalServerError, "Failed to update status")
	}

	//user.Status = "online"
	// Respond with the profile data
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":        user.ID,
		"username":  user.Username,
		"email":     user.Email,
		"avatarUrl": user.AvatarURL,
		"status":    user.Status,
		"friends":   user.Friends,
	})
}

func GetUsersHandler(c *fiber.Ctx) error {
	// Get the userID from the context (set by the auth middleware)
	userID, ok := c.Locals("userID").(primitive.ObjectID)
	if !ok {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	var dbUsers []models.User = []models.User{}
	// Use Find to retrieve multiple users, excluding the one with the given userID
	cursor, err := config.GetUserCollection().Find(c.Context(), bson.M{"_id": bson.M{"$ne": userID}})
	if err != nil {
		return utils.RespondWithError(c, fiber.StatusInternalServerError, "Failed to fetch users data")
	}
	defer cursor.Close(c.Context())

	if err := cursor.All(c.Context(), &dbUsers); err != nil {
		return utils.RespondWithError(c, fiber.StatusInternalServerError, "Failed to decode users")
	}

	safeUsers := make([]models.SafeUser, 0, len(dbUsers))
	for _, u := range dbUsers {
		friendIDs := make([]string, len(u.Friends))
		for i, fid := range u.Friends {
			friendIDs[i] = fid.Hex()
		}
		safeUsers = append(safeUsers, models.SafeUser{
			ID:        u.ID.Hex(),
			Username:  u.Username,
			Email:     u.Email,
			AvatarURL: u.AvatarURL,
			Status:    string(u.Status),
			Friends:   friendIDs,
			CreatedAt: u.CreatedAt.Format(time.RFC3339),
		})
	}

	// Return the list of users
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"users": safeUsers,
	})
}
