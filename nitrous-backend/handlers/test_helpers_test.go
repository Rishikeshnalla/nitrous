package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"nitrous-backend/config"
	"nitrous-backend/database"
	"nitrous-backend/models"
	"nitrous-backend/utils"

	"github.com/gin-gonic/gin"
)

func setupHandlersTestEnv() {
	gin.SetMode(gin.TestMode)
	config.AppConfig.JWTSecret = "test-secret"
	database.Events = nil
	database.Categories = nil
	database.Journeys = nil
	database.MerchItems = nil
	database.Users = nil
	database.Teams = nil
	database.Reminders = nil
	database.Orders = nil
}

func makeToken(t *testing.T, userID string) string {
	t.Helper()
	token, err := utils.GenerateJWT(userID)
	if err != nil {
		t.Fatalf("failed generating token: %v", err)
	}
	return token
}

func performJSONRequest(r http.Handler, method string, path string, body any, token string) *httptest.ResponseRecorder {
	var payload []byte
	if body != nil {
		payload, _ = json.Marshal(body)
	}

	req := httptest.NewRequest(method, path, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func performRawRequest(r http.Handler, method string, path string, rawBody string, token string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, bytes.NewBufferString(rawBody))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func sampleEvent(id string) models.Event {
	return models.Event{
		ID:       id,
		Title:    "Sample Event",
		Location: "Test Track",
		Date:     time.Now().Add(24 * time.Hour).UTC(),
		IsLive:   false,
		Category: "motorsport",
		Time:     "12:00 UTC",
	}
}
