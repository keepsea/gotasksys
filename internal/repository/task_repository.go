// internal/repository/task_repository.go
package repository

import (
	"gotasksys/internal/config"
	"gotasksys/internal/model"

	"github.com/google/uuid"
)

func CreateTask(task *model.Task) error {
	result := config.DB.Create(task)
	return result.Error
}

func ListTasks() ([]model.Task, error) {
	var tasks []model.Task
	result := config.DB.Find(&tasks)
	return tasks, result.Error
}

// FindTaskByID 根据ID查找单个任务
func FindTaskByID(id uint) (model.Task, error) {
	var task model.Task
	result := config.DB.First(&task, id) // GORM的First方法会根据主键查询
	return task, result.Error
}

// UpdateTask 保存对任务的修改
func UpdateTask(task *model.Task) error {
	// Save会更新所有字段，即使它们是零值
	result := config.DB.Save(task)
	return result.Error
}

// UpdateTaskFields 更新任务的指定字段
func UpdateTaskFields(id uint, updates map[string]interface{}) error {
	result := config.DB.Model(&model.Task{}).Where("id = ?", id).Updates(updates)
	return result.Error
}

// DeleteTask 根据ID删除任务
func DeleteTask(id uint) error {
	// GORM的Delete方法会根据主键删除记录
	result := config.DB.Delete(&model.Task{}, id)
	return result.Error
}

// FindInProgressTasksByAssigneeID 根据负责人ID查找所有进行中的任务
func FindInProgressTasksByAssigneeID(assigneeID uuid.UUID) ([]model.Task, error) {
	var tasks []model.Task
	result := config.DB.Where("assignee_id = ? AND status = ?", assigneeID, "in_progress").Find(&tasks)
	return tasks, result.Error
}
