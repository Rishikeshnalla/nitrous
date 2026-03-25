package handlers

import (
	"net/http"
	"testing"

	"nitrous-backend/database"
	"nitrous-backend/middleware"
	"nitrous-backend/models"

	"github.com/gin-gonic/gin"
)

func TestCreateEventEndpoint(t *testing.T) {
	setupHandlersTestEnv()
	database.Users = []models.User{
		{ID: "admin-1", Email: "admin@example.com", Role: "admin"},
		{ID: "user-1", Email: "user@example.com", Role: "user"},
	}

	r := gin.New()
	r.POST("/events", middleware.AuthMiddleware(), middleware.AdminMiddleware(), CreateEvent)

	validEvent := sampleEvent("")

	w := performJSONRequest(r, http.MethodPost, "/events", validEvent, "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 unauthorized, got %d", w.Code)
	}

	nonAdminToken := makeToken(t, "user-1")
	w = performJSONRequest(r, http.MethodPost, "/events", validEvent, nonAdminToken)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 forbidden for non-admin, got %d", w.Code)
	}

	token := makeToken(t, "admin-1")
	w = performRawRequest(r, http.MethodPost, "/events", "{", token)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid payload, got %d", w.Code)
	}

	w = performJSONRequest(r, http.MethodPost, "/events", validEvent, token)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 for successful create, got %d", w.Code)
	}
	if len(database.Events) != 1 {
		t.Fatalf("expected 1 event after create, got %d", len(database.Events))
	}
}

func TestUpdateEventEndpoint(t *testing.T) {
	setupHandlersTestEnv()
	database.Events = []models.Event{sampleEvent("event-1")}
	database.Users = []models.User{
		{ID: "admin-1", Email: "admin@example.com", Role: "admin"},
		{ID: "user-1", Email: "user@example.com", Role: "user"},
	}

	r := gin.New()
	r.PUT("/events/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), UpdateEvent)

	updated := sampleEvent("")
	updated.Title = "Updated Title"

	w := performJSONRequest(r, http.MethodPut, "/events/event-1", updated, "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 unauthorized, got %d", w.Code)
	}

	nonAdminToken := makeToken(t, "user-1")
	w = performJSONRequest(r, http.MethodPut, "/events/event-1", updated, nonAdminToken)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 forbidden for non-admin, got %d", w.Code)
	}

	token := makeToken(t, "admin-1")
	w = performRawRequest(r, http.MethodPut, "/events/event-1", "{", token)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid payload, got %d", w.Code)
	}

	w = performJSONRequest(r, http.MethodPut, "/events/nope", updated, token)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for not found, got %d", w.Code)
	}

	w = performJSONRequest(r, http.MethodPut, "/events/event-1", updated, token)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for successful update, got %d", w.Code)
	}
	if database.Events[0].Title != "Updated Title" {
		t.Fatalf("expected title to be updated")
	}
}

func TestDeleteEventEndpoint(t *testing.T) {
	setupHandlersTestEnv()
	database.Events = []models.Event{sampleEvent("event-1")}
	database.Users = []models.User{
		{ID: "admin-1", Email: "admin@example.com", Role: "admin"},
		{ID: "user-1", Email: "user@example.com", Role: "user"},
	}

	r := gin.New()
	r.DELETE("/events/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), DeleteEvent)

	w := performJSONRequest(r, http.MethodDelete, "/events/event-1", nil, "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 unauthorized, got %d", w.Code)
	}

	nonAdminToken := makeToken(t, "user-1")
	w = performJSONRequest(r, http.MethodDelete, "/events/event-1", nil, nonAdminToken)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 forbidden for non-admin, got %d", w.Code)
	}

	token := makeToken(t, "admin-1")
	w = performJSONRequest(r, http.MethodDelete, "/events/nope", nil, token)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 not found, got %d", w.Code)
	}

	w = performJSONRequest(r, http.MethodDelete, "/events/event-1", nil, token)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for successful delete, got %d", w.Code)
	}
	if len(database.Events) != 0 {
		t.Fatalf("expected event list to be empty after delete")
	}
}
