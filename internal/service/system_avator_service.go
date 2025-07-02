// internal/service/system_avatar_service.go (最终锁定版)
package service

import (
	"errors"
	"gotasksys/internal/model"
	"gotasksys/internal/repository"

	"github.com/google/uuid"
)

// ListSystemAvatarsService 根据参数决定是获取全部还是仅获取启用的头像
func ListSystemAvatarsService(all bool) ([]model.SystemAvatar, error) {
	return repository.ListSystemAvatars(all)
}

// CreateSystemAvatarService 创建一个新的系统头像
// 【修正】参数变为(model.SystemAvatar)
func CreateSystemAvatarService(avatar model.SystemAvatar) (model.SystemAvatar, error) {
	// 新头像默认为启用状态
	avatar.IsActive = true
	err := repository.CreateSystemAvatar(&avatar)
	return avatar, err
}

// UpdateSystemAvatarService 更新一个系统头像
// 【修正】参数变为(uuid.UUID, model.SystemAvatar)
func UpdateSystemAvatarService(id uuid.UUID, input model.SystemAvatar) (model.SystemAvatar, error) {
	avatar, err := repository.FindSystemAvatarByID(id)
	if err != nil {
		return model.SystemAvatar{}, errors.New("system avatar not found")
	}

	// 更新字段
	avatar.URL = input.URL
	avatar.Description = input.Description
	avatar.IsActive = input.IsActive

	err = repository.UpdateSystemAvatar(&avatar)
	return avatar, err
}

// DeleteSystemAvatarService 删除一个系统头像
func DeleteSystemAvatarService(id uuid.UUID) error {
	// 在这里可以增加业务逻辑，例如检查该头像是否被用户默认使用等
	return repository.DeleteSystemAvatar(id)
}
