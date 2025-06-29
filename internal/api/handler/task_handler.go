// internal/api/handler/task_handler.go

package handler

import (
	"gotasksys/internal/model"
	"gotasksys/internal/service" // 引入service层
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// CreateTaskInput 结构体保持不变
type CreateTaskInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Priority    string `json:"priority" binding:"required"`
	Effort      int    `json:"effort"`
	TaskTypeID  string `json:"task_type_id" binding:"required,uuid"`
}

// CreateTask 现在只负责参数解析和调用service
func CreateTask(c *gin.Context) {
	var input CreateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	creatorIDStr, _ := c.Get("user_id")
	creatorID, _ := uuid.Parse(creatorIDStr.(string))

	taskTypeID, _ := uuid.Parse(input.TaskTypeID)

	// 将输入参数组装成模型对象
	taskModel := model.Task{
		Title:       input.Title,
		Description: input.Description,
		Priority:    input.Priority,
		Effort:      input.Effort,
		TaskTypeID:  taskTypeID,
		CreatorID:   creatorID,
	}

	// 将业务逻辑委托给service层
	createdTask, err := service.CreateTaskService(taskModel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdTask)
}

// ListTasks 现在只负责调用service并返回结果
func ListTasks(c *gin.Context) {
	tasks, err := service.ListTasksService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// GetTask 获取单个任务的详情
func GetTask(c *gin.Context) {
	// 从URL路径中获取id参数, e.g., /tasks/123
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// 将业务委托给service层
	task, err := service.GetTaskByIDService(uint(id))
	if err != nil {
		// gorm.ErrRecordNotFound 是一个常见的错误，我们应该返回404
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// UpdateTaskInput 定义了更新任务时允许输入的参数
// 注意：这里的字段应该和Service层中允许更新的字段对应
type UpdateTaskInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Priority    string `json:"priority" binding:"required"`
	Effort      int    `json:"effort"`
}

// UpdateTask 更新一个已存在的任务
func UpdateTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var input UpdateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateData := model.Task{
		Title:       input.Title,
		Description: input.Description,
		Priority:    input.Priority,
		Effort:      input.Effort,
	}

	updatedTask, err := service.UpdateTaskService(uint(id), updateData)
	if err != nil {
		if err.Error() == "record not found" { // gorm.ErrRecordNotFound 的字符串形式
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	c.JSON(http.StatusOK, updatedTask)
}

// DeleteTask 删除一个任务
func DeleteTask(c *gin.Context) {
	// === 权限校验 ===
	userRole, _ := c.Get("user_role")
	if userRole != "manager" && userRole != "system_admin" { // 允许管理员和系统管理员删除
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied. Only managers can delete tasks."})
		return
	}
	// =================

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	err = service.DeleteTaskService(uint(id))
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

// ApproveTask 批准一个待审核的任务
func ApproveTask(c *gin.Context) {
	// 权限校验
	userRole, _ := c.Get("user_role")
	if userRole != "manager" && userRole != "system_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied. Only managers can approve tasks."})
		return
	}

	// 获取任务ID
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// 获取审批人ID
	reviewerIDStr, _ := c.Get("user_id")
	reviewerID, _ := uuid.Parse(reviewerIDStr.(string))

	// 调用Service层处理业务逻辑
	err = service.ApproveTaskService(uint(taskID), reviewerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task approved successfully and moved to pool."})
}

// ClaimTask 允许用户从任务池中领取任务
func ClaimTask(c *gin.Context) {
	// 任何登录的用户都可以领取任务，所以我们只检查登录状态，不校验角色

	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// 从中间件获取当前操作用户的ID
	assigneeIDStr, _ := c.Get("user_id")
	assigneeID, _ := uuid.Parse(assigneeIDStr.(string))

	// 调用Service层处理业务逻辑
	err = service.ClaimTaskService(uint(taskID), assigneeID)
	if err != nil {
		// 如果任务状态不对或已被分配，返回409 Conflict更合适
		if err.Error() == "task is not available to be claimed" || err.Error() == "task has already been assigned" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "task not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to claim task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task claimed successfully."})
}

// CompleteTask 允许负责人完成任务并提交评价
func CompleteTask(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// 从中间件获取当前操作用户的ID
	currentUserIDStr, _ := c.Get("user_id")
	currentUserID, _ := uuid.Parse(currentUserIDStr.(string))

	// 调用Service层处理业务逻辑
	err = service.CompleteTaskService(uint(taskID), currentUserID)
	if err != nil {
		// 根据错误类型返回不同的HTTP状态码
		switch err.Error() {
		case "task not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "permission denied: you are not the assignee of this task":
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case "task is not in progress":
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete task"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task completed and submitted for evaluation."})
}

// EvaluateTaskInput 定义了评价任务时需要输入的参数
type EvaluateTaskInput struct {
	// 我们直接使用 datatypes.JSON 来接收任意结构的JSON评价数据
	// 前端可以传入 {"timeliness": 5, "quality": 4.5, ...} 这样的格式
	Evaluation datatypes.JSON `json:"evaluation" binding:"required"`
}

// EvaluateTask 允许管理者评价一个已完成的任务
func EvaluateTask(c *gin.Context) {
	// 权限校验
	userRole, _ := c.Get("user_role")
	if userRole != "manager" && userRole != "system_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied. Only managers can evaluate tasks."})
		return
	}

	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var input EvaluateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用Service层处理业务逻辑
	err = service.EvaluateTaskService(uint(taskID), input.Evaluation)
	if err != nil {
		if err.Error() == "task not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "task is not pending evaluation" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to evaluate task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task evaluated successfully and marked as completed."})
}
