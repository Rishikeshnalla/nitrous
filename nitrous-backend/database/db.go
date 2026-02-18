package database

import (
	"log"
	"nitrous-backend/models"
	"time"

	"github.com/google/uuid"
)

// In-memory storage for prototype
// Replace with actual DB connection for production
var (
	Events     []models.Event
	Categories []models.Category
	Journeys   []models.Journey
	MerchItems []models.MerchItem
	Users      []models.User
)

func InitDB() {
	log.Println("Initializing in-memory database...")
	
	// Seed data
	seedEvents()
	seedCategories()
	seedJourneys()
	seedMerch()
	
	log.Println("✓ Database initialized with seed data")
}

func CloseDB() {
	log.Println("Closing database connection...")
}

func seedEvents() {
	Events = []models.Event{
		{
			ID:       uuid.New().String(),
			Title:    "NASCAR Daytona 500",
			Location: "Daytona International Speedway · Florida",
			Date:     time.Now().Add(10 * 24 * time.Hour),
			IsLive:   true,
			Category: "motorsport",
			Time:     "15:00 UTC",
		},
		{
			ID:       uuid.New().String(),
			Title:    "Dakar Rally — Stage 9",
			Location: "Al Ula → Ha'il · Saudi Arabia",
			Date:     time.Now().Add(-2 * 24 * time.Hour),
			IsLive:   false,
			Category: "offroad",
			Time:     "09:00 UTC",
		},
		{
			ID:       uuid.New().String(),
			Title:    "World Dirt Track Championship",
			Location: "Knob Noster · Missouri, USA",
			Date:     time.Now().Add(5 * 24 * time.Hour),
			IsLive:   true,
			Category: "motorsport",
			Time:     "18:00 UTC",
		},
		{
			ID:       uuid.New().String(),
			Title:    "Speed Boat Cup — Finals",
			Location: "Lake Como · Italy",
			Date:     time.Now().Add(14 * 24 * time.Hour),
			IsLive:   false,
			Category: "water",
			Time:     "14:00 UTC",
		},
		{
			ID:       uuid.New().String(),
			Title:    "Red Bull Skydive Series — Rd. 3",
			Location: "Interlaken Drop Zone · Switzerland",
			Date:     time.Now().Add(20 * 24 * time.Hour),
			IsLive:   false,
			Category: "air",
			Time:     "11:30 UTC",
		},
		{
			ID:       uuid.New().String(),
			Title:    "Crop Duster Air Racing",
			Location: "Bakersfield Airfield · California",
			Date:     time.Now().Add(26 * 24 * time.Hour),
			IsLive:   false,
			Category: "air",
			Time:     "16:00 UTC",
		},
	}
}

func seedCategories() {
	Categories = []models.Category{
		{
			ID:          uuid.New().String(),
			Name:        "MOTORSPORT",
			Slug:        "motorsport",
			Icon:        "🏎️",
			LiveCount:   24,
			Description: "NASCAR · F1 · Dirt · Rally",
			Color:       "cyan",
		},
		{
			ID:          uuid.New().String(),
			Name:        "WATER",
			Slug:        "water",
			Icon:        "🌊",
			LiveCount:   8,
			Description: "Speed Boats · Jet Ski · Surf",
			Color:       "blue",
		},
		{
			ID:          uuid.New().String(),
			Name:        "AIR & SKY",
			Slug:        "air",
			Icon:        "🪂",
			LiveCount:   5,
			Description: "Skydive · Air Race · Wing",
			Color:       "purple",
		},
		{
			ID:          uuid.New().String(),
			Name:        "OFF-ROAD",
			Slug:        "offroad",
			Icon:        "🏔️",
			LiveCount:   12,
			Description: "Dakar · Baja · Enduro",
			Color:       "orange",
		},
	}
}

func seedJourneys() {
	Journeys = []models.Journey{
		{
			ID:          uuid.New().String(),
			Title:       "DAYTONA PIT CREW EXPERIENCE",
			Category:    "MOTORSPORT · BEHIND THE SCENES",
			Description: "Go behind the wall at Daytona 500. Watch pit stops up close, meet the crew chiefs, and ride the pace car on track.",
			Badge:       "EXCLUSIVE",
			SlotsLeft:   12,
			Date:        time.Now().Add(10 * 24 * time.Hour),
			Price:       2400,
		},
		{
			ID:          uuid.New().String(),
			Title:       "DAKAR DESERT CONVOY",
			Category:    "RALLY · DESERT EXPEDITION",
			Description: "Ride a support vehicle through the Dakar stages. Sleep under the stars, eat with the team, and feel the dust.",
			Badge:       "MEMBERS ONLY",
			SlotsLeft:   6,
			Date:        time.Now().Add(345 * 24 * time.Hour),
			Price:       5800,
		},
		{
			ID:          uuid.New().String(),
			Title:       "RED BULL TANDEM SKYDIVE",
			Category:    "AIR · EXTREME SPORT",
			Description: "Jump with a Red Bull certified instructor at 15,000ft. Camera-equipped, full debrief, and a story you'll never forget.",
			Badge:       "LIMITED",
			SlotsLeft:   3,
			Date:        time.Now().Add(20 * 24 * time.Hour),
			Price:       1200,
		},
	}
}

func seedMerch() {
	MerchItems = []models.MerchItem{
		{ID: uuid.New().String(), Name: "Team Hoodie", Icon: "👕", Price: 89, Category: "apparel"},
		{ID: uuid.New().String(), Name: "NITROUS Cap", Icon: "🧢", Price: 42, Category: "apparel"},
		{ID: uuid.New().String(), Name: "Racing Jacket", Icon: "🏎️", Price: 189, Category: "apparel"},
		{ID: uuid.New().String(), Name: "Pit Watch", Icon: "⌚", Price: 249, Category: "accessories"},
		{ID: uuid.New().String(), Name: "Gear Backpack", Icon: "🎒", Price: 120, Category: "accessories"},
		{ID: uuid.New().String(), Name: "Drop Keychain", Icon: "🏆", Price: 28, Category: "collectibles"},
	}
}
