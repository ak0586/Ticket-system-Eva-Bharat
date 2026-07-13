package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken creates a new JWT for a logged-in user.
// It returns the token string and an error if something fails.
func GenerateToken(userID uint) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	
	// Claims are the "payload" data stored inside the JWT.
	claims := jwt.MapClaims{
		"user_id": userID,                                     // Store the user ID
		"exp":     time.Now().Add(24 * time.Hour).Unix(),      // Token expires in 24 hours
	}
	
	// Create the token using the HMAC SHA-256 signing algorithm.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// Sign the token with our secret key to prevent tampering.
	return token.SignedString([]byte(secret))
}

// ParseToken validates an incoming JWT string and extracts the user_id.
func ParseToken(tokenStr string) (uint, error) {
	secret := os.Getenv("JWT_SECRET")
	
	// Parse the token and verify the signature method matches what we expect.
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC (HS256) to prevent algorithm downgrade attacks.
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	
	// If parsing failed or the token is expired/invalid, reject it.
	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}

	// Extract the payload (claims) from the token.
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	// The JWT library parses numbers as float64 by default. 
	// We convert it back to a uint (unsigned integer) which matches our database ID type.
	userID := uint(claims["user_id"].(float64))
	return userID, nil
}
