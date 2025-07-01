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

// UpdateTaskType 更新一个任务类型的名称或启用状态
func UpdateTaskType(id uuid.UUID, name string, isEnabled bool) error {
	return config.DB.Model(&model.TaskType{}).Where("id = ?", id).Updates(map[string]interface{}{
		"name":       name,
		"is_enabled": isEnabled,
	}).Error
}

// IsTaskTypeInUse 检查一个任务类型是否已被任何任务使用
func IsTaskTypeInUse(id uuid.UUID) (bool, error) {
	var count int64
	result := config.DB.Model(&model.Task{}).Where("task_type_id = ?", id).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

// DeleteTaskType 删除一个任务类型
func DeleteTaskType(id uuid.UUID) error {
	return config.DB.Where("id = ?", id).Delete(&model.TaskType{}).Error
}

// ListAllActivePeriodicTasks 获取所有启用的计划任务
func ListAllActivePeriodicTasks() ([]model.PeriodicTask, error) {
	var periodicTasks []model.PeriodicTask
	result := config.DB.Where("is_active = ?", true).Find(&periodicTasks)
	return periodicTasks, result.Error
}
