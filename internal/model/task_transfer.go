// internal/model/task_transfer.go
package model

import (
	"time"

	"github.com/google/uuid"
)

type TaskTransfer struct {
	ID                     uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TaskID                 uint      `gorm:"not null"`
	FromUserID             uuid.UUID `gorm:"not null"`
	ToUserID               uuid.UUID `gorm:"not null"`
	EffortSpentByInitiator int       `gorm:"not null"`
	Status                 string    `gorm:"type:varchar(50);not null;default:'pending'"`
	CreatedAt              time.Time
	UpdatedAt              time.Time
}
