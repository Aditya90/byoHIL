package main

import (
	"fmt"
	"log"
	"time"

	"backend/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func initDatabase() {
	// DSN matches the docker-compose.yml configuration
	dsn := "host=localhost user=hiluser password=hilpassword dbname=hildb port=5432 sslmode=disable TimeZone=UTC"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to postgres database:", err)
	}

	// Auto-migrate the Node schema
	err = db.AutoMigrate(&models.Node{})
	if err != nil {
		log.Fatal("Failed to auto migrate database schema:", err)
	}
	fmt.Println("Database connection and migration successful.")
}

func main() {
	// 1. Initialize DB
	initDatabase()

	// 2. Setup Fiber
	app := fiber.New()
	app.Use(logger.New())

	// 3. API Routes
	api := app.Group("/api/v1")

	// Helper to get nodes
	api.Get("/nodes", func(c *fiber.Ctx) error {
		var nodes []models.Node
		db.Find(&nodes)
		return c.JSON(nodes)
	})

	// Registration Endpoint
	api.Post("/nodes/register", func(c *fiber.Ctx) error {
		var req models.RegisterRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request payload"})
		}

		if req.Hostname == "" {
			return c.Status(400).JSON(fiber.Map{"error": "hostname is required"})
		}

		// Find existing node to avoid re-assigning ports
		var node models.Node
		result := db.Where("id = ?", req.Hostname).First(&node)

		if result.Error != nil {
			// Node doesn't exist, assignment logic
			// Find max port currently assigned to incrementally assign
			var lastNode models.Node
			port := 22000 // Starting port for reverse SSH tunnels
			if db.Order("assigned_ssh_port desc").First(&lastNode).Error == nil {
				if lastNode.AssignedSSHPort >= 22000 {
					port = lastNode.AssignedSSHPort + 1
				}
			}

			// Create fresh node
			node = models.Node{
				ID:              req.Hostname,
				Hostname:        req.Hostname,
				Status:          "online",
				AssignedSSHPort: port,
				LastSeenAt:      time.Now(),
			}
			db.Create(&node)
		} else {
			// Update existing node
			node.Status = "online"
			node.LastSeenAt = time.Now()
			db.Save(&node)
		}

		// Return port assignment so the agent can spin up autossh
		return c.JSON(fiber.Map{
			"status":            "registered",
			"assigned_ssh_port": node.AssignedSSHPort,
		})
	})

	// 4. Start Server
	log.Fatal(app.Listen(":8080"))
}
