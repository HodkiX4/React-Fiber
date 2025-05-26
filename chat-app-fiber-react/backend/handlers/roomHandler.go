package handlers

import (
	"log"
	"strings"
	"time"

	"github.com/HodkiX4/react-go-projects/config"
	"github.com/HodkiX4/react-go-projects/models"
	"github.com/HodkiX4/react-go-projects/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetChatRoomsHandler handles fetching all the chat rooms for a user
func GetChatRoomsHandler(c *fiber.Ctx) error {
	// Get the user ID from the context (set by the AuthMiddleware)
	userID := c.Locals("userID").(primitive.ObjectID)
	log.Println("user ID:", userID)

	//Find rooms where the user is a member OR the room is public
	filter := bson.M{
		"$or": []bson.M{
			{"members": userID},
			{"isPrivate": false},
		},
	}

	// Fetch all chat rooms the user is part of
	cursor, err := config.GetRoomCollection().Find(c.Context(), filter)
	if err != nil {
		return utils.RespondWithError(c, fiber.StatusInternalServerError, "Failed to fetch chat rooms")
	}
	defer cursor.Close(c.Context())

	var rooms []models.Room = []models.Room{}

	for cursor.Next(c.Context()) {
		var room models.Room
		if err := cursor.Decode(&room); err != nil {
			return utils.RespondWithError(c, fiber.StatusInternalServerError, "Failed to decode chat room data")
		}
		rooms = append(rooms, room)
	}

	// Check for cursor iteration errors
	if err := cursor.Err(); err != nil {
		return utils.RespondWithError(c, fiber.StatusInternalServerError, "Error iterating through chat rooms")
	}

	// Respond with the list of chat rooms
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Chat rooms fetched successfully",
		"rooms":   rooms,
	})
}

func CreateChatRoomHandler(c *fiber.Ctx) error {

	var room models.Room

	// Parse incoming request body to populate the room struct
	if err := c.BodyParser(&room); err != nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Creators ID
	creatorID := c.Locals("userID").(primitive.ObjectID)

	// Add creator to members
	room.Members = append(room.Members, creatorID)

	// Assign the creator as the first admin
	room.Admins = append(room.Admins, creatorID)

	// Assign the creator as a moderator
	room.Moderators = append(room.Moderators, creatorID)

	// Validate required fields
	if room.Name == "" || len(room.Members) == 0 {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Room name and members are required")
	}

	// Enforce max name length
	if len(room.Name) > 50 {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Room name is to long (max 50 characters)")
	}

	// Make room name lowercase
	room.Name = strings.ToLower(room.Name)

	// Validate member ObjectIDs
	validatedMembers := make([]primitive.ObjectID, 0, len(room.Members))
	for _, member := range room.Members {
		objID, err := primitive.ObjectIDFromHex(member.Hex())
		if err != nil {
			return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid member ID")
		}
		validatedMembers = append(validatedMembers, objID)
	}
	room.Members = validatedMembers

	// Set the createdAt time to the current time
	room.CreatedAt = time.Now()

	// Insert room into the database
	result, err := config.GetRoomCollection().InsertOne(c.Context(), room)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return utils.RespondWithError(c, fiber.StatusConflict, "Chat room already exists")
		}
		return utils.RespondWithError(c, fiber.StatusInternalServerError, "Failed to create chat room")
	}

	room.ID = result.InsertedID.(primitive.ObjectID)
	// Respond with success
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Chat room created successfully",
		"room":    room,
	})

}

func DeleteChatRoomHandler(c *fiber.Ctx) error {
	roomIDParam := c.Params("roomId")
	roomID, err := primitive.ObjectIDFromHex(roomIDParam)
	if err != nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid room ID")
	}

	userID := c.Locals("userID").(primitive.ObjectID)

	var room models.Room
	roomCollection := config.GetRoomCollection()
	err = roomCollection.FindOne(c.Context(), bson.M{
		"_id": roomID,
	}).Decode(&room)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return utils.RespondWithError(c, fiber.StatusNotFound, "Chat room not found")
		}
		return utils.RespondWithError(c, fiber.StatusInternalServerError, "Failed to fetch room")
	}

	isAdmin := false
	for _, adminID := range room.Admins {
		if adminID == userID {
			isAdmin = true
		}
	}

	if !isAdmin {
		return utils.RespondWithError(c, fiber.StatusUnauthorized, "You are not authorized to delete this room")
	}

	_, err = roomCollection.DeleteOne(c.Context(), bson.M{
		"_id": roomID,
	})
	if err != nil {
		return utils.RespondWithError(c, fiber.StatusInternalServerError, "Failed to delete chat room")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Chat room deleted successfully",
	})
}
