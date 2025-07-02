// internal/service/profile_service.go
package service

import (
	"errors"
	"gotasksys/internal/repository"
	"gotasksys/pkg/utils"

	"github.com/google/uuid"
)

// UpdateMyProfileService 处理用户更新自己姓名的头像的逻辑
func UpdateMyProfileService(userID uuid.UUID, realName, avatar, email, team string) error {
	user, err := repository.FindUserByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	updates := make(map[string]interface{})
	if realName != "" {
		updates["real_name"] = realName
	}
	if email != "" {
		updates["email"] = email
	}
	if team != "" {
		updates["team"] = team
	}

	// 核心业务规则：禁止任何人（包括自己）通过此接口修改system_admin的头像
	if user.Role == "system_admin" {
		if avatar != "" && avatar != user.Avatar {
			return errors.New("system admin avatar cannot be changed")
		}
	} else {
		updates["avatar"] = avatar
	}

	return repository.UpdateUserProfile(userID, updates)
}

// ChangeMyPasswordService 处理用户修改自己密码的逻辑
func ChangeMyPasswordService(userID uuid.UUID, oldPassword, newPassword string) error {
	if len(newPassword) < 12 {
		return errors.New("new password does not meet complexity requirements")
	}
	user, err := repository.FindUserByID(userID)
	if err != nil {
		return errors.New("user not found")
	}
	// 验证旧密码是否正确
	if !utils.CheckPasswordHash(oldPassword, user.PasswordHash) {
		return errors.New("old password is incorrect")
	}
	newHash, err := utils.HashPassword(newPassword)
	if err != nil {
		return errors.New("failed to hash new password")
	}
	return repository.UpdateUserPassword(userID, newHash)
}
