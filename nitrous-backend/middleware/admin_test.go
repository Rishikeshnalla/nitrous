package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"nitrous-backend/config"
	"nitrous-backend/database"
	"nitrous-backend/models"
	"nitrous-backend/utils"

	"github.com/gin-gonic/gin"
)

func TestAdminMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	config.AppConfig.JWTSecret = "test-secret"
	database.Users = nil

	r := gin.New()
	r.GET("/admin-only", AuthMiddleware(), AdminMiddleware(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/admin-only", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for missing auth header, got %d", w.Code)
	}

	userToken, err := utils.GenerateJWT("user-1")
	if err != nil {
		t.Fatalf("failed to generate user token: %v", err)
	}

	req = httptest.NewRequest(http.MethodGet, "/admin-only", nil)
	req.Header.Set("Authorization", "Bearer "+userToken)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 when token user does not exist, got %d", w.Code)
	}

	database.Users = []models.User{{ID: "user-1", Email: "user@example.com", Role: "user"}}

	req = httptest.NewRequest(http.MethodGet, "/admin-only", nil)
	req.Header.Set("Authorization", "Bearer "+userToken)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for non-admin user, got %d", w.Code)
	}

	database.Users = append(database.Users, models.User{ID: "admin-1", Email: "admin@example.com", Role: "admin"})
	adminToken, err := utils.GenerateJWT("admin-1")
	if err != nil {
		t.Fatalf("failed to generate admin token: %v", err)
	}

	req = httptest.NewRequest(http.MethodGet, "/admin-only", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for admin user, got %d", w.Code)
	}

	var payload map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if payload["message"] != "ok" {
		t.Fatalf("expected success payload from admin-only route")
	}
}
