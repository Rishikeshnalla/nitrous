package utils

import (
	"testing"
	"time"

	"nitrous-backend/config"

	"github.com/golang-jwt/jwt/v5"
)

func TestJWTUtility(t *testing.T) {
	config.AppConfig.JWTSecret = "test-secret"

	token, err := GenerateJWT("user-1")
	if err != nil {
		t.Fatalf("expected token generation success, got error: %v", err)
	}
	if token == "" {
		t.Fatalf("expected non-empty token")
	}

	claims, err := ValidateJWT(token)
	if err != nil {
		t.Fatalf("expected token validation success, got error: %v", err)
	}
	if claims.UserID != "user-1" {
		t.Fatalf("expected user-1 in claims, got %s", claims.UserID)
	}

	expiredClaims := Claims{
		UserID: "user-1",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	expiredTokenStr, err := expiredToken.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		t.Fatalf("failed to sign expired token: %v", err)
	}
	if _, err := ValidateJWT(expiredTokenStr); err == nil {
		t.Fatalf("expected expired token to be rejected")
	}

	if _, err := ValidateJWT("this-is-not-a-jwt"); err == nil {
		t.Fatalf("expected invalid token to be rejected")
	}
}
