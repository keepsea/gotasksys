// internal/repository/admin_repository.go
package repository

import (
	"gotasksys/internal/config"
	"gotasksys/internal/model"

	"github.com/google/uuid"
)

func ListTaskTypes() ([]model.TaskType, error) {
	var types []model.TaskType
	result := config.DB.Order("created_at asc").Find(&types)
	return types, result.Error
}

func CreateTaskType(taskType *model.TaskType) error {
	result := config.DB.Create(taskType)
	return result.Error
}

func UpdateTaskType(id uuid.UUID, name string, isEnabled bool) error {
	result := config.DB.Model(&model.TaskType{}).Where("id = ?", id).Updates(map[string]interface{}{
		"name":       name,
		"is_enabled": isEnabled,
	})
	return result.Error
}
