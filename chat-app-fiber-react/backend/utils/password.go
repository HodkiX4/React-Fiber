package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/argon2"
)

// HashPassword creates a hash from a password using the provided salt.
func HashPassword(password string, salt []byte) string {
	// Generate Argon2id hash with 1 iteration, 64MB memory, 4 threads, 32-byte key
	hashedPassword := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	// Encode hash to base64 string
	return base64.StdEncoding.EncodeToString(hashedPassword)
}

// GenerateSalt creates a secure random salt of the given size.
func GenerateSalt(size int) ([]byte, error) {
	salt := make([]byte, size)
	_, err := rand.Read(salt) // Fill slice with random bytes
	if err != nil {
		return nil, err
	}
	return salt, nil
}

// ComparePassword checks if a plain password matches the stored hash using the stored salt.
func ComparePassword(plainPassword string, encodedSalt string, encodedHash string) (bool, error) {
	// Decode base64 salt
	salt, err := base64.StdEncoding.DecodeString(encodedSalt)
	if err != nil {
		return false, err
	}

	// Decode base64 hash
	expectedHash, err := base64.StdEncoding.DecodeString(encodedHash)
	if err != nil {
		return false, err
	}

	// Re-hash the input password using the decoded salt
	providedHash := argon2.IDKey([]byte(plainPassword), salt, 1, 64*1024, 4, 32)

	// Compare the re-hashed password with the stored hash
	return bytes.Equal(providedHash, expectedHash), nil
}
