// internal/repository/task_repository.go (最终完整版)
package repository

import (
	"gotasksys/internal/config"
	"gotasksys/internal/model"

	"github.com/google/uuid"
)

// --- 已有函数 ---
func CreateTask(task *model.Task) error {
	result := config.DB.Create(task)
	return result.Error
}

func FindTaskByID(id uint) (model.Task, error) {
	var task model.Task
	result := config.DB.Preload("Creator").Preload("Assignee").First(&task, id)
	return task, result.Error
}

func UpdateTask(task *model.Task) error {
	result := config.DB.Save(task)
	return result.Error
}

func UpdateTaskFields(id uint, updates map[string]interface{}) error {
	result := config.DB.Model(&model.Task{}).Where("id = ?", id).Updates(updates)
	return result.Error
}

func DeleteTask(id uint) error {
	result := config.DB.Delete(&model.Task{}, id)
	return result.Error
}

// --- 为人员看板新增的函数 (之前被遗漏) ---
// FindInProgressTasksByAssigneeID 根据负责人ID查找所有进行中的任务
func FindInProgressTasksByAssigneeID(assigneeID uuid.UUID) ([]model.Task, error) {
	var tasks []model.Task
	result := config.DB.Where("assignee_id = ? AND status = ?", assigneeID, "in_progress").Find(&tasks)
	return tasks, result.Error
}

// --- 为列表精细化查询新增的函数 ---

// ListAllTasks 获取所有任务 (供 admin/manager 使用)
func ListAllTasks() ([]model.Task, error) {
	var tasks []model.Task
	result := config.DB.Preload("Creator").Preload("Assignee").Order("created_at desc").Find(&tasks)
	return tasks, result.Error
}

// ListTasksForExecutor 获取执行者能看到的任务
func ListTasksForExecutor(executorID uuid.UUID) ([]model.Task, error) {
	var tasks []model.Task
	result := config.DB.Preload("Creator").Preload("Assignee").
		Where("status = ?", "in_pool").
		Or("assignee_id = ?", executorID).
		Order("created_at desc").
		Find(&tasks)
	return tasks, result.Error
}

// ListTasksForCreator 获取创建者能看到的任务
func ListTasksForCreator(creatorID uuid.UUID) ([]model.Task, error) {
	var tasks []model.Task
	publicStatuses := []string{"in_pool", "in_progress", "pending_evaluation", "completed"}

	result := config.DB.Preload("Creator").Preload("Assignee").
		Where("status IN (?)", publicStatuses).
		Or("creator_id = ? AND status IN (?)", creatorID, []string{"pending_review", "rejected"}).
		Order("created_at desc").
		Find(&tasks)
	return tasks, result.Error
}
