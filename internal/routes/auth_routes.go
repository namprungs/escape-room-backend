package routes

import (
	"errors"
	"net/http"
	"strings"

	"github.com/FieldPs/escape-room-backend/internal/auth"
	"github.com/FieldPs/escape-room-backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterAuthRoutes sets up authentication-related endpoints
func RegisterAuthRoutes(r gin.IRouter, db *gorm.DB) {
	r.POST("/register", registerHandler(db))
	r.POST("/login", loginHandler(db))
}

// AuthMiddleware protects routes with JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := auth.ValidateJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		c.Set("userID", uint(claims.UserID))
		c.Next()
	}
}

// registerHandler handles user registration
func registerHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required,min=6"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Hash the password
		hash, err := auth.HashPassword(input.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		// Use transaction to create both user and solved puzzle record
		err = db.Transaction(func(tx *gorm.DB) error {
			// Create user
			user := models.User{
				Username:     input.Username,
				PasswordHash: hash,
			}

			if err := tx.Create(&user).Error; err != nil {
				return err // This will rollback the transaction
			}

			// Create initial UserSolvedPuzzle record
			var totalPuzzles int64
			if err := tx.Model(&models.Puzzle{}).Count(&totalPuzzles).Error; err != nil {
				return err
			}

			solvedPuzzle := models.UserSolvedPuzzle{
				UserID:        user.ID,
				SolvedPuzzles: 0,
				TotalPuzzles:  uint(totalPuzzles),
				CurrentStreak: 0,
				BestStreak:    0,
			}

			if err := tx.Create(&solvedPuzzle).Error; err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				c.JSON(http.StatusConflict, gin.H{"error": "Username already taken"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Registration failed"})
			}
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "User registered successfully",
		})
	}
}

// loginHandler handles user login
func loginHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Define input structure with validation tags
		var input struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		// Bind JSON input - THIS WAS MISSING
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		// Trim whitespace
		input.Username = strings.TrimSpace(input.Username)
		input.Password = strings.TrimSpace(input.Password)

		// Find user by username
		var user models.User
		if err := db.Where("username = ?", input.Username).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Simulate password check to prevent timing attacks
				auth.CheckPasswordHash("dummy_password", "$2a$10$dummyhashdummyhashdummyhashdu")
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			}
			return
		}

		// Verify password
		if !auth.CheckPasswordHash(input.Password, user.PasswordHash) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		// Generate JWT token
		token, err := auth.GenerateJWT(int(user.ID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to generate token",
				"details": err.Error(),
			})
			return
		}

		// Successful login response
		c.JSON(http.StatusOK, gin.H{
			"token": token,
			"user": gin.H{
				"username": user.Username,
				// Never include sensitive information here
			},
		})
	}
}
