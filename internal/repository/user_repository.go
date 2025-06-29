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

// FindAllActiveMembers 查找所有需要被看板管理的活跃成员
func FindAllActiveMembers() ([]model.User, error) {
	var users []model.User
	// 我们假设看板只需要展示 manager 和 member
	result := config.DB.Where("role = ? OR role = ?", "manager", "member").Find(&users)
	return users, result.Error
}
