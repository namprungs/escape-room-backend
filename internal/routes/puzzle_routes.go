package routes

import (
	"net/http"

	"github.com/FieldPs/escape-room-backend/internal/puzzle"
	"github.com/FieldPs/escape-room-backend/internal/stats"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterPuzzleRoutes sets up puzzle and stats endpoints
func RegisterPuzzleRoutes(r gin.IRouter, db *gorm.DB) {
	// Protected routes under /api
	authGroup := r.Group("/", AuthMiddleware())
	{
		authGroup.GET("/stats", statsHandler(db))
		authGroup.POST("/submit_answer", SubmitAnswerHandler(db))
	}
}

// statsHandler retrieves user statistics
func statsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uint)
		stats, err := stats.GetUserStats(db, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stats"})
			return
		}
		c.JSON(http.StatusOK, stats)
	}
}

func SubmitAnswerHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uint)

		var req puzzle.AnswerRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request",
				"details": err.Error(), // This will show binding errors
			})
			return
		}

		res, err := puzzle.CheckAnswer(db, userID, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if res.Correct {
			c.JSON(http.StatusOK, res)
		} else {
			c.JSON(http.StatusBadRequest, res)
		}
	}
}
