package models

import "time"

type Log struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	JobID     uint `gorm:"index"`
	RunID     uint `gorm:"index"`
	Type      string
	Message   string
}
