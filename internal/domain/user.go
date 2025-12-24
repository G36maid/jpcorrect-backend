package domain

import (
	"context"
)

// User represents the jpcorrect.user table
type User struct {
	UserID int    `json:"user_id" gorm:"column:user_id;primaryKey;autoIncrement"`
	Name   string `json:"name" gorm:"column:name;not null"`
}

// TableName overrides the default table name
func (User) TableName() string {
	return "jpcorrect.user"
}

type UserRepository interface {
	GetByID(ctx context.Context, userID int) (*User, error)
	GetByName(ctx context.Context, name string) ([]*User, error)

	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, userID int) error
}
