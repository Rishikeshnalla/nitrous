package handlers

import (
	"net/http"
	"nitrous-backend/database"

	"github.com/gin-gonic/gin"
)

// GetCategories returns all categories
func GetCategories(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"categories": database.Categories,
		"count":      len(database.Categories),
	})
}

// GetCategoryBySlug returns a single category by slug
func GetCategoryBySlug(c *gin.Context) {
	slug := c.Param("slug")
	
	for _, category := range database.Categories {
		if category.Slug == slug {
			c.JSON(http.StatusOK, category)
			return
		}
	}
	
	c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
}

// GetJourneys returns all journeys
func GetJourneys(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"journeys": database.Journeys,
		"count":    len(database.Journeys),
	})
}

// GetJourneyByID returns a single journey
func GetJourneyByID(c *gin.Context) {
	id := c.Param("id")
	
	for _, journey := range database.Journeys {
		if journey.ID == id {
			c.JSON(http.StatusOK, journey)
			return
		}
	}
	
	c.JSON(http.StatusNotFound, gin.H{"error": "Journey not found"})
}

// BookJourney handles journey booking
func BookJourney(c *gin.Context) {
	id := c.Param("id")
	
	for i, journey := range database.Journeys {
		if journey.ID == id {
			if journey.SlotsLeft <= 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "No slots available"})
				return
			}
			
			database.Journeys[i].SlotsLeft--
			
			c.JSON(http.StatusOK, gin.H{
				"message": "Journey booked successfully",
				"journey": database.Journeys[i],
			})
			return
		}
	}
	
	c.JSON(http.StatusNotFound, gin.H{"error": "Journey not found"})
}

// GetMerchItems returns all merch items
func GetMerchItems(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"items": database.MerchItems,
		"count": len(database.MerchItems),
	})
}

// GetMerchItemByID returns a single merch item
func GetMerchItemByID(c *gin.Context) {
	id := c.Param("id")
	
	for _, item := range database.MerchItems {
		if item.ID == id {
			c.JSON(http.StatusOK, item)
			return
		}
	}
	
	c.JSON(http.StatusNotFound, gin.H{"error": "Merch item not found"})
}
