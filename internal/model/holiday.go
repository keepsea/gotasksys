// internal/model/holiday.go
package model

import (
	"time"

	"github.com/google/uuid"
)

type Holiday struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	HolidayDate time.Time `gorm:"type:date;not null;unique"`
	Description string    `gorm:"type:varchar(255)"`
	CreatedAt   time.Time
}
