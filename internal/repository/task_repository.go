// internal/repository/task_repository.go
package repository

import (
	"gotasksys/internal/config"
	"gotasksys/internal/model"
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

// DeleteTask 根据ID删除任务
func DeleteTask(id uint) error {
	// GORM的Delete方法会根据主键删除记录
	result := config.DB.Delete(&model.Task{}, id)
	return result.Error
}
