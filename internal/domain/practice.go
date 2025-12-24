package domain

import (
	"context"
	"time"
)

// Practice represents the jpcorrect.practice table
type Practice struct {
	PracticeID int       `json:"practice_id" gorm:"column:practice_id;primaryKey;autoIncrement"`
	UserID     int       `json:"user_id" gorm:"column:user_id"`
	Date       time.Time `json:"date" gorm:"column:date;type:date"`
	Duration   float64   `json:"duration" gorm:"column:duration"`
}

// TableName overrides the default table name
func (Practice) TableName() string {
	return "jpcorrect.practice"
}

type PracticeRepository interface {
	GetByID(ctx context.Context, practiceID int) (*Practice, error)
	GetByUserID(ctx context.Context, userID int) ([]*Practice, error)

	Create(ctx context.Context, practice *Practice) error
	Update(ctx context.Context, practice *Practice) error
	Delete(ctx context.Context, practiceID int) error
}
