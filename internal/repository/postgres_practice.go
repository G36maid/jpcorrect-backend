package repository

import (
	"context"
	"errors"

	"jpcorrect-backend/internal/domain"

	"gorm.io/gorm"
)

type postgresPracticeRepository struct {
	db *gorm.DB
}

func NewPostgresPractice(conn *Connection) domain.PracticeRepository {
	return &postgresPracticeRepository{db: conn.DB}
}

func (p *postgresPracticeRepository) GetByID(ctx context.Context, practiceID int) (*domain.Practice, error) {
	var practice domain.Practice
	err := p.db.WithContext(ctx).Where("practice_id = ?", practiceID).First(&practice).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &practice, nil
}

func (p *postgresPracticeRepository) GetByUserID(ctx context.Context, userID int) ([]*domain.Practice, error) {
	var practices []*domain.Practice
	err := p.db.WithContext(ctx).Where("user_id = ?", userID).Find(&practices).Error
	if err != nil {
		return nil, err
	}
	if len(practices) == 0 {
		return nil, domain.ErrNotFound
	}
	return practices, nil
}

func (p *postgresPracticeRepository) Create(ctx context.Context, practice *domain.Practice) error {
	return p.db.WithContext(ctx).Create(practice).Error
}

func (p *postgresPracticeRepository) Update(ctx context.Context, practice *domain.Practice) error {
	result := p.db.WithContext(ctx).Model(practice).Where("practice_id = ?", practice.PracticeID).Updates(map[string]interface{}{
		"user_id":  practice.UserID,
		"date":     practice.Date,
		"duration": practice.Duration,
	})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (p *postgresPracticeRepository) Delete(ctx context.Context, practiceID int) error {
	result := p.db.WithContext(ctx).Where("practice_id = ?", practiceID).Delete(&domain.Practice{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
