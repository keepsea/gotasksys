// internal/model/task_type.go

package model

import (
	"time"

	"github.com/google/uuid"
)

// TaskType 定义了任务类型的数据结构
type TaskType struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name      string    `gorm:"type:varchar(255);unique_not_null" json:"name"`
	IsEnabled bool      `gorm:"default:true" json:"is_enabled"`
	CreatedAt time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null;default:now()" json:"updated_at"`
}
