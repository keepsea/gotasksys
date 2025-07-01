// internal/service/admin_service.go
package service

import (
	"errors"
	"gotasksys/internal/model"
	"gotasksys/internal/repository"
	"gotasksys/pkg/utils"

	"github.com/google/uuid"
)

// 我们暂时把service和repository的函数一一对应
func ListTaskTypesService() ([]model.TaskType, error) {
	return repository.ListTaskTypes()
}

func CreateTaskTypeService(name string) (model.TaskType, error) {
	taskType := model.TaskType{
		Name:      name,
		IsEnabled: true,
	}
	err := repository.CreateTaskType(&taskType)
	return taskType, err
}

// --- 用户管理 ---
func ListUsersService() ([]model.User, error) {
	return repository.ListAllUsers()
}

func UpdateUserRoleService(userID uuid.UUID, newRole string) error {
	// 在这里可以加入不允许将最后一个admin降级的逻辑等，V1.0暂时简化
	return repository.UpdateUserRole(userID, newRole)
}

func DeleteUserService(userID uuid.UUID) error {
	// 在这里可以加入不允许删除自己的逻辑等，V1.0暂时简化
	return repository.DeleteUser(userID)
}

func ResetPasswordService(userID uuid.UUID, newPassword string) error {
	if len(newPassword) < 12 {
		return errors.New("new password must be at least 12 characters long")
	}
	newPasswordHash, err := utils.HashPassword(newPassword)
	if err != nil {
		return errors.New("failed to hash new password")
	}
	return repository.UpdateUserPassword(userID, newPasswordHash)
}

// UpdateTaskTypeService 封装了更新任务类型的业务逻辑
func UpdateTaskTypeService(id uuid.UUID, name string, isEnabled bool) error {
	// 此处可添加更多业务逻辑，如名称是否重复等
	return repository.UpdateTaskType(id, name, isEnabled)
}

// DeleteTaskTypeService 封装了删除任务类型的业务逻辑
func DeleteTaskTypeService(id uuid.UUID) error {
	// 核心业务规则：如果一个类型正在被使用，则不允许删除
	inUse, err := repository.IsTaskTypeInUse(id)
	if err != nil {
		return err
	}
	if inUse {
		return errors.New("cannot delete task type: it is currently in use by one or more tasks")
	}
	return repository.DeleteTaskType(id)
}
