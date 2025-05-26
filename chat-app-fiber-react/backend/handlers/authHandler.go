package handlers

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/HodkiX4/react-go-projects/config"
	"github.com/HodkiX4/react-go-projects/models"
	"github.com/HodkiX4/react-go-projects/repositories"
	"github.com/HodkiX4/react-go-projects/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
SignupHandler handles user registration requests.

It:
  - Parses and validates the request body.
  - Checks for duplicate username or email.
  - Generates a password salt and hashes the password.
  - Inserts the new user document into the MongoDB users collection.

@route POST /api/auth/signup
@returns 201 Created on success, or an appropriate error response.
*/
func SignupHandler(c *fiber.Ctx) error {
	log.Println("Signup attempt")

	// Set a timeout context for DB operations
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var signupRequest models.SignupRequest

	// Parse incoming JSON request to SignupRequest struct
	if err := c.BodyParser(&signupRequest); err != nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate input fields (email, username, password)
	if err := utils.ValidateSignupInput(&signupRequest); err != nil {
		return err
	}

	// Check if the username already exists
	_, err := repositories.GetUserByUsername(ctx, signupRequest.Username)
	if err != mongo.ErrNoDocuments {
		if err != nil {
			return utils.RespondWithError(c, fiber.StatusInternalServerError, "Error checking username")
		}
		return utils.RespondWithError(c, fiber.StatusConflict, "Username already taken")
	}

	// Check if the email is already registered
	_, err = repositories.GetUserByEmail(ctx, signupRequest.Email)
	if err != mongo.ErrNoDocuments {
		if err != nil {
			return utils.RespondWithError(c, fiber.StatusInternalServerError, "Error checking for existing user")
		}
		return utils.RespondWithError(c, fiber.StatusConflict, "Email already in use")
	}

	// Generate a salt and hash the password
	salt, err := utils.GenerateSalt(32)
	if err != nil {
		log.Printf("Error generating salt: %v", err)
		return utils.RespondWithError(c, fiber.StatusInternalServerError, "Internal server error")
	}
	hashedPassword := utils.HashPassword(signupRequest.Password, salt)

	// Create the default avatar URL using UI Avatars
	defaultAvatar := fmt.Sprintf("https://ui-avatars.com/api/?name=%s&background=random", url.QueryEscape(signupRequest.Username))

	// Create the user model instance
	user := models.User{
		Email:     signupRequest.Email,
		Username:  signupRequest.Username,
		Password:  hashedPassword,
		Salt:      base64.StdEncoding.EncodeToString(salt),
		AvatarURL: defaultAvatar,
		CreatedAt: time.Now(),
	}

	// Insert the user into the database
	_, err = config.GetUserCollection().InsertOne(ctx, user)
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		return utils.RespondWithError(c, fiber.StatusInternalServerError, "Could not create user")
	}

	// Respond with success
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Signup successful",
	})
}

/*
LoginHandler handles user login requests.

It:
  - Parses and sanitizes input credentials.
  - Checks if a user with the given email exists.
  - Validates the provided password against the stored hash and salt.
  - If valid, creates and returns a JWT token and sets it as an HTTP-only cookie.

@route POST /api/auth/login
@returns 200 OK with JWT token on success, or an appropriate error response.
*/
func LoginHandler(c *fiber.Ctx) error {
	log.Println("Login attempt")

	// Set a timeout context for DB operations
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var loginRequest models.LoginRequest

	// Parse the login request body
	if err := c.BodyParser(&loginRequest); err != nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Trim and normalize input
	loginRequest.Email = strings.ToLower(strings.TrimSpace(loginRequest.Email))
	loginRequest.Password = strings.TrimSpace(loginRequest.Password)

	// Ensure email and password are not empty
	if loginRequest.Email == "" || loginRequest.Password == "" {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Fill out all fields")
	}

	// Attempt to retrieve the user by email
	existingUser, err := repositories.GetUserByEmail(ctx, loginRequest.Email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return utils.RespondWithError(c, fiber.StatusUnauthorized, "Invalid email or password")
		}
		return utils.RespondWithError(c, fiber.StatusInternalServerError, "Error finding user")
	}

	// Compare provided password with stored hash
	isMatch, err := utils.ComparePassword(loginRequest.Password, existingUser.Salt, existingUser.Password)
	if err != nil {
		return utils.RespondWithError(c, fiber.StatusInternalServerError, "Error comparing password")
	}
	if !isMatch {
		return utils.RespondWithError(c, fiber.StatusUnauthorized, "Invalid email or password")
	}

	// Generate JWT token for authenticated session
	token, err := utils.CreateToken(existingUser.ID.Hex(), existingUser.Email)
	if err != nil {
		log.Printf("Error signing the token %v", err)
		return utils.RespondWithError(c, fiber.StatusInternalServerError, "Could not generate token")
	}

	// Determine if running in production (secure cookie)
	isProduction := os.Getenv("MOD") == "production"

	// Set token as HTTP-only cookie
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		HTTPOnly: true,
		Secure:   isProduction,
		SameSite: "Lax",
	})

	// Return token in response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"token":   token,
	})
}

func LogoutHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(primitive.ObjectID)
	if !ok {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	// Set status to offline in DB
	_, err := config.GetUserCollection().UpdateOne(
		c.Context(),
		bson.M{"_id": userID},
		bson.M{"$set": bson.M{"status": models.StatusOffline}},
	)
	if err != nil {
		return utils.RespondWithError(c, fiber.StatusInternalServerError, "Failed to update user status")
	}

	// Broadcast status update via WebSocket
	/*
		if websocket.WebSocketServer != nil {
			websocket.WebSocketServer.BroadcastStatusChange(userID.Hex(), "offline")
		}
	*/

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User logged out successfully",
	})
}
