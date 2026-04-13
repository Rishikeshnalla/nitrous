package handlers

import (
	"net/http"
	"nitrous-backend/database"
	"nitrous-backend/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SetReminder(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.SetReminderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.RemindAt.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Reminder time must be in the future"})
		return
	}

	database.Mu.RLock()

	eventExists := false
	for _, event := range database.Events {
		if event.ID == req.EventID {
			eventExists = true
			break
		}
	}
	database.Mu.RUnlock()

	if !eventExists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	reminder := models.Reminder{
		ID:        uuid.New().String(),
		UserID:    userID,
		EventID:   eventID,
		CreatedAt: time.Now(),
	}
	database.AppendReminder(r)

	database.Mu.Lock()
	database.Reminders = append(database.Reminders, reminder)
	database.Mu.Unlock()

	c.JSON(http.StatusCreated, reminder)
}

func DeleteReminder(c *gin.Context) {
	if !database.DeleteReminder(c.GetString("userID"), c.Param("id")) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reminder not found"})
		return
	}

	reminderID := c.Param("id")
	database.Mu.Lock()

	for i, reminder := range database.Reminders {
		if reminder.ID == reminderID {
			if reminder.UserID != userID.(string) {
				database.Mu.Unlock()
				c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
				return
			}

			database.Reminders = append(database.Reminders[:i], database.Reminders[i+1:]...)
			database.Mu.Unlock()
			c.JSON(http.StatusOK, gin.H{"message": "Reminder deleted"})
			return
		}
	}
	database.Mu.Unlock()

	c.JSON(http.StatusNotFound, gin.H{"error": "Reminder not found"})
}

func GetMyReminders(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	database.Mu.RLock()
	defer database.Mu.RUnlock()

	var reminders []models.Reminder
	for _, reminder := range database.Reminders {
		if reminder.UserID == userID.(string) {
			reminders = append(reminders, reminder)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"reminders": reminders,
		"count":     len(reminders),
	})
}
