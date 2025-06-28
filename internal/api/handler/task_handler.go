// internal/api/handler/task_handler.go

package handler

import (
	"gotasksys/internal/model"
	"gotasksys/internal/service" // 引入service层
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
