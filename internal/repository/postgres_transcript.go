package repository

import (
	"context"
	"errors"

	"jpcorrect-backend/internal/domain"

	"gorm.io/gorm"
)

type postgresTranscriptRepository struct {
	db *gorm.DB
}

func NewPostgresTranscript(conn *Connection) domain.TranscriptRepository {
	return &postgresTranscriptRepository{db: conn.DB}
}

func (p *postgresTranscriptRepository) GetByID(ctx context.Context, transcriptID int) (*domain.Transcript, error) {
	var transcript domain.Transcript
	err := p.db.WithContext(ctx).Where("transcript_id = ?", transcriptID).First(&transcript).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &transcript, nil
}

func (p *postgresTranscriptRepository) GetByMistakeID(ctx context.Context, mistakeID int) (*domain.Transcript, error) {
	var transcript domain.Transcript
	err := p.db.WithContext(ctx).Where("mistake_id = ?", mistakeID).First(&transcript).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &transcript, nil
}

func (p *postgresTranscriptRepository) Create(ctx context.Context, transcript *domain.Transcript) error {
	return p.db.WithContext(ctx).Create(transcript).Error
}

func (p *postgresTranscriptRepository) Update(ctx context.Context, transcript *domain.Transcript) error {
	result := p.db.WithContext(ctx).Model(transcript).Where("transcript_id = ?", transcript.TranscriptID).Updates(map[string]interface{}{
		"mistake_id": transcript.MistakeID,
		"content":    transcript.Content,
		"furigana":   transcript.Furigana,
		"accent":     transcript.Accent,
	})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (p *postgresTranscriptRepository) Delete(ctx context.Context, transcriptID int) error {
	result := p.db.WithContext(ctx).Where("transcript_id = ?", transcriptID).Delete(&domain.Transcript{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
