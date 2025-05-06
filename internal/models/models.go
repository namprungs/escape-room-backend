package models

import (
	"time"

	"github.com/lib/pq"
)

// In models/user.go
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"unique" json:"username"`
	Password     string    `gorm:"-" json:"password"` // Only for input, not stored
	PasswordHash string    `json:"-"`                 // Only stored in DB
	CreatedAt    time.Time `json:"created_at"`
}

type Puzzle struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `json:"title"`
	Content   string         `json:"content"`
	Solution  string         `json:"-"`
	Subjects  pq.StringArray `json:"subjects" gorm:"type:text[]"`
	CreatedAt time.Time      `json:"created_at"`
}

type UserPuzzle struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	UserID   uint      `gorm:"index" json:"user_id"`
	PuzzleID uint      `gorm:"index" json:"puzzle_id"`
	SolvedAt time.Time `json:"solved_at"`
	Puzzle   Puzzle    `gorm:"foreignKey:PuzzleID" json:"-"` // For Preload
}

type UserSolvedPuzzle struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"uniqueIndex" json:"user_id"` // One row per user
	SolvedPuzzles uint      `json:"solved_puzzles"`
	TotalPuzzles  uint      `json:"total_puzzles"`
	CurrentStreak uint      `json:"current_streak"`
	BestStreak    uint      `json:"best_streak"`
	LastSolvedAt  time.Time `json:"last_solved_at"`
}
