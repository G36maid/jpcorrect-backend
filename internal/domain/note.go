package domain

import (
	"context"
)

// Note represents the jpcorrect.note table
type Note struct {
	NoteID     int    `json:"note_id" gorm:"column:note_id;primaryKey;autoIncrement"`
	PracticeID int    `json:"practice_id" gorm:"column:practice_id;not null"`
	Content    string `json:"content" gorm:"column:content"`
}

// TableName overrides the default table name
func (Note) TableName() string {
	return "jpcorrect.note"
}

type NoteRepository interface {
	GetByID(ctx context.Context, noteID int) (*Note, error)
	GetByPracticeID(ctx context.Context, practiceID int) (*Note, error)

	Create(ctx context.Context, note *Note) error
	Update(ctx context.Context, note *Note) error
	Delete(ctx context.Context, noteID int) error
}
