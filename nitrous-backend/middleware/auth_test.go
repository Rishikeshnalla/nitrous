package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"nitrous-backend/config"
	"nitrous-backend/utils"

	"github.com/gin-gonic/gin"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	config.AppConfig.JWTSecret = "test-secret"

	r := gin.New()
	r.GET("/protected", AuthMiddleware(), func(c *gin.Context) {
		uid, ok := c.Get("userID")
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "missing userID"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"userID": uid})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for missing header, got %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "BadFormat token")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for malformed bearer format, got %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer not-a-token")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for invalid token, got %d", w.Code)
	}

	token, err := utils.GenerateJWT("user-123")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	req = httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for valid token, got %d", w.Code)
	}

	var payload map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if payload["userID"] != "user-123" {
		t.Fatalf("expected middleware to attach userID user-123, got %q", payload["userID"])
	}
}
