package domain

import (
	"context"
)

// Mistake represents the jpcorrect.mistake table
type Mistake struct {
	MistakeID     int     `json:"mistake_id" gorm:"column:mistake_id;primaryKey;autoIncrement"`
	PracticeID    int     `json:"practice_id" gorm:"column:practice_id;not null"`
	UserID        int     `json:"user_id" gorm:"column:user_id;not null"`
	StartTime     float64 `json:"start_time" gorm:"column:start_time;not null;default:0"`
	EndTime       float64 `json:"end_time" gorm:"column:end_time;not null;default:0"`
	MistakeStatus string  `json:"mistake_status" gorm:"column:mistake_status;type:jpcorrect.mistake_status;not null"`
	MistakeType   string  `json:"mistake_type" gorm:"column:mistake_type;type:jpcorrect.mistake_type;not null"`
}

// TableName overrides the default table name
func (Mistake) TableName() string {
	return "jpcorrect.mistake"
}

type MistakeRepository interface {
	GetByID(ctx context.Context, mistakeID int) (*Mistake, error)
	GetByPracticeID(ctx context.Context, practiceID int) ([]*Mistake, error)
	GetByUserID(ctx context.Context, userID int) ([]*Mistake, error)

	Create(ctx context.Context, m *Mistake) error
	Update(ctx context.Context, m *Mistake) error
	Delete(ctx context.Context, mistakeID int) error
}
