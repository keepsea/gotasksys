// internal/repository/user_repository.go
package repository

import (
	"gotasksys/internal/config"
	"gotasksys/internal/model"

	"github.com/google/uuid"
)

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

func CreateUser(user *model.User) error {
	result := config.DB.Create(user)
	return result.Error
}

// --- 【V1.2 最终修正】获取所有活跃的团队成员 ---
// FindAllActiveMembers 获取所有应在人员看板上展示的成员 (manager 和 executor)
func FindAllActiveMembers() ([]model.User, error) {
	var members []model.User
	// 定义需要展示在看板上的角色
	rolesToShow := []string{"manager", "executor"}

	// 查询所有角色为 manager 或 executor 的用户
	result := config.DB.Where("role IN (?)", rolesToShow).Find(&members)
	if result.Error != nil {
		return nil, result.Error
	}
	return members, nil
}
