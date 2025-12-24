package repository

import (
	"context"
	"errors"

	"jpcorrect-backend/internal/domain"

	"gorm.io/gorm"
)

type postgresUserRepository struct {
	db *gorm.DB
}

func NewPostgresUser(conn *Connection) domain.UserRepository {
	return &postgresUserRepository{db: conn.DB}
}

func (u *postgresUserRepository) GetByID(ctx context.Context, userID int) (*domain.User, error) {
	var user domain.User
	err := u.db.WithContext(ctx).Where("user_id = ?", userID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (u *postgresUserRepository) GetByName(ctx context.Context, name string) ([]*domain.User, error) {
	var users []*domain.User
	err := u.db.WithContext(ctx).Where("name = ?", name).Find(&users).Error
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, domain.ErrNotFound
	}
	return users, nil
}

func (u *postgresUserRepository) Create(ctx context.Context, user *domain.User) error {
	return u.db.WithContext(ctx).Create(user).Error
}

func (u *postgresUserRepository) Update(ctx context.Context, user *domain.User) error {
	result := u.db.WithContext(ctx).Model(user).Where("user_id = ?", user.UserID).Updates(map[string]interface{}{
		"name": user.Name,
	})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *postgresUserRepository) Delete(ctx context.Context, userID int) error {
	result := u.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&domain.User{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
