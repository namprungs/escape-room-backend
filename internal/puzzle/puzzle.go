package puzzle

import (
	"errors"
	"fmt"
	"time"

	"github.com/FieldPs/escape-room-backend/internal/models"
	"gorm.io/gorm"
)

type AnswerRequest struct {
	PuzzleID uint   `json:"puzzle_id" binding:"required"`
	Answer   string `json:"answer" binding:"required"`
}

type AnswerResponse struct {
	Correct       bool      `json:"correct"`
	Message       string    `json:"message"`
	CurrentStreak uint      `json:"current_streak,omitempty"`
	BestStreak    uint      `json:"best_streak,omitempty"`
	SolvedAt      time.Time `json:"solved_at,omitempty"`
}

func CheckAnswer(db *gorm.DB, userID uint, req AnswerRequest) (*AnswerResponse, error) {
	// Verify puzzle exists and get solution
	var p models.Puzzle
	if err := db.First(&p, req.PuzzleID).Error; err != nil {
		return nil, errors.New("puzzle not found")
	}
	fmt.Println("Puzzle ID:", p.ID, "Solution:", p.Solution)
	// Initialize response
	res := &AnswerResponse{
		Correct: req.Answer == p.Solution,
		Message: "Incorrect answer",
	}

	// Check if already solved
	if exists, err := alreadySolved(db, userID, req.PuzzleID); err != nil {
		return nil, err
	} else if exists {
		res.Message = "Already solved"
		return res, nil
	}

	if !res.Correct {
		return res, nil
	}

	// Process correct answer
	if err := recordSolve(db, userID, req.PuzzleID, res); err != nil {
		return nil, err
	}

	res.Message = "Correct answer!"
	return res, nil
}

// Helper functions
func alreadySolved(db *gorm.DB, userID uint, puzzleID uint) (bool, error) {
	var count int64
	err := db.Model(&models.UserPuzzle{}).
		Where("user_id = ? AND puzzle_id = ?", userID, puzzleID).
		Count(&count).Error
	return count > 0, err
}

func recordSolve(db *gorm.DB, userID uint, puzzleID uint, res *AnswerResponse) error {
	now := time.Now()

	return db.Transaction(func(tx *gorm.DB) error {
		// Create solve record
		if err := tx.Create(&models.UserPuzzle{
			UserID:   userID,
			PuzzleID: puzzleID,
			SolvedAt: now,
		}).Error; err != nil {
			return err
		}

		// 2. Get or create the user's stats record
		var stats models.UserSolvedPuzzle
		result := tx.Where(models.UserSolvedPuzzle{UserID: userID}).Attrs(models.UserSolvedPuzzle{
			SolvedPuzzles: 0,
			TotalPuzzles:  0, // Will update this below
			CurrentStreak: 0,
			BestStreak:    0,
		}).FirstOrCreate(&stats)

		if result.Error != nil {
			return fmt.Errorf("failed to get/create UserSolvedPuzzle: %w", result.Error)
		}

		// 3. Get total puzzles count
		var totalPuzzles int64
		if err := tx.Model(&models.Puzzle{}).Count(&totalPuzzles).Error; err != nil {
			return fmt.Errorf("failed to count puzzles: %w", err)
		}

		// Calculate streak
		stats.CurrentStreak = calculateStreak(stats.LastSolvedAt, now, stats.CurrentStreak)
		if stats.CurrentStreak > stats.BestStreak {
			stats.BestStreak = stats.CurrentStreak
		}

		// 5. Update all stats
		stats.SolvedPuzzles++
		stats.TotalPuzzles = uint(totalPuzzles) // Update with current total
		stats.LastSolvedAt = now

		// 6. Save the updated stats
		if err := tx.Model(&stats).Updates(models.UserSolvedPuzzle{
			SolvedPuzzles: stats.SolvedPuzzles,
			TotalPuzzles:  stats.TotalPuzzles,
			CurrentStreak: stats.CurrentStreak,
			BestStreak:    stats.BestStreak,
			LastSolvedAt:  stats.LastSolvedAt,
		}).Error; err != nil {
			return fmt.Errorf("failed to update UserSolvedPuzzle: %w", err)
		}

		// 7. Set response values
		res.CurrentStreak = stats.CurrentStreak
		res.BestStreak = stats.BestStreak
		res.SolvedAt = now

		return nil
	})
}

func calculateStreak(lastSolved time.Time, current time.Time, currentStreak uint) uint {
	if lastSolved.IsZero() {
		return 1
	}

	lastDay := lastSolved.Truncate(24 * time.Hour)
	currentDay := current.Truncate(24 * time.Hour)

	if lastDay.Add(24 * time.Hour).Equal(currentDay) {
		return currentStreak + 1
	}
	return 1
}
