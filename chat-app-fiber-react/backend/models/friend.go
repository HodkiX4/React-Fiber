package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FriendRequest struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	SenderID   primitive.ObjectID `bson:"senderId" json:"senderId"`
	ReceiverID primitive.ObjectID `bson:"receiverId" json:"receiverId"`
	Status     string             `bson:"status" json:"status"` // e.g., "pending", "accepted", "declined"
	CreatedAt  time.Time          `bson:"createdAt"`
}
