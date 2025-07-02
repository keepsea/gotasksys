// internal/repository/system_avatar_repository.go
package repository

import (
	"gotasksys/internal/config"
	"gotasksys/internal/model"

	"github.com/google/uuid"
)

// ListSystemAvatars 获取系统内置头像列表
// all: true - 获取所有头像 (供管理员使用)
// all: false - 只获取已启用的头像 (供普通用户选择)
func ListSystemAvatars(all bool) ([]model.SystemAvatar, error) {
	var avatars []model.SystemAvatar
	db := config.DB
	if !all {
		db = db.Where("is_active = ?", true)
	}
	err := db.Order("created_at asc").Find(&avatars).Error
	return avatars, err
}

// CreateSystemAvatar 在数据库中创建一条新的头像记录
func CreateSystemAvatar(avatar *model.SystemAvatar) error {
	return config.DB.Create(avatar).Error
}

// FindSystemAvatarByID 根据ID查找一个内置头像
func FindSystemAvatarByID(id uuid.UUID) (model.SystemAvatar, error) {
	var avatar model.SystemAvatar
	err := config.DB.First(&avatar, "id = ?", id).Error
	return avatar, err
}

// UpdateSystemAvatar 更新一个内置头像的信息
func UpdateSystemAvatar(avatar *model.SystemAvatar) error {
	return config.DB.Save(avatar).Error
}

// DeleteSystemAvatar 从数据库中删除一个内置头像
func DeleteSystemAvatar(id uuid.UUID) error {
	return config.DB.Where("id = ?", id).Delete(&model.SystemAvatar{}).Error
}
