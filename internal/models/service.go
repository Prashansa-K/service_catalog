// internal/models/service.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type Service struct {
	ID           uint           `gorm:"primaryKey"`
	Name         string         `gorm:"not null;unique"`
	Description  string         `gorm:"type:text"`
	CreatedAt    time.Time      `gorm:"not null"`
	DeletedAt    gorm.DeletedAt `gorm:"default:null"`
	VersionCount int            `gorm:"default:0"`
	Versions     []Version      `gorm:"foreignKey:ServiceID;references:ID"`
}
