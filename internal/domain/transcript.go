package domain

import (
	"context"
)

// Transcript represents the jpcorrect.transcript table
type Transcript struct {
	TranscriptID int    `json:"transcript_id" gorm:"column:transcript_id;primaryKey;autoIncrement"`
	MistakeID    int    `json:"mistake_id" gorm:"column:mistake_id;not null"`
	Content      string `json:"content" gorm:"column:content"`
	Furigana     string `json:"furigana" gorm:"column:furigana"`
	Accent       string `json:"accent" gorm:"column:accent"`
}

// TableName overrides the default table name
func (Transcript) TableName() string {
	return "jpcorrect.transcript"
}

type TranscriptRepository interface {
	GetByID(ctx context.Context, transcriptID int) (*Transcript, error)
	GetByMistakeID(ctx context.Context, mistakeID int) (*Transcript, error)

	Create(ctx context.Context, transcript *Transcript) error
	Update(ctx context.Context, transcript *Transcript) error
	Delete(ctx context.Context, transcriptID int) error
}
