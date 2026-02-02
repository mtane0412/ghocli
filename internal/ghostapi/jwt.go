/**
 * jwt.go
 * JWT generation for Ghost Admin API
 *
 * Ghost Admin API requires JWT tokens signed with HS256 algorithm.
 * Token expiration is 5 minutes.
 */

package ghostapi

import (
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateJWT generates a JWT token for Ghost Admin API.
// keyID: ID part of the Admin API key
// secret: Secret part of the Admin API key
func GenerateJWT(keyID, secret string) (string, error) {
	if keyID == "" {
		return "", errors.New("key ID is empty")
	}
	if secret == "" {
		return "", errors.New("secret is empty")
	}

	// Current time (in seconds)
	now := time.Now().Unix()

	// Set JWT claims
	claims := jwt.MapClaims{
		"iat": now,           // Issued at
		"exp": now + 5*60,    // Expiration (5 minutes later)
		"aud": "/admin/",     // Ghost Admin API path
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Set key ID in header
	token.Header["kid"] = keyID

	// Decode secret from hex to binary
	secretBytes, err := hex.DecodeString(secret)
	if err != nil {
		return "", errors.New("failed to decode secret from hex")
	}

	// Sign with decoded secret
	tokenString, err := token.SignedString(secretBytes)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
