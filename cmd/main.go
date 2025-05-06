package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/FieldPs/escape-room-backend/internal/routes"
	"github.com/FieldPs/escape-room-backend/migrations"
	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: Could not load .env file")
	}

	// Construct PostgreSQL DSN
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	// Connect to PostgreSQL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Silent), // Disables all SQL logging
	})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	// 1. Run migrations FIRST
	if err := migrations.MigrateAll(db); err != nil {
		log.Fatal("Migration failed:", err)
	}

	// Set up Gin router
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // เปลี่ยนเป็น URL ของ frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.OPTIONS("/*any", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	routes.SetupRoutes(r, db)

	// Run server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
