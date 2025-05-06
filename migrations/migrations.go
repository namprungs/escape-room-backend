package migrations

import (
	"math/rand"
	"time"

	"gorm.io/gorm"

	"github.com/FieldPs/escape-room-backend/internal/auth"
	"github.com/FieldPs/escape-room-backend/internal/models"
	"github.com/lib/pq"
)

func MigrateAll(db *gorm.DB) error {
	// 1. Auto-migrate all models FIRST
	err := db.AutoMigrate(
		&models.Puzzle{},
		&models.User{},
		&models.UserPuzzle{},
		&models.UserSolvedPuzzle{},
	)
	if err != nil {
		return err
	}

	// 2. Now seed data
	return SeedData(db)
}

// Available subjects to randomize from
var availableSubjects = []string{"Physics", "Chemistry", "Biology", "Math", "Thai", "English", "Social"}

func randomSubjects() pq.StringArray {
	count := 3 // Number of subjects to pick
	subjects := make([]string, count)

	// Shuffle and pick first 'count' elements
	rand.Shuffle(len(availableSubjects), func(i, j int) {
		availableSubjects[i], availableSubjects[j] = availableSubjects[j], availableSubjects[i]
	})

	copy(subjects, availableSubjects[:count])
	return pq.StringArray(subjects)
}

func SeedPuzzles(db *gorm.DB) error {
	puzzles := []models.Puzzle{
		{
			ID:        1,
			Title:     "First Puzzle",
			Content:   "This is the first puzzle.",
			Solution:  "101",
			Subjects:  randomSubjects(),
			CreatedAt: time.Now().Add(-148 * time.Hour),
		},
		{
			ID:        2,
			Title:     "Second Puzzle",
			Content:   "This is the second puzzle.",
			Solution:  "202",
			Subjects:  randomSubjects(),
			CreatedAt: time.Now().Add(-120 * time.Hour),
		},
		{
			ID:        3,
			Title:     "Third Puzzle",
			Content:   "This is the third puzzle.",
			Solution:  "202",
			Subjects:  randomSubjects(),
			CreatedAt: time.Now().Add(-96 * time.Hour),
		},
		{
			ID:        4,
			Title:     "Forth Puzzle",
			Content:   "This is the Forth puzzle.",
			Solution:  "202",
			Subjects:  randomSubjects(),
			CreatedAt: time.Now().Add(-72 * time.Hour),
		},
		{
			ID:        5,
			Title:     "Fifth Puzzle",
			Content:   "This is the Fifth puzzle.",
			Solution:  "202",
			Subjects:  randomSubjects(),
			CreatedAt: time.Now().Add(-48 * time.Hour),
		},
		{
			ID:        6,
			Title:     "Sixth Puzzle",
			Content:   "This is the Sixth puzzle.",
			Solution:  "202",
			Subjects:  randomSubjects(),
			CreatedAt: time.Now().Add(-24 * time.Hour),
		},
		{
			ID:        7,
			Title:     "Demo Puzzle",
			Content:   "This is a demo puzzle.",
			Solution:  "751857",
			Subjects:  pq.StringArray{"Physics", "Math", "English"},
			CreatedAt: time.Now(),
		},
	}

	for _, puzzle := range puzzles {
		if err := db.FirstOrCreate(&puzzle, models.Puzzle{ID: puzzle.ID}).Error; err != nil {
			return err
		}
	}
	return nil
}

func SeedUserPuzzles(db *gorm.DB) error {
	now := time.Now()

	// User 1 solves puzzles 2, 3, 4 at different times
	userPuzzles := []models.UserPuzzle{
		{
			UserID:   1,
			PuzzleID: 3,
			SolvedAt: now.Add(-96 * time.Hour), // 3 days ago
		},
		{
			UserID:   1,
			PuzzleID: 4,
			SolvedAt: now.Add(-72 * time.Hour), // 2 days ago
		},
		{
			UserID:   1,
			PuzzleID: 5,
			SolvedAt: now.Add(-48 * time.Hour), // 1 day ago
		},
		{
			UserID:   1,
			PuzzleID: 6,
			SolvedAt: now.Add(-24 * time.Hour), // 1 day ago
		},
		// Add other user's puzzle solves if needed
	}

	for _, up := range userPuzzles {
		// Only create if this user+puzzle combination doesn't exist
		result := db.Where(
			models.UserPuzzle{UserID: up.UserID, PuzzleID: up.PuzzleID},
		).FirstOrCreate(&up)

		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}

func SeedUsers(db *gorm.DB) error {
	// Seed initial users
	user := models.User{
		Username: "testUser1",
		PasswordHash: func() string {
			hash, err := auth.HashPassword("pass123")
			if err != nil {
				panic(err) // Handle error appropriately
			}
			return hash
		}(),
		CreatedAt: time.Now().Add(-168 * time.Hour),
	}

	result := db.FirstOrCreate(&user, models.User{Username: user.Username})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func SeedUserSolvedPuzzles(db *gorm.DB) error {
	now := time.Now()

	// Get total puzzle count
	var totalPuzzles int64
	db.Model(&models.Puzzle{}).Count(&totalPuzzles)

	userStat := models.UserSolvedPuzzle{
		// User 1 stats (solved 3 puzzles, current streak 3)
		UserID:        1,
		SolvedPuzzles: 4,
		TotalPuzzles:  uint(totalPuzzles),
		CurrentStreak: 4,
		BestStreak:    4,
		LastSolvedAt:  now.Add(-24 * time.Hour), // Matches last solve time
	}

	if err := db.Where(
		models.UserSolvedPuzzle{UserID: userStat.UserID},
	).Assign(userStat).FirstOrCreate(&userStat).Error; err != nil {
		return err
	}
	return nil
}

func SeedData(db *gorm.DB) error {
	// Call all seed functions in order
	if err := SeedPuzzles(db); err != nil {
		return err
	}
	if err := SeedUsers(db); err != nil {
		return err
	}
	if err := SeedUserPuzzles(db); err != nil {
		return err
	}
	return SeedUserSolvedPuzzles(db)
}
