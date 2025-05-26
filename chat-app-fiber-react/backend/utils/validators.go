package utils

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/HodkiX4/react-go-projects/models"
	"github.com/gofiber/fiber/v2"
)

// Regular expression to validate usernames (3-20 characters, alphanumeric + underscore).
var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)

// IsValidUsername checks if the username is alphanumeric and 3-20 characters long.
func IsValidUsername(username string) bool {
	username = strings.TrimSpace(username)
	return usernameRegex.MatchString(username)
}

// Regular expression to validate email addresses.
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// IsValidEmail checks if the provided email matches a common email pattern.
func IsValidEmail(email string) bool {
	return emailRegex.MatchString(strings.TrimSpace(email))
}

/*
Regular expression to ensure minimum length and character type validity for passwords.
Note: This regex alone does not enforce presence of letters, digits, and special chars.
That logic is implemented separately in the function below.
*/
var passwordRegex = regexp.MustCompile(`^[A-Za-z0-9!@#$%^&*()_+={}\[\]:;"'<>,.?/-]{8,}$`)

/*
IsValidPassword checks if the password is:
- At least 8 characters long
- Contains at least one letter
- Contains at least one digit
- Contains at least one special character
*/
func IsValidPassword(password string) bool {
	if !passwordRegex.MatchString(password) {
		return false
	}

	hasLetter, hasDigit, hasSpecial := false, false, false

	for _, ch := range password {
		if unicode.IsLetter(ch) {
			hasLetter = true
		}
		if unicode.IsDigit(ch) {
			hasDigit = true
		}
		if strings.ContainsAny(string(ch), "!@#$%^&*()_+={}[\\]:;\"'<>,.?/-") {
			hasSpecial = true
		}
	}

	return hasLetter && hasDigit && hasSpecial
}

/*
ValidateSignupInput validates the fields of a SignupRequest struct.
It returns a Fiber-compatible error if any field is invalid.
This includes:
- Trimming and lowercasing input
- Ensuring username meets format and length rules
- Validating email format
- Ensuring password meets complexity requirements
*/

func ValidateSignupInput(signupRequest *models.SignupRequest) error {
	// Trim spaces
	signupRequest.Username = strings.TrimSpace(signupRequest.Username)
	signupRequest.Email = strings.TrimSpace(signupRequest.Email)
	signupRequest.Password = strings.TrimSpace(signupRequest.Password)

	// Normalize email to lowercase
	signupRequest.Email = strings.ToLower(signupRequest.Email)

	// Check username
	if !IsValidUsername(signupRequest.Username) {
		return RespondWithError(nil, fiber.StatusBadRequest, "Username must be alphanumeric and between 3 and 20 characters long")
	}

	// Check email
	if !IsValidEmail(signupRequest.Email) {
		return RespondWithError(nil, fiber.StatusBadRequest, "Invalid email format")
	}

	// Check password
	if !IsValidPassword(signupRequest.Password) {
		return RespondWithError(nil, fiber.StatusBadRequest, "Password must be at least 8 characters long and contain at least one letter, one number, and one special character")
	}

	return nil
}
