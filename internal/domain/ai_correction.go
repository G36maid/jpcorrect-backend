package domain

import (
	"context"
)

// AICorrection represents the jpcorrect.ai_correction table
type AICorrection struct {
	AICorrectionID int    `json:"ai_correction_id" gorm:"column:ai_correction_id;primaryKey;autoIncrement"`
	MistakeID      int    `json:"mistake_id" gorm:"column:mistake_id;not null"`
	Content        string `json:"content" gorm:"column:content"`
}

// TableName overrides the default table name
func (AICorrection) TableName() string {
	return "jpcorrect.ai_correction"
}

type AICorrectionRepository interface {
	GetByID(ctx context.Context, aiCorrectionID int) (*AICorrection, error)
	GetByMistakeID(ctx context.Context, mistakeID int) (*AICorrection, error)

	Create(ctx context.Context, aiCorrection *AICorrection) error
	Update(ctx context.Context, aiCorrection *AICorrection) error
	Delete(ctx context.Context, aiCorrectionID int) error
}
