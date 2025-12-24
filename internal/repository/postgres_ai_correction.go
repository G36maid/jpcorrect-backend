package repository

import (
	"context"
	"errors"

	"jpcorrect-backend/internal/domain"

	"gorm.io/gorm"
)

type postgresAICorrectionRepository struct {
	db *gorm.DB
}

func NewPostgresAICorrection(conn *Connection) domain.AICorrectionRepository {
	return &postgresAICorrectionRepository{db: conn.DB}
}

func (p *postgresAICorrectionRepository) GetByID(ctx context.Context, aiCorrectionID int) (*domain.AICorrection, error) {
	var aiCorrection domain.AICorrection
	err := p.db.WithContext(ctx).Where("ai_correction_id = ?", aiCorrectionID).First(&aiCorrection).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &aiCorrection, nil
}

func (p *postgresAICorrectionRepository) GetByMistakeID(ctx context.Context, mistakeID int) (*domain.AICorrection, error) {
	var aiCorrection domain.AICorrection
	err := p.db.WithContext(ctx).Where("mistake_id = ?", mistakeID).First(&aiCorrection).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &aiCorrection, nil
}

func (p *postgresAICorrectionRepository) Create(ctx context.Context, aiCorrection *domain.AICorrection) error {
	return p.db.WithContext(ctx).Create(aiCorrection).Error
}

func (p *postgresAICorrectionRepository) Update(ctx context.Context, aiCorrection *domain.AICorrection) error {
	result := p.db.WithContext(ctx).Model(aiCorrection).Where("ai_correction_id = ?", aiCorrection.AICorrectionID).Updates(map[string]interface{}{
		"mistake_id": aiCorrection.MistakeID,
		"content":    aiCorrection.Content,
	})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (p *postgresAICorrectionRepository) Delete(ctx context.Context, aiCorrectionID int) error {
	result := p.db.WithContext(ctx).Where("ai_correction_id = ?", aiCorrectionID).Delete(&domain.AICorrection{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
