// internal/model/periodic_task.go
package model

import (
	"time"

	"github.com/google/uuid"
)

type PeriodicTask struct {
	ID                uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Title             string     `gorm:"type:varchar(255);not null" json:"title"`
	Description       string     `gorm:"type:text" json:"description"`
	CronExpression    string     `gorm:"type:varchar(255);not null" json:"cron_expression"`
	DefaultAssigneeID *uuid.UUID `json:"default_assignee_id"`
	DefaultEffort     int        `gorm:"not null" json:"default_effort"`
	DefaultPriority   string     `gorm:"type:varchar(50);not null" json:"default_priority"`
	DefaultTaskTypeID *uuid.UUID `json:"default_task_type_id"`
	IsActive          bool       `gorm:"not null;default:true" json:"is_active"`
	CreatedByID       uuid.UUID  `gorm:"not null" json:"created_by_id"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}
