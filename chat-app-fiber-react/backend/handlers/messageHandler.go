package handlers

import (
	"log"
	"slices"
	"time"

	"github.com/HodkiX4/react-go-projects/config"
	"github.com/HodkiX4/react-go-projects/models"
	"github.com/HodkiX4/react-go-projects/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SendMessageHandler handles incoming POST requests to send a message to a chat room.
// It validates the input, parses the sender ID from the context, and stores the message in MongoDB.
func SendMessageHandler(c *fiber.Ctx) error {
	var payload struct {
		RoomID  string `json:"roomId"`  // ID of the chat room
		Content string `json:"content"` // Message content
	}

	// Parse and validate the JSON body
	if err := c.BodyParser(&payload); err != nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if payload.RoomID == "" || payload.Content == "" {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Room ID and content are required")
	}

	// Convert room ID from string to MongoDB ObjectID
	roomID, err := primitive.ObjectIDFromHex(payload.RoomID)
	if err != nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid Room ID")
	}

	// Extract sender ID from the request context
	senderID := c.Locals("userID").(primitive.ObjectID)

	// Create the message object
	message := models.Message{
		RoomID:    roomID,
		SenderID:  senderID,
		Content:   payload.Content,
		CreatedAt: time.Now(),
	}

	// Insert the message into the database
	res, err := config.GetMessageCollection().InsertOne(c.Context(), message)
	if err != nil {
		return utils.RespondWithError(c, fiber.StatusInternalServerError, "Failed to send message")
	}

	// Set the inserted ID on the message object
	message.ID = res.InsertedID.(primitive.ObjectID)

	// Return the created message
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": message,
	})
}

func DeleteMessageHandler(c *fiber.Ctx) error {
	log.Println("Inside")
	messageIdStr := c.Params("messageId")
	userID := c.Locals("userID").(primitive.ObjectID)

	// Convert messageId to ObjectID
	messageId, err := primitive.ObjectIDFromHex(messageIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid message ID format",
		})
	}

	// Fetch message
	var message models.Message
	messageCollection := config.GetMessageCollection()
	err = messageCollection.FindOne(c.Context(), bson.M{"_id": messageId}).Decode(&message)
	if err != nil {
		return utils.RespondWithError(c, fiber.StatusNotFound, "Message not found")
	}

	// Fetch room
	var room models.Room
	roomCollection := config.GetRoomCollection()
	err = roomCollection.FindOne(c.Context(), bson.M{"_id": message.RoomID}).Decode(&room)
	if err != nil {
		return utils.RespondWithError(c, fiber.StatusNotFound, "Room not found")
	}
	// Check permissons
	isOwner := message.SenderID == userID
	isAdminOrMod := slices.Contains(room.Admins, userID) || slices.Contains(room.Moderators, userID)

	if !isOwner && !isAdminOrMod {
		return utils.RespondWithError(c, fiber.StatusForbidden, "Not authorized to delete this message")
	}

	_, err = messageCollection.DeleteOne(c.Context(), bson.M{"_id": messageId})
	if err != nil {
		return utils.RespondWithError(c, fiber.StatusInternalServerError, "Failed to delete")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Message deleted",
	})
}

// GetMessagesHandler handles GET requests to fetch the latest 50 messages from a specific chat room.
// Messages are sorted in reverse chronological order.
func GetMessagesHandler(c *fiber.Ctx) error {
	// Parse and validate the room ID from URL parameters
	roomIDStr := c.Params("roomId")
	roomID, err := primitive.ObjectIDFromHex(roomIDStr)
	if err != nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid Room ID")
	}

	// Build query filter and options
	filter := bson.M{"roomId": roomID}
	findOptions := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}}). // Sort by createdAt descending
		SetLimit(50)                                    // Limit to 50 messages

	// Query the database for messages
	cursor, err := config.GetMessageCollection().Find(c.Context(), filter, findOptions)
	if err != nil {
		return utils.RespondWithError(c, fiber.StatusInternalServerError, "Failed to fetch messages")
	}
	defer cursor.Close(c.Context())

	// Decode the results
	var messages []models.Message
	for cursor.Next(c.Context()) {
		var msg models.Message
		if err := cursor.Decode(&msg); err != nil {
			return utils.RespondWithError(c, fiber.StatusInternalServerError, "Failed to decode message")
		}
		messages = append(messages, msg)
	}

	// Return the list of messages
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"messages": messages,
	})
}
