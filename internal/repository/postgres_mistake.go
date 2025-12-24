package repository

import (
	"context"
	"errors"

	"jpcorrect-backend/internal/domain"

	"gorm.io/gorm"
)

type postgresMistakeRepository struct {
	db *gorm.DB
}

func NewPostgresMistake(conn *Connection) domain.MistakeRepository {
	return &postgresMistakeRepository{db: conn.DB}
}

func (p *postgresMistakeRepository) GetByID(ctx context.Context, mistakeID int) (*domain.Mistake, error) {
	var mistake domain.Mistake
	err := p.db.WithContext(ctx).Where("mistake_id = ?", mistakeID).First(&mistake).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &mistake, nil
}

func (p *postgresMistakeRepository) GetByPracticeID(ctx context.Context, practiceID int) ([]*domain.Mistake, error) {
	var mistakes []*domain.Mistake
	err := p.db.WithContext(ctx).Where("practice_id = ?", practiceID).Find(&mistakes).Error
	if err != nil {
		return nil, err
	}
	if len(mistakes) == 0 {
		return nil, domain.ErrNotFound
	}
	return mistakes, nil
}

func (p *postgresMistakeRepository) GetByUserID(ctx context.Context, userID int) ([]*domain.Mistake, error) {
	var mistakes []*domain.Mistake
	err := p.db.WithContext(ctx).Where("user_id = ?", userID).Find(&mistakes).Error
	if err != nil {
		return nil, err
	}
	if len(mistakes) == 0 {
		return nil, domain.ErrNotFound
	}
	return mistakes, nil
}

func (p *postgresMistakeRepository) Create(ctx context.Context, m *domain.Mistake) error {
	return p.db.WithContext(ctx).Create(m).Error
}

func (p *postgresMistakeRepository) Update(ctx context.Context, m *domain.Mistake) error {
	result := p.db.WithContext(ctx).Model(m).Where("mistake_id = ?", m.MistakeID).Updates(map[string]interface{}{
		"practice_id":    m.PracticeID,
		"user_id":        m.UserID,
		"start_time":     m.StartTime,
		"end_time":       m.EndTime,
		"mistake_status": m.MistakeStatus,
		"mistake_type":   m.MistakeType,
	})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (p *postgresMistakeRepository) Delete(ctx context.Context, mistakeID int) error {
	result := p.db.WithContext(ctx).Where("mistake_id = ?", mistakeID).Delete(&domain.Mistake{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
