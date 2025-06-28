// internal/model/user.go

package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Username     string    `gorm:"type:varchar(255);unique_not_null" json:"username"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"` // json:"-" 表示在返回JSON时不包含此字段
	RealName     string    `gorm:"type:varchar(255);not null" json:"real_name"`
	Role         string    `gorm:"type:varchar(50);not null" json:"role"`
	Avatar       string    `gorm:"type:varchar(255)" json:"avatar,omitempty"` // omitempty 表示如果字段为空则在JSON中忽略
	CreatedAt    time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt    time.Time `gorm:"not null;default:now()" json:"updated_at"`
}
