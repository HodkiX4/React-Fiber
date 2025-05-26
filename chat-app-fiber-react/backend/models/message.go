package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	RoomID    primitive.ObjectID `bson:"roomId" json:"roomId"`
	SenderID  primitive.ObjectID `bson:"senderId" json:"senderId"`
	Content   string             `bson:"content" json:"content"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}
