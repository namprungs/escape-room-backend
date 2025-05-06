package stats

import (
	"fmt"
	"math"
	"time"

	"github.com/FieldPs/escape-room-backend/internal/models"
	"gorm.io/gorm"
)

type SubjectStat struct {
	Subject    string `json:"subject"`
	Solved     int
	Total      int
	Percentage float64 `json:"percentage"`
}

type UserStatsResponse struct {
	SubjectStats  map[string]SubjectStat `json:"subject_stats"`
	CurrentStreak uint                   `json:"current_streak"`
	BestStreak    uint                   `json:"best_streak"`
	LastSolvedAt  time.Time              `json:"last_solved_at"`
}

func GetUserStats(db *gorm.DB, userID uint) (*UserStatsResponse, error) {
	// Initialize response with data from UserSolvedPuzzle
	var solvedPuzzle models.UserSolvedPuzzle
	if err := db.Where("user_id = ?", userID).First(&solvedPuzzle).Error; err != nil {
		return nil, fmt.Errorf("failed to get user solved puzzles: %w", err)
	}

	response := &UserStatsResponse{
		CurrentStreak: solvedPuzzle.CurrentStreak,
		BestStreak:    solvedPuzzle.BestStreak,
		LastSolvedAt:  solvedPuzzle.LastSolvedAt,
		SubjectStats:  make(map[string]SubjectStat),
	}

	// Get subject-level stats
	var userPuzzles []models.UserPuzzle
	if err := db.Preload("Puzzle").
		Where("user_id = ?", userID).
		Find(&userPuzzles).Error; err != nil {
		return nil, fmt.Errorf("failed to get user puzzles: %w", err)
	}

	// Calculate subject stats
	subjectTotals := make(map[string]int)
	subjectSolved := make(map[string]int)

	// First get all possible subjects and totals
	var puzzles []models.Puzzle
	if err := db.Find(&puzzles).Error; err != nil {
		return nil, fmt.Errorf("failed to get puzzles: %w", err)
	}

	for _, p := range puzzles {
		for _, subject := range p.Subjects {
			subjectTotals[subject]++
		}
	}

	// Then count solved per subject
	for _, up := range userPuzzles {
		for _, subject := range up.Puzzle.Subjects {
			subjectSolved[subject]++
		}
	}

	// Build subject stats
	for subject, total := range subjectTotals {
		percentage := math.Round(float64(subjectSolved[subject])/float64(total)*100*100) / 100

		response.SubjectStats[subject] = SubjectStat{
			Subject:    subject,
			Total:      total,
			Solved:     subjectSolved[subject],
			Percentage: percentage,
		}
	}

	return response, nil
}
