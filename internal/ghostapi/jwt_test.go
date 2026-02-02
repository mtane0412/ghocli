/**
 * jwt_test.go
 * Test code for JWT generation functionality for Ghost Admin API
 */

package ghostapi

import (
	"encoding/hex"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TestGenerateJWT_GenerateCorrectFormatToken generates a token in the correct format
func TestGenerateJWT_GenerateCorrectFormatToken(t *testing.T) {
	// Test API key information
	keyID := "64fac5417c4c6b0001234567"
	secret := "89abcdef01234567890123456789abcd01234567890123456789abcdef0123"

	// Generate JWT
	token, err := GenerateJWT(keyID, secret)
	if err != nil {
		t.Fatalf("Failed to generate JWT: %v", err)
	}

	// Verify token is not empty
	if token == "" {
		t.Error("Generated token is empty")
	}

	// Verify token consists of three parts (header.payload.signature)
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Errorf("Number of token parts = %d; want 3", len(parts))
	}
}

// TestGenerateJWT_VerifyToken verifies the token
func TestGenerateJWT_VerifyToken(t *testing.T) {
	keyID := "64fac5417c4c6b0001234567"
	secret := "89abcdef01234567890123456789abcd01234567890123456789abcdef0123"

	// Generate JWT
	tokenString, err := GenerateJWT(keyID, secret)
	if err != nil {
		t.Fatalf("Failed to generate JWT: %v", err)
	}

	// Parse and verify token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify signing algorithm is HS256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			t.Errorf("Unexpected signing algorithm: %v", token.Header["alg"])
		}
		// Decode secret from hex to binary
		secretBytes, err := hex.DecodeString(secret)
		if err != nil {
			t.Fatalf("Failed to decode secret: %v", err)
		}
		return secretBytes, nil
	})

	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}

	if !token.Valid {
		t.Error("Token is invalid")
	}

	// Verify claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("Failed to retrieve claims")
	}

	// Verify aud claim is the Ghost Admin API path
	aud, ok := claims["aud"].(string)
	if !ok || aud != "/admin/" {
		t.Errorf("aud = %q; want %q", aud, "/admin/")
	}

	// Verify iat claim is near current time
	iat, ok := claims["iat"].(float64)
	if !ok {
		t.Fatal("Failed to retrieve iat claim")
	}
	iatTime := time.Unix(int64(iat), 0)
	if time.Since(iatTime) > 10*time.Second {
		t.Errorf("iat is too old: %v", iatTime)
	}

	// Verify exp claim is iat + 5 minutes
	exp, ok := claims["exp"].(float64)
	if !ok {
		t.Fatal("Failed to retrieve exp claim")
	}
	expectedExp := int64(iat) + 5*60
	if int64(exp) != expectedExp {
		t.Errorf("exp = %d; want %d", int64(exp), expectedExp)
	}
}

// TestGenerateJWT_HeaderContainsKid verifies kid is included in header
func TestGenerateJWT_HeaderContainsKid(t *testing.T) {
	keyID := "64fac5417c4c6b0001234567"
	secret := "89abcdef01234567890123456789abcd01234567890123456789abcdef0123"

	// Generate JWT
	tokenString, err := GenerateJWT(keyID, secret)
	if err != nil {
		t.Fatalf("Failed to generate JWT: %v", err)
	}

	// Parse token (without verification)
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}

	// Verify kid in header
	kid, ok := token.Header["kid"].(string)
	if !ok {
		t.Fatal("kid is not included in header")
	}
	if kid != keyID {
		t.Errorf("kid = %q; want %q", kid, keyID)
	}
}

// TestGenerateJWT_ErrorWithEmptyKeyID tests error with empty key ID
func TestGenerateJWT_ErrorWithEmptyKeyID(t *testing.T) {
	_, err := GenerateJWT("", "secret")
	if err == nil {
		t.Error("No error returned with empty key ID")
	}
}

// TestGenerateJWT_ErrorWithEmptySecret tests error with empty secret
func TestGenerateJWT_ErrorWithEmptySecret(t *testing.T) {
	_, err := GenerateJWT("keyid", "")
	if err == nil {
		t.Error("No error returned with empty secret")
	}
}
