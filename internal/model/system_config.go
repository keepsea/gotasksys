// internal/model/system_config.go
package model

import "time"

type SystemConfig struct {
	ConfigKey   string `gorm:"primaryKey"`
	ConfigValue string `gorm:"not null"`
	Description string
	UpdatedAt   time.Time `gorm:"not null;default:now()"`
}
