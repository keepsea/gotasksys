// internal/model/system_avatar.go
package model

import (
	"time"

	"github.com/google/uuid"
)

// SystemAvatar 定义了内置头像库的数据结构
type SystemAvatar struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	URL         string    `gorm:"type:varchar(255);not null;unique" json:"url"`
	Description string    `gorm:"type:varchar(255)" json:"description"`
	IsActive    bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `gorm:"not null;default:now()"`
}
