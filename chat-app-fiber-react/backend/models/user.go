package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserStatus string

const (
	StatusOnline  UserStatus = "online"
	StatusOffline UserStatus = "offline"
	StatusAway    UserStatus = "away"
	StatusBusy    UserStatus = "busy"
)

/*
func IsValidStatus(status string) bool {
	switch UserStatus(status) {
	case StatusOnline, StatusOffline, StatusAway, StatusBusy:
		return true
	default:
		return false
	}
}
*/

type SignupRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Username  string               `bson:"username" json:"username"`
	Email     string               `bson:"email" json:"email"`
	Password  string               `bson:"password,omitempty"`
	Salt      string               `bson:"salt,omitempty"`
	AvatarURL string               `bson:"avatarUrl,omitempty" json:"avatarUrl,omitempty"`
	Status    UserStatus           `bson:"status,omitempty" json:"status,omitempty"`
	Friends   []primitive.ObjectID `bson:"friends,omitempty" json:"friends,omitempty"`
	CreatedAt time.Time            `bson:"createdAt" json:"createdAt"`
}

type SafeUser struct {
	ID        string   `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	AvatarURL string   `json:"avatarUrl,omitempty"`
	Status    string   `json:"status,omitempty"`
	Friends   []string `json:"friends,omitempty"`
	CreatedAt string   `json:"createdAt,omitempty"`
}
