package repositories

import (
	"context"

	"github.com/HodkiX4/react-go-projects/config"
	"github.com/HodkiX4/react-go-projects/models"
	"go.mongodb.org/mongo-driver/bson"
)

// GetUserByField queries the database by a specific field (email or username) and returns the user.
func GetUserByField(ctx context.Context, field string, value any) (models.User, error) {
	var existingUser models.User
	filter := bson.M{field: value}
	err := config.GetUserCollection().FindOne(ctx, filter).Decode(&existingUser)
	return existingUser, err
}

// Queries the database by email field using the GetUserByField function and returns the user
func GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	return GetUserByField(ctx, "email", email)
}

// Queries the database by username field using the GetUserByField function and returns the user
func GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	return GetUserByField(ctx, "username", username)
}
