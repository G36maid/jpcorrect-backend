package repository

import (
	"context"
	"errors"

	"jpcorrect-backend/internal/domain"

	"gorm.io/gorm"
)

type postgresNoteRepository struct {
	db *gorm.DB
}

func NewPostgresNote(conn *Connection) domain.NoteRepository {
	return &postgresNoteRepository{db: conn.DB}
}

func (p *postgresNoteRepository) GetByID(ctx context.Context, noteID int) (*domain.Note, error) {
	var note domain.Note
	err := p.db.WithContext(ctx).Where("note_id = ?", noteID).First(&note).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &note, nil
}

func (p *postgresNoteRepository) GetByPracticeID(ctx context.Context, practiceID int) (*domain.Note, error) {
	var note domain.Note
	err := p.db.WithContext(ctx).Where("practice_id = ?", practiceID).First(&note).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &note, nil
}

func (p *postgresNoteRepository) Create(ctx context.Context, note *domain.Note) error {
	return p.db.WithContext(ctx).Create(note).Error
}

func (p *postgresNoteRepository) Update(ctx context.Context, note *domain.Note) error {
	result := p.db.WithContext(ctx).Model(note).Where("note_id = ?", note.NoteID).Updates(map[string]interface{}{
		"practice_id": note.PracticeID,
		"content":     note.Content,
	})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (p *postgresNoteRepository) Delete(ctx context.Context, noteID int) error {
	result := p.db.WithContext(ctx).Where("note_id = ?", noteID).Delete(&domain.Note{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
