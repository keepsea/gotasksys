// internal/model/leave.go
package model

import (
	"time"

	"github.com/google/uuid"
)

type Leave struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"not null;index" json:"user_id"`
	StartDate time.Time `gorm:"type:date;not null" json:"start_date"`
	EndDate   time.Time `gorm:"type:date;not null" json:"end_date"`
	Reason    string    `gorm:"type:text" json:"reason"`
	Status    string    `gorm:"type:varchar(50);not null;default:'approved'" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
