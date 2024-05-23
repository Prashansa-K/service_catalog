// internal/models/version.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type Version struct {
	ID          uint           `gorm:"primaryKey"`
	ServiceID   uint           `gorm:"not null;index"`
	Name        string         `gorm:"not null"`
	CreatedAt   time.Time      `gorm:"not null"`
	DeletedAt   gorm.DeletedAt `gorm:"default:null"`
	Description string         `gorm:"default:null"`
	Service     *Service       `gorm:"foreignKey:ServiceID;references:ID"`
}
