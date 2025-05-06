package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes configures all API endpoints
func SetupRoutes(r *gin.Engine, db *gorm.DB) {

	apiV1 := r.Group("/api/v1")
	{
		RegisterAuthRoutes(apiV1, db)
		RegisterPuzzleRoutes(apiV1, db)
	}

	r.GET("/healthz", func(c *gin.Context) {
		sqlDB, err := db.DB()
		if err != nil {
			c.JSON(500, gin.H{"status": "unhealthy", "error": err.Error()})
			return
		}

		if err := sqlDB.Ping(); err != nil {
			c.JSON(500, gin.H{"status": "unhealthy", "error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"status": "healthy"})
	})
}
