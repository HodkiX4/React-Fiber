package handlers

import (
	"time"

	"github.com/HodkiX4/react-go-projects/config"
	"github.com/HodkiX4/react-go-projects/models"
	"github.com/HodkiX4/react-go-projects/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SendFriendRequest(c *fiber.Ctx) error {
	senderID, ok := c.Locals("userID").(primitive.ObjectID)
	if !ok {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid sender ID")
	}

	var body struct {
		ReceiverID string `json:"receiverId"`
	}

	if err := c.BodyParser(&body); err != nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	receiverID, err := primitive.ObjectIDFromHex(body.ReceiverID)
	if err != nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid receiver ID format")
	}

	if senderID == receiverID {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Cannot send request to self")
	}

	userCol := config.GetUserCollection()
	var sender models.User
	if err := userCol.FindOne(c.Context(), bson.M{"_id": senderID}).Decode(&sender); err == nil {
		for _, friend := range sender.Friends {
			if friend == receiverID {
				return utils.RespondWithError(c, fiber.StatusConflict, "You are already friends")
			}
		}
	}

	reqCol := config.GetFriendRequestCollection()
	filter := bson.M{
		"$or": []bson.M{
			{"senderID": senderID, "receiverID": receiverID, "status": "pending"},
			{"senderID": receiverID, "receiverID": senderID, "status": "pending"},
		},
	}
	count, err := reqCol.CountDocuments(c.Context(), filter)
	if err != nil || count > 0 {
		return utils.RespondWithError(c, fiber.StatusConflict, "Friend request already exists")
	}

	request := models.FriendRequest{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Status:     "pending",
		CreatedAt:  time.Now(),
	}

	_, err = reqCol.InsertOne(c.Context(), request)
	if err != nil {
		return utils.RespondWithError(c, fiber.StatusInternalServerError, "Failed to send friend request")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Friend request sent",
	})
}

func AcceptFriendRequest(c *fiber.Ctx) error {
	requestID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid request ID")
	}

	userID, ok := c.Locals("userID").(primitive.ObjectID)
	if !ok {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	reqCol := config.GetFriendRequestCollection()
	var request models.FriendRequest

	err = reqCol.FindOne(c.Context(), bson.M{"_id": requestID, "status": "pending"}).Decode(&request)
	if err != nil {
		return utils.RespondWithError(c, fiber.StatusNotFound, "Friend request not found")
	}

	if request.ReceiverID != userID {
		return utils.RespondWithError(c, fiber.StatusForbidden, "You are not authorize to accept this request")
	}

	_, err = reqCol.UpdateOne(c.Context(), bson.M{"_id": requestID}, bson.M{"$set": bson.M{"status": "accepted"}})
	if err != nil {
		return utils.RespondWithError(c, fiber.StatusInternalServerError, "Failed to accept request")
	}

	// Update each user's friend list
	userCol := config.GetUserCollection()
	_, _ = userCol.UpdateOne(c.Context(), bson.M{"_id": request.SenderID}, bson.M{"$addToSet": bson.M{"friends": request.ReceiverID}})
	_, _ = userCol.UpdateOne(c.Context(), bson.M{"_id": request.ReceiverID}, bson.M{"$addToSet": bson.M{"friends": request.SenderID}})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Friend request accepted",
	})
}

// RejectFriendRequest - POST /friend-request/reject/:id
func RejectFriendRequest(c *fiber.Ctx) error {
	requestID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid request ID")
	}

	userID, ok := c.Locals("userID").(primitive.ObjectID)
	if !ok {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	reqCol := config.GetFriendRequestCollection()
	var request models.FriendRequest

	err = reqCol.FindOne(c.Context(), bson.M{"_id": requestID, "status": "pending"}).Decode(&request)
	if err != nil {
		return utils.RespondWithError(c, fiber.StatusNotFound, "Pending request not found")
	}

	if request.ReceiverID != userID {
		return utils.RespondWithError(c, fiber.StatusForbidden, "You are not authorized to reject this request")
	}

	_, err = reqCol.UpdateOne(c.Context(), bson.M{"_id": requestID}, bson.M{"$set": bson.M{"status": "rejected"}})
	if err != nil {
		return utils.RespondWithError(c, fiber.StatusInternalServerError, "Failed to reject request")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Friend request rejected",
	})
}

// GetPendingFriendRequests - GET /friend-request/pending
func GetPendingFriendRequests(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(primitive.ObjectID)
	if !ok {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	collection := config.GetFriendRequestCollection()

	// Fetch incoming requests
	incomingCursor, err := collection.Find(c.Context(), bson.M{
		"receiverId": userID,
		"status":     "pending",
	})
	if err != nil {
		return utils.RespondWithError(c, fiber.StatusInternalServerError, "Failed to fetch incoming requests")
	}
	defer incomingCursor.Close(c.Context())

	var incomingRequests []models.FriendRequest
	if err := incomingCursor.All(c.Context(), &incomingRequests); err != nil {
		return utils.RespondWithError(c, fiber.StatusInternalServerError, "Failed to parse incoming requests")
	}

	// Fetch outgoing requests
	outgoingCursor, err := collection.Find(c.Context(), bson.M{
		"senderId": userID,
		"status":   "pending",
	})
	if err != nil {
		return utils.RespondWithError(c, fiber.StatusInternalServerError, "Failed to fetch outgoing requests")
	}
	defer outgoingCursor.Close(c.Context())

	var outgoingRequests []models.FriendRequest
	if err := outgoingCursor.All(c.Context(), &outgoingRequests); err != nil {
		return utils.RespondWithError(c, fiber.StatusInternalServerError, "Failed to parse outgoing requests")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":  "Friend requests fetched successfully",
		"incoming": incomingRequests,
		"outgoing": outgoingRequests,
	})
}
