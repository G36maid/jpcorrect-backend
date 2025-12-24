package repository

import (
	"gorm.io/gorm"
)

// Connection wraps *gorm.DB to provide a consistent interface
// This allows us to maintain the same repository pattern while using GORM
type Connection struct {
	DB *gorm.DB
}

// NewConnection creates a new Connection wrapper around a gorm.DB instance
func NewConnection(db *gorm.DB) *Connection {
	return &Connection{DB: db}
}
