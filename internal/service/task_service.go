// internal/service/task_service.go

package service

import (
	"errors"
	"gotasksys/internal/model"
	"gotasksys/internal/repository"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// CreateTaskService 封装了创建任务的业务逻辑
func CreateTaskService(input model.Task) (model.Task, error) {
	// 业务规则1：新任务的初始状态总是 "pending_review"
	input.Status = "pending_review"

	// 业务规则2：原始工作量等于预估工作量
	input.OriginalEffort = input.Effort

	// 调用仓储层将任务存入数据库
	err := repository.CreateTask(&input)
	if err != nil {
		return model.Task{}, err
	}

	return input, nil
}

// ListTasksService 封装了获取任务列表的业务逻辑
func ListTasksService() ([]model.Task, error) {
	// 目前没有复杂逻辑，直接调用仓储层
	// 未来可以在这里加入权限校验、筛选等逻辑
	tasks, err := repository.ListTasks()
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// GetTaskByIDService 封装了根据ID获取任务的业务逻辑
func GetTaskByIDService(id uint) (model.Task, error) {
	// 目前直接调用仓储层，未来可加入权限校验等
	return repository.FindTaskByID(id)
}

// UpdateTaskService 封装了更新任务的业务逻辑
func UpdateTaskService(id uint, input model.Task) (model.Task, error) {
	// 1. 先根据ID查找出要更新的任务
	task, err := repository.FindTaskByID(id)
	if err != nil {
		return model.Task{}, err // 如果任务不存在，则返回错误
	}

	// 2. 更新字段 (这里我们先简单地更新几个核心字段)
	task.Title = input.Title
	task.Description = input.Description
	task.Priority = input.Priority
	task.Effort = input.Effort

	// 3. 将修改后的完整任务对象存回数据库
	err = repository.UpdateTask(&task)
	if err != nil {
		return model.Task{}, err
	}

	return task, nil
}

// DeleteTaskService 封装了删除任务的业务逻辑
func DeleteTaskService(id uint) error {
	// 1. 先确保任务存在
	_, err := repository.FindTaskByID(id)
	if err != nil {
		return err // 任务不存在
	}

	// 2. 调用仓储层删除任务
	return repository.DeleteTask(id)
}

// ApproveTaskService 封装了审批任务的业务逻辑
func ApproveTaskService(taskID uint, reviewerID uuid.UUID) error {
	// 1. 确保任务存在且状态正确
	task, err := repository.FindTaskByID(taskID)
	if err != nil {
		return errors.New("task not found")
	}
	if task.Status != "pending_review" {
		return errors.New("task is not in pending_review status")
	}

	// 2. 准备要更新的字段
	updates := map[string]interface{}{
		"status":      "in_pool", // 新状态
		"reviewer_id": reviewerID,
		"approved_at": time.Now(),
	}

	// 3. 调用仓储层更新任务
	return repository.UpdateTaskFields(taskID, updates)
}

// ClaimTaskService 封装了领取任务的业务逻辑
func ClaimTaskService(taskID uint, assigneeID uuid.UUID) error {
	// 1. 查找任务并进行业务规则校验
	task, err := repository.FindTaskByID(taskID)
	if err != nil {
		return errors.New("task not found")
	}
	if task.Status != "in_pool" {
		return errors.New("task is not available to be claimed")
	}
	if task.AssigneeID != nil {
		return errors.New("task has already been assigned")
	}

	// 2. 准备要更新的字段
	updates := map[string]interface{}{
		"status":      "in_progress",
		"assignee_id": assigneeID,
		"claimed_at":  time.Now(),
	}

	// 3. 更新数据库
	return repository.UpdateTaskFields(taskID, updates)
}

// CompleteTaskService 封装了完成任务并提交评价的业务逻辑
func CompleteTaskService(taskID uint, currentUserID uuid.UUID) error {
	// 1. 查找任务
	task, err := repository.FindTaskByID(taskID)
	if err != nil {
		return errors.New("task not found")
	}

	// 2. 业务规则校验
	if task.Status != "in_progress" {
		return errors.New("task is not in progress")
	}
	// 权限校验：只有当前负责人才能完成任务
	if task.AssigneeID == nil || *task.AssigneeID != currentUserID {
		return errors.New("permission denied: you are not the assignee of this task")
	}

	// 3. 准备更新
	updates := map[string]interface{}{
		"status": "pending_evaluation",
	}

	// 4. 更新数据库
	return repository.UpdateTaskFields(taskID, updates)
}

// EvaluateTaskService 封装了评价任务的业务逻辑
func EvaluateTaskService(taskID uint, evaluationData datatypes.JSON) error {
	// 1. 查找任务并校验状态
	task, err := repository.FindTaskByID(taskID)
	if err != nil {
		return errors.New("task not found")
	}
	if task.Status != "pending_evaluation" {
		return errors.New("task is not pending evaluation")
	}

	// 2. 准备更新
	updates := map[string]interface{}{
		"status":       "completed",
		"evaluation":   evaluationData,
		"completed_at": time.Now(),
	}

	// 3. 更新数据库
	return repository.UpdateTaskFields(taskID, updates)
}
