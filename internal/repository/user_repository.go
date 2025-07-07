// internal/repository/user_repository.go (最终锁定版)

package repository

import (
	"gotasksys/internal/config"
	"gotasksys/internal/model"

	"github.com/google/uuid"
)

// 核心查询
func FindUserByID(id uuid.UUID) (model.User, error) {
	var user model.User
	result := config.DB.First(&user, "id = ?", id)
	return user, result.Error
}

// 根据用户名查找用户
func FindUserByUsername(username string) (model.User, error) {
	var user model.User
	result := config.DB.Where("username = ?", username).First(&user)
	return user, result.Error
}

// 查询所有用户
func ListAllUsers() ([]model.User, error) {
	var users []model.User
	result := config.DB.Order("created_at asc").Find(&users)
	return users, result.Error
}

// 查询所有活跃成员(manager和executor角色)
func FindAllActiveMembers() ([]model.User, error) {
	var members []model.User
	rolesToShow := []string{"manager", "executor"}
	result := config.DB.Where("role IN (?)", rolesToShow).Find(&members)
	return members, result.Error
}

// CreateUser 创建一个新用户
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

// UpdateUserProfile 更新用户的非敏感信息
func UpdateUserProfile(userID uuid.UUID, updates map[string]interface{}) error {
	// 如果没有任何需要更新的字段，则直接返回成功
	if len(updates) == 0 {
		return nil
	}
	// GORM的Updates方法，可以安全地只更新map中存在的字段
	return config.DB.Model(&model.User{}).Where("id = ?", userID).Updates(updates).Error
}

// DeleteUser 删除一个指定用户
func DeleteUser(userID uuid.UUID) error {
	// 注意：这是一个物理删除
	return config.DB.Where("id = ?", userID).Delete(&model.User{}).Error
}

// HasUnfinishedTasks 检查一个用户是否还有未完成的任务
func HasUnfinishedTasks(userID uuid.UUID) (bool, error) {
	var count int64
	// 查询该用户作为负责人(assignee)的、状态不为'completed'的任务数量
	result := config.DB.Model(&model.Task{}).
		Where("assignee_id = ? AND status != ?", userID, "completed").
		Count(&count)

	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil
}
