// internal/model/user.go

package model

import (
	"encoding/json"
	"gotasksys/internal/config"
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
	Email        string    `gorm:"type:varchar(255);unique" json:"email,omitempty"`
	Team         string    `gorm:"type:varchar(255)" json:"team,omitempty"`
	CreatedAt    time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt    time.Time `gorm:"not null;default:now()" json:"updated_at"`
}

// MarshalJSON 自定义User结构体的JSON序列化方法
// 当任何时候User对象被转为JSON（如API返回时），此方法会被自动调用
func (u User) MarshalJSON() ([]byte, error) {
	type Alias User // 创建一个别名类型，防止在方法内再次调用MarshalJSON导致无限递归

	// 如果从数据库取出的头像为空字符串，则根据角色赋予默认头像
	avatar := u.Avatar
	if avatar == "" {
		if u.Role == "system_admin" {
			avatar = config.DefaultAvatarAdmin
		} else {
			avatar = config.DefaultAvatarUser
		}
	}

	return json.Marshal(&struct {
		// 覆盖原始的Avatar字段，确保即使原始值为空，最终JSON也有值
		Avatar string `json:"avatar"`
		// 嵌入原始的User所有可序列化字段
		*Alias
	}{
		Alias:  (*Alias)(&u),
		Avatar: avatar,
	})
}
