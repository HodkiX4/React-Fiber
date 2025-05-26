package config

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName                = "Chat_App"
	userCollName          = "users"
	roomCollName          = "rooms"
	messageCollName       = "messages"
	friendRequestCollName = "friend_requests"
)

// MongoDB collection references, initialized after connecting to the database
var userCollection *mongo.Collection
var roomCollection *mongo.Collection
var messageCollection *mongo.Collection
var friendRequestCollection *mongo.Collection

/**
 * GetUserCollection returns a reference to the users collection.
 * This getter prevents accidental reassignment of the internal collection variable.
 *
 * @returns { *mongo.Collection } - The MongoDB collection for users.
 */
func GetUserCollection() *mongo.Collection {
	return userCollection
}

/**
 * GetRoomCollection returns a reference to the rooms collection.
 *
 * @returns { *mongo.Collection } - The MongoDB collection for rooms.
 */
func GetRoomCollection() *mongo.Collection {
	return roomCollection
}

/**
 * GetMessageCollection returns a reference to the messages collection.
 *
 * @returns { *mongo.Collection } - The MongoDB collection for messages.
 */
func GetMessageCollection() *mongo.Collection {
	return messageCollection
}

/**
 * GetFriendRequestCollection returns a reference to the friend_requests collection.
 *
 * @returns { *mongo.Collection } - The MongoDB collection for friend requests.
 */
func GetFriendRequestCollection() *mongo.Collection {
	return friendRequestCollection
}

/**
 * ConnectDB initializes a connection to MongoDB using the URI
 * specified in the environment variable MONGODB_URI.
 * It also initializes references to all required collections and
 * creates necessary indexes for them.
 *
 * Logs fatal errors if connection or indexing fails.
 */
func ConnectDB() {
	// Attempt to retrieve the MongoDB URI from environment variables
	MONGODB_URI, exists := os.LookupEnv("MONGODB_URI")
	if !exists {
		log.Fatal("MONGODB_URI not set in environment")
	}

	// Create client options using the MongoDB URI
	clientOptions := options.Client().ApplyURI(MONGODB_URI)

	// Create a context with a timeout to ensure the connection doesn't hang indefinitely
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Ensure that the context is properly canceled after use

	// Attempt to connect to MongoDB with the provided client options
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ping the MongoDB server to verify connection and authentication
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	// Log successful connection
	log.Println("[MongoDB] Connected to database")

	// Initialize the collection variables using the connected client's database
	userCollection = client.Database(dbName).Collection(userCollName)
	roomCollection = client.Database(dbName).Collection(roomCollName)
	messageCollection = client.Database(dbName).Collection(messageCollName)
	friendRequestCollection = client.Database(dbName).Collection(friendRequestCollName)

	// Call the index creation logic to set up necessary database indexes
	createIndexes(ctx)
}

/**
 * createIndexes sets up the required indexes for each collection.
 *
 * @param ctx - The context used for index creation.
 */
func createIndexes(ctx context.Context) {
	if err := createUserIndexes(ctx); err != nil {
		log.Fatalf("Error creating user indexes: %v", err)
	}

	if err := createRoomIndexes(ctx); err != nil {
		log.Fatalf("Error creating room indexes: %v", err)
	}

	if err := createMessageIndexes(ctx); err != nil {
		log.Fatalf("Error creating message indexes: %v", err)
	}

	if err := createFriendRequestIndexes(ctx); err != nil {
		log.Fatalf("Error creating message indexes: %v", err)
	}
	log.Println("[MongoDB] Indexes created successfully")
}

/**
 * createUserIndexes creates unique indexes on the 'email' and 'username' fields
 * to ensure data integrity in the users collection.
 *
 * @param ctx - The context used for index creation.
 * @returns { error } - Any error encountered during index creation.
 */
func createUserIndexes(ctx context.Context) error {
	emailIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "email", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	usernameIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "username", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	_, err := userCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{emailIndex, usernameIndex})
	return err
}

/**
 * createRoomIndexes creates a unique index on the 'name' field
 * in the rooms collection.
 *
 * @param ctx - The context used for index creation.
 * @returns { error } - Any error encountered during index creation.
 */
func createRoomIndexes(ctx context.Context) error {
	roomIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "name", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	_, err := roomCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{roomIndex})
	return err
}

/**
 * createMessageIndexes creates a compound index on 'roomId' and 'createdAt'
 * for efficient querying and sorting of messages.
 *
 * @param ctx - The context used for index creation.
 * @returns { error } - Any error encountered during index creation.
 */
func createMessageIndexes(ctx context.Context) error {
	messageIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "roomId", Value: 1},
			{Key: "createdAt", Value: -1},
		},
		Options: options.Index().SetName("roomId_createdAt_index"),
	}
	_, err := messageCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{messageIndex})
	return err
}

func createFriendRequestIndexes(ctx context.Context) error {
	// Prevent duplicate friend requests
	uniqueRequestIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "senderId", Value: 1},
			{Key: "receiverId", Value: 1},
		},
		Options: options.Index().SetUnique(true).SetName("sender_receiver_unique"),
	}

	// For efficient lookups
	senderIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "senderId", Value: 1},
		},
	}
	receiverIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "receiverId", Value: 1},
		},
	}

	_, err := friendRequestCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		uniqueRequestIndex,
		senderIndex,
		receiverIndex,
	})
	return err
}
