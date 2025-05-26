package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Room struct {
	ID         primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name       string               `bson:"name" json:"name"`
	CreatedAt  time.Time            `bson:"createdAt" json:"createdAt"`
	Members    []primitive.ObjectID `bson:"members" json:"members"`
	IsPrivate  bool                 `bson:"isPrivate" json:"isPrivate"`
	Admins     []primitive.ObjectID `bson:"admins" json:"admins"`
	Moderators []primitive.ObjectID `bson:"moderators" json:"moderators"`
}
