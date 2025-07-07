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

// GetTotalEffortOfSubtasks 获取一个父任务下所有子任务的工时总和
func GetTotalEffortOfSubtasks(parentTaskID uint) (int64, error) {
	var totalEffort int64
	// 使用 GORM 的 Select 和 Where 来构建 SUM 查询
	result := config.DB.Model(&model.Task{}).
		Where("parent_task_id = ?", parentTaskID).
		Select("COALESCE(SUM(effort), 0)"). // COALESCE 确保在没有子任务时返回0而不是NULL
		Row().
		Scan(&totalEffort)

	if result != nil {
		return 0, result
	}
	return totalEffort, nil
}

// CountIncompleteSubtasks 获取一个父任务下未完成的子任务数量
func CountIncompleteSubtasks(parentTaskID uint) (int64, error) {
	var count int64
	// 我们定义 "未完成" 的状态是不等于 'completed'
	result := config.DB.Model(&model.Task{}).
		Where("parent_task_id = ? AND status != ?", parentTaskID, "completed").
		Count(&count)

	return count, result.Error
}

// PerformanceMetrics 定义了从数据库聚合查询返回的结构
type PerformanceMetrics struct {
	AvgTimeliness    float64
	AvgQuality       float64
	AvgCollaboration float64
	AvgComplexity    float64
}

// GetPerformanceMetricsForUser 获取一个用户所有已完成任务的各项评价平均分
func GetPerformanceMetricsForUser(userID uuid.UUID) (PerformanceMetrics, error) {
	var metrics PerformanceMetrics

	// 我们使用原生SQL查询，因为JSON字段的聚合操作非常复杂，原生SQL更清晰高效
	query := `
		SELECT 
			COALESCE(AVG((evaluation->>'timeliness')::numeric), 0) as avg_timeliness,
			COALESCE(AVG((evaluation->>'quality')::numeric), 0) as avg_quality,
			COALESCE(AVG((evaluation->>'collaboration')::numeric), 0) as avg_collaboration,
			COALESCE(AVG((evaluation->>'complexity')::numeric), 0) as avg_complexity
		FROM 
			tasks
		WHERE 
			assignee_id = ? AND status = 'completed' AND evaluation IS NOT NULL;
	`

	result := config.DB.Raw(query, userID).Scan(&metrics)
	if result.Error != nil {
		return PerformanceMetrics{}, result.Error
	}

	return metrics, nil
}

// BatchUpdateSubtasksAssignee 批量更新一个主任务下，特定原负责人的所有子任务的新负责人
func BatchUpdateSubtasksAssignee(parentTaskID uint, oldAssigneeID, newAssigneeID uuid.UUID) error {
	result := config.DB.Model(&model.Task{}).
		Where("parent_task_id = ? AND assignee_id = ?", parentTaskID, oldAssigneeID).
		Update("assignee_id", newAssigneeID)

	return result.Error
}
