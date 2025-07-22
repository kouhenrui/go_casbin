package audit

import (
	"time"

	"gorm.io/gorm"
)

type AuditLog struct {
	gorm.Model
	Action    string    `gorm:"not null"`
	TableName string    `gorm:"not null"`
	RecordID  uint      `gorm:"not null"`
	OldData   string    `gorm:"type:json"`
	NewData   string    `gorm:"type:json"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}