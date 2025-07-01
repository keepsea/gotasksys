// internal/repository/user_repository.go (最终锁定版)

package repository

import (
	"gotasksys/internal/config"
	"gotasksys/internal/model"

	"github.com/google/uuid"
)

// --- 已有函数 ---
func FindUserByUsername(username string) (model.User, error) {
	var user model.User
	result := config.DB.Where("username = ?", username).First(&user)
	return user, result.Error
}

func FindUserByID(id uuid.UUID) (model.User, error) {
	var user model.User
	result := config.DB.First(&user, "id = ?", id)
	return user, result.Error
}

func FindAllActiveMembers() ([]model.User, error) {
	var members []model.User
	rolesToShow := []string{"manager", "executor"}
	result := config.DB.Where("role IN (?)", rolesToShow).Find(&members)
	return members, result.Error
}

// --- 为后台管理新增的函数 ---

// ListAllUsers 获取系统内所有用户
func ListAllUsers() ([]model.User, error) {
	var users []model.User
	result := config.DB.Order("created_at asc").Find(&users)
	return users, result.Error
}

// CreateUser 创建一个新用户 (复用已有的CreateUser)
func CreateUser(user *model.User) error {
	result := config.DB.Create(user)
	return result.Error
}

// UpdateUserRole 更新指定用户的角色
func UpdateUserRole(userID uuid.UUID, newRole string) error {
	return config.DB.Model(&model.User{}).Where("id = ?", userID).Update("role", newRole).Error
}

// UpdateUserPassword 更新指定用户的密码
func UpdateUserPassword(userID uuid.UUID, newPasswordHash string) error {
	return config.DB.Model(&model.User{}).Where("id = ?", userID).Update("password_hash", newPasswordHash).Error
}

// DeleteUser 删除一个指定用户
func DeleteUser(userID uuid.UUID) error {
	// 注意：这是一个物理删除
	return config.DB.Where("id = ?", userID).Delete(&model.User{}).Error
}
