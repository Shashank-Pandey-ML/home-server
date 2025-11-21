package models

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel provides common fields for all models
type BaseModel struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// User represents a user in the system
type User struct {
	BaseModel
	Email    string `json:"email" gorm:"uniqueIndex;not null"`
	Name     string `json:"name" gorm:"not null"`
	Password string `json:"-" gorm:"not null"` // omit in JSON
	IsAdmin  bool   `json:"is_admin" gorm:"default:false"`
}

// TableName returns the table name for User model
func (User) TableName() string {
	return "users"
}

// BeforeCreate hook for User model
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// Add any pre-creation logic here
	return nil
}

// BeforeUpdate hook for User model
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	// Add any pre-update logic here
	return nil
}
