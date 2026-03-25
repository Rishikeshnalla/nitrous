package handlers

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"nitrous-backend/database"
	"nitrous-backend/middleware"
	"nitrous-backend/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func TestRegisterFlow(t *testing.T) {
	setupHandlersTestEnv()

	r := gin.New()
	r.POST("/auth/register", Register)

	validReq := map[string]any{"email": "new@example.com", "password": "securepass123", "name": "New User"}
	w := performJSONRequest(r, http.MethodPost, "/auth/register", validReq, "")
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	w = performJSONRequest(r, http.MethodPost, "/auth/register", validReq, "")
	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409 for duplicate email, got %d", w.Code)
	}

	invalidReq := map[string]any{"email": "bad-email", "password": "short", "name": "x"}
	w = performJSONRequest(r, http.MethodPost, "/auth/register", invalidReq, "")
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid payload, got %d", w.Code)
	}
}

func TestLoginFlow(t *testing.T) {
	setupHandlersTestEnv()

	hash, err := bcrypt.GenerateFromPassword([]byte("securepass123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	database.Users = []models.User{
		{ID: "user-1", Email: "known@example.com", PasswordHash: string(hash), Name: "Known", CreatedAt: time.Now()},
	}

	r := gin.New()
	r.POST("/auth/login", Login)

	validReq := map[string]any{"email": "known@example.com", "password": "securepass123"}
	w := performJSONRequest(r, http.MethodPost, "/auth/login", validReq, "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	wrongPassReq := map[string]any{"email": "known@example.com", "password": "wrongpass"}
	w = performJSONRequest(r, http.MethodPost, "/auth/login", wrongPassReq, "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for wrong password, got %d", w.Code)
	}

	unknownReq := map[string]any{"email": "unknown@example.com", "password": "securepass123"}
	w = performJSONRequest(r, http.MethodPost, "/auth/login", unknownReq, "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for unknown user, got %d", w.Code)
	}

	invalidReq := map[string]any{"email": "not-an-email", "password": "x"}
	w = performJSONRequest(r, http.MethodPost, "/auth/login", invalidReq, "")
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid payload, got %d", w.Code)
	}
}

func TestGetCurrentUserFlow(t *testing.T) {
	setupHandlersTestEnv()

	database.Users = []models.User{
		{ID: "user-1", Email: "me@example.com", Name: "Me", CreatedAt: time.Now()},
	}

	r := gin.New()
	r.GET("/auth/me", middleware.AuthMiddleware(), GetCurrentUser)

	validToken := makeToken(t, "user-1")
	w := performJSONRequest(r, http.MethodGet, "/auth/me", nil, validToken)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var user models.User
	if err := json.Unmarshal(w.Body.Bytes(), &user); err != nil {
		t.Fatalf("failed to decode user: %v", err)
	}
	if user.ID != "user-1" {
		t.Fatalf("expected user-1, got %s", user.ID)
	}

	w = performJSONRequest(r, http.MethodGet, "/auth/me", nil, "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for missing token, got %d", w.Code)
	}

	w = performJSONRequest(r, http.MethodGet, "/auth/me", nil, "not-a-token")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for invalid token, got %d", w.Code)
	}

	missingUserToken := makeToken(t, "user-404")
	w = performJSONRequest(r, http.MethodGet, "/auth/me", nil, missingUserToken)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for user-not-found path, got %d", w.Code)
	}
}
