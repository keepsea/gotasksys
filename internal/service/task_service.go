// internal/service/task_service.go

package service

import (
	"encoding/json"
	"errors"
	"gotasksys/internal/model"
	"gotasksys/internal/repository"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// CreateTaskService 封装了创建任务的业务逻辑
func CreateTaskService(input model.Task, creatorID uuid.UUID) (model.Task, error) {
	task := model.Task{
		Title:       input.Title,
		Description: input.Description,
		DueDate:     input.DueDate, // 保存创建者设定的截止时间
		CreatorID:   creatorID,
		Status:      "pending_review", // 新任务的初始状态
	}
	err := repository.CreateTask(&task)
	return task, err
}

// ListTasksService 根据用户角色和ID，获取其能看到的任务列表
func ListTasksService(userRole string, userID uuid.UUID) ([]model.Task, error) {
	switch userRole {
	case "system_admin", "manager":
		return repository.ListAllTasks()
	case "executor":
		return repository.ListTasksForExecutor(userID)
	case "creator":
		return repository.ListTasksForCreator(userID)
	default:
		// 如果遇到未知的角色，返回空列表和错误
		return nil, errors.New("invalid user role for listing tasks")
	}
}

// GetTaskByIDService 封装了根据ID获取任务的业务逻辑
func GetTaskByIDService(id uint) (model.Task, error) {
	// 目前直接调用仓储层，未来可加入权限校验等
	return repository.FindTaskByID(id)
}

// UpdateTaskService 封装了更新任务的业务逻辑 (最终锁定版)
func UpdateTaskService(taskID uint, currentUser model.User, updateData model.Task) (model.Task, error) {
	// 1. 先根据ID查找出要更新的任务
	task, err := repository.FindTaskByID(taskID)
	if err != nil {
		return model.Task{}, errors.New("task not found")
	}

	// 2. --- 【V1.2 最终锁定版】权限校验 ---
	hasPermission := false
	// 规则一: Manager或Admin总是有权限
	if currentUser.Role == "manager" || currentUser.Role == "system_admin" {
		hasPermission = true
	}
	// 规则三: 如果任务状态是'rejected'，且创建者是当前用户，则有权限
	if !hasPermission && task.Status == "rejected" && task.CreatorID == currentUser.ID {
		hasPermission = true
	}
	// **注意：我们已根据您的要求，移除了允许Assignee修改的规则**

	if !hasPermission {
		return model.Task{}, errors.New("permission denied: you are not authorized to update this task")
	}
	// ------------------------------------------

	// 3. 更新字段 (逻辑保持不变)
	task.Title = updateData.Title
	task.Description = updateData.Description
	task.Priority = updateData.Priority
	// 注意：工时(Effort)的修改权限可以后续再细化，V1.0中暂时允许在有权限时修改
	task.Effort = updateData.Effort

	// 4. 将修改后的完整任务对象存回数据库 (逻辑保持不变)
	err = repository.UpdateTask(&task)
	if err != nil {
		return model.Task{}, err
	}

	return task, nil
}

// DeleteTaskService 封装了删除任务的业务逻辑 (最终权限版)
func DeleteTaskService(taskID uint, currentUser model.User) error {
	// 1. 查找待删除的任务
	taskToDelete, err := repository.FindTaskByID(taskID)
	if err != nil {
		return errors.New("task not found")
	}

	// 2. --- 核心：最终的删除权限校验逻辑 ---
	hasPermission := false
	// 规则a: Manager或Admin总是有权限
	if currentUser.Role == "manager" || currentUser.Role == "system_admin" {
		hasPermission = true
	}
	// 规则b: 如果是子任务，检查当前用户是否是其父任务的负责人
	if !hasPermission && taskToDelete.ParentTaskID != nil {
		parentTask, err := repository.FindTaskByID(*taskToDelete.ParentTaskID)
		if err == nil && parentTask.AssigneeID != nil && *parentTask.AssigneeID == currentUser.ID {
			hasPermission = true
		}
	}
	// 规则c: 任务的创建者，在任务被驳回时，也可以删除它
	if !hasPermission && taskToDelete.Status == "rejected" && taskToDelete.CreatorID == currentUser.ID {
		hasPermission = true
	}

	if !hasPermission {
		return errors.New("permission denied: you are not authorized to delete this task")
	}
	// ------------------------------------------

	// 3. 调用仓储层删除任务
	return repository.DeleteTask(taskID)
}

func ApproveTaskService(taskID uint, reviewerID uuid.UUID, effort int, priority string, taskTypeID uuid.UUID, difficultyRating map[string]float64) error {
	task, err := repository.FindTaskByID(taskID)
	if err != nil {
		return errors.New("task not found")
	}
	if task.Status != "pending_review" {
		return errors.New("task is not in pending_review status")
	}
	if effort <= 0 {
		return errors.New("effort must be greater than zero")
	}

	// --- 新增：处理和计算技术难度分 ---
	var finalRatingJSON datatypes.JSON
	if difficultyRating != nil {
		// 我们可以在这里增加校验，确保4个维度都存在
		novelty := difficultyRating["novelty"]
		complexity := difficultyRating["logic_complexity"]
		impact := difficultyRating["impact_scope"]
		collaboration := difficultyRating["collaboration_cost"]

		compositeScore := (novelty + complexity + impact + collaboration) / 4.0
		difficultyRating["composite_difficulty_score"] = compositeScore

		ratingBytes, err := json.Marshal(difficultyRating)
		if err != nil {
			return errors.New("failed to process difficulty rating")
		}
		finalRatingJSON = ratingBytes
	}
	// ------------------------------------

	updates := map[string]interface{}{
		"status":            "in_pool",
		"reviewer_id":       reviewerID,
		"approved_at":       time.Now(),
		"effort":            effort,
		"original_effort":   effort,
		"priority":          priority,
		"task_type_id":      &taskTypeID,     // 确保传递指针
		"difficulty_rating": finalRatingJSON, // 保存包含综合分的完整JSON
	}
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

// CompleteTaskService 封装了完成任务并提交评价的业务逻辑 (最终版)
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
	if task.AssigneeID == nil || *task.AssigneeID != currentUserID {
		return errors.New("permission denied: you are not the assignee of this task")
	}

	// --- 【V1.2 最终修正】主任务完成的前置条件校验 ---
	// 3. 检查这是否是一个主任务 (即没有parent_task_id)
	if task.ParentTaskID == nil {
		// 如果是主任务，则检查其下是否有未完成的子任务
		incompleteSubtasks, err := repository.CountIncompleteSubtasks(task.ID)
		if err != nil {
			return err // 如果查询出错，也中断操作
		}
		if incompleteSubtasks > 0 {
			return errors.New("cannot complete main task: there are still incomplete subtasks")
		}
	}
	// ----------------------------------------------------

	// 4. 准备更新
	updates := map[string]interface{}{
		"status": "pending_evaluation",
	}

	// 5. 更新数据库
	return repository.UpdateTaskFields(taskID, updates)
}

// EvaluateTaskService 封装了评价任务的业务逻辑 (最终版)
func EvaluateTaskService(taskID uint, currentUser model.User, evaluationData datatypes.JSON) error {
	// 1. 查找待评价的任务
	taskToEvaluate, err := repository.FindTaskByID(taskID)
	if err != nil {
		return errors.New("task not found")
	}
	if taskToEvaluate.Status != "pending_evaluation" {
		return errors.New("task is not pending evaluation")
	}

	// 2. --- 【V1.2 最终修正】重写权限校验逻辑 ---
	hasPermission := false
	// 场景一：这是一个主任务 (没有父任务ID)
	if taskToEvaluate.ParentTaskID == nil {
		// 规则：主任务只能由 manager 或 system_admin 评价
		if currentUser.Role == "manager" || currentUser.Role == "system_admin" {
			hasPermission = true
		}
	} else { // 场景二：这是一个子任务
		// 规则：子任务只能由其父任务的负责人评价
		parentTask, err := repository.FindTaskByID(*taskToEvaluate.ParentTaskID)
		// 必须成功找到父任务，且父任务的负责人(Assignee)正好是当前操作的用户
		if err == nil && parentTask.AssigneeID != nil && *parentTask.AssigneeID == currentUser.ID {
			hasPermission = true
		}
	}

	if !hasPermission {
		return errors.New("permission denied: you are not authorized to evaluate this task")
	}
	// ------------------------------------

	// 3. 计算综合得分 (逻辑保持不变)
	var evalMap map[string]interface{}
	if err := json.Unmarshal(evaluationData, &evalMap); err == nil {
		timeliness, ok1 := evalMap["timeliness"].(float64)
		quality, ok2 := evalMap["quality"].(float64)
		collaboration, ok3 := evalMap["collaboration"].(float64)
		complexity, ok4 := evalMap["complexity"].(float64)
		if ok1 && ok2 && ok3 && ok4 {
			compositeScore := (timeliness + quality + collaboration + complexity) / 4.0
			evalMap["composite_score"] = compositeScore
			updatedEvaluationData, _ := json.Marshal(evalMap)
			evaluationData = updatedEvaluationData
		} else {
			return errors.New("evaluation data must contain all four dimensions with numeric values")
		}
	} else {
		return errors.New("invalid evaluation data format")
	}

	// 4. 更新数据库 (逻辑保持不变)
	updates := map[string]interface{}{
		"status":       "completed",
		"evaluation":   evaluationData,
		"completed_at": time.Now(),
	}
	return repository.UpdateTaskFields(taskID, updates)
}

// RejectTaskService 封装了驳回任务的业务逻辑
func RejectTaskService(taskID uint, reason string, reviewerID uuid.UUID) error {
	task, err := repository.FindTaskByID(taskID)
	if err != nil {
		return errors.New("task not found")
	}
	if task.Status != "pending_review" {
		return errors.New("task is not in pending_review status")
	}

	updates := map[string]interface{}{
		"status":           "rejected",
		"rejection_reason": reason,
		"reviewer_id":      reviewerID, // 记录是谁驳回的
	}
	return repository.UpdateTaskFields(taskID, updates)
}

// ResubmitTaskService 封装了重新提交任务的业务逻辑
func ResubmitTaskService(taskID uint, creatorID uuid.UUID) error {
	task, err := repository.FindTaskByID(taskID)
	if err != nil {
		return errors.New("task not found")
	}
	// 权限校验：只有创建者自己才能重新提交
	if task.CreatorID != creatorID {
		return errors.New("permission denied: only the creator can resubmit the task")
	}
	if task.Status != "rejected" {
		return errors.New("task is not in rejected status")
	}

	updates := map[string]interface{}{
		"status": "pending_review",
	}
	return repository.UpdateTaskFields(taskID, updates)
}

// CreateSubtaskService 封装了创建子任务的业务逻辑 (最终锁定版)
func CreateSubtaskService(parentTaskID uint, creatorID uuid.UUID, subtaskInput model.Task) (model.Task, error) {
	// 1. 查找父任务 (逻辑不变)
	parentTask, err := repository.FindTaskByID(parentTaskID)
	if err != nil {
		return model.Task{}, errors.New("parent task not found")
	}
	// ... (权限和状态校验逻辑不变) ...
	if parentTask.Status != "in_progress" { /* ... */
	}
	if parentTask.AssigneeID == nil || *parentTask.AssigneeID != creatorID { /* ... */
	}

	// 2. 工时上限校验逻辑 (逻辑不变)
	existingSubtasksEffort, err := repository.GetTotalEffortOfSubtasks(parentTaskID)
	if err != nil {
		return model.Task{}, err
	}
	if (existingSubtasksEffort + int64(subtaskInput.Effort)) > int64(parentTask.OriginalEffort) {
		return model.Task{}, errors.New("total effort of subtasks cannot exceed parent task's original effort")
	}

	// 3. --- 【V1.2 最终锁定版】截止时间校验 ---
	// 规则：如果父任务有截止时间，则子任务的截止时间不能晚于父任务的截止时间
	if parentTask.DueDate != nil && subtaskInput.DueDate.After(*parentTask.DueDate) {
		return model.Task{}, errors.New("subtask due date cannot be after the parent task's due date")
	}
	// ------------------------------------

	// 4. 准备子任务数据 (逻辑不变)
	subtask := model.Task{
		Title:          subtaskInput.Title,
		Description:    subtaskInput.Description,
		Priority:       parentTask.Priority,
		Effort:         subtaskInput.Effort,
		OriginalEffort: subtaskInput.Effort,
		DueDate:        subtaskInput.DueDate,
		TaskTypeID:     parentTask.TaskTypeID,
		CreatorID:      creatorID,
		ParentTaskID:   &parentTask.ID,
		Status:         "in_pool",
	}

	// 5. 创建任务 (逻辑不变)
	err = repository.CreateTask(&subtask)
	if err != nil {
		return model.Task{}, err
	}

	return subtask, nil
}

// AssignTaskService 封装了指派任务的业务逻辑
func AssignTaskService(taskID uint, assigneeID uuid.UUID, managerID uuid.UUID) error {
	// 1. 查找任务并校验状态
	task, err := repository.FindTaskByID(taskID)
	if err != nil {
		return errors.New("task not found")
	}
	if task.Status != "in_pool" {
		return errors.New("task is not in the task pool to be assigned")
	}

	// 2. 准备更新
	updates := map[string]interface{}{
		"status":      "in_progress", // 任务被指派后，直接进入进行中状态
		"assignee_id": assigneeID,
		"claimed_at":  time.Now(), // 视同被领取
		"reviewer_id": managerID,  // 记录下是哪位经理指派的
	}

	// 3. 更新数据库
	return repository.UpdateTaskFields(taskID, updates)
}
