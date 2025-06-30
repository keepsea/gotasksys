// internal/api/handler/task_handler.go

package handler

import (
	"gotasksys/internal/model"
	"gotasksys/internal/repository"
	"gotasksys/internal/service" // 引入service层
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// CreateTaskInput (最终版): 只包含用户创建时需要提供的信息
type CreateTaskInput struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	DueDate     *time.Time `json:"due_date" binding:"required"`
}

// ApproveTaskInput (最终版): 包含Manager审批时需要设定的所有信息
type ApproveTaskInput struct {
	Effort     int    `json:"effort" binding:"required,gt=0"`
	Priority   string `json:"priority" binding:"required"`
	TaskTypeID string `json:"task_type_id" binding:"required,uuid"`
}

// CreateTask 现在只负责参数解析和调用service
// CreateTask 创建任务时不再处理工时
func CreateTask(c *gin.Context) {
	var input CreateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	creatorID, _ := uuid.Parse(c.GetString("user_id"))

	taskModel := model.Task{
		Title:       input.Title,
		Description: input.Description,
		DueDate:     input.DueDate,
	}

	createdTask, err := service.CreateTaskService(taskModel, creatorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task", "details": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, createdTask)
}

// ListTasks 获取任务列表，现在会根据用户角色返回不同内容
func ListTasks(c *gin.Context) {
	// 从中间件中获取用户信息
	userRole, _ := c.Get("user_role")
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	// 调用重构后的Service，并传入用户信息
	tasks, err := service.ListTasksService(userRole.(string), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks", "details": err.Error()})
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

// UpdateTask 更新一个已存在的任务 (最终权限版)
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

	// 从Context中获取当前用户的ID，并从数据库中查出完整的用户信息
	currentUserID, _ := uuid.Parse(c.GetString("user_id"))
	currentUser, err := repository.FindUserByID(currentUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve current user info"})
		return
	}

	updateData := model.Task{
		Title:       input.Title,
		Description: input.Description,
		Priority:    input.Priority,
		Effort:      input.Effort,
	}

	// 调用Service层，并传入当前用户信息用于权限判断
	updatedTask, err := service.UpdateTaskService(uint(id), currentUser, updateData)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		if err.Error() == "permission denied: you are not authorized to update this task" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
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

// ApproveTask 批准任务时接收并设定工时 (最终版)
func ApproveTask(c *gin.Context) {
	// ... (函数前半部分的代码保持不变) ...
	userRole, _ := c.Get("user_role")
	if userRole != "manager" && userRole != "system_admin" { /* ... */
	}
	taskID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	reviewerID, _ := uuid.Parse(c.GetString("user_id"))

	var input ApproveTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	taskTypeID, _ := uuid.Parse(input.TaskTypeID)

	// --- 这里的组装逻辑有微调 ---
	updateData := model.Task{
		Effort:     input.Effort,
		Priority:   input.Priority,
		TaskTypeID: &taskTypeID, // <-- 修改在这里：因为模型是*uuid.UUID，所以这里要传递指针
	}
	// -------------------------

	err := service.ApproveTaskService(uint(taskID), reviewerID, updateData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task approved successfully."})
}

// ClaimTask 允许用户从任务池中领取任务
func ClaimTask(c *gin.Context) {
	// --- 新增：角色权限校验 ---
	userRole, _ := c.Get("user_role")
	if userRole != "executor" && userRole != "manager" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied. Only executors or managers can claim tasks."})
		return
	}
	// ---------------------------

	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	assigneeIDStr, _ := c.Get("user_id")
	assigneeID, _ := uuid.Parse(assigneeIDStr.(string))

	err = service.ClaimTaskService(uint(taskID), assigneeID)
	// ...后续错误处理逻辑保持不变...
	if err != nil {
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

// CompleteTask 允许负责人完成任务并提交评价 (最终版)
func CompleteTask(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}
	currentUserIDStr, _ := c.Get("user_id")
	currentUserID, _ := uuid.Parse(currentUserIDStr.(string))

	err = service.CompleteTaskService(uint(taskID), currentUserID)
	if err != nil {
		// 根据错误类型返回不同的HTTP状态码
		switch err.Error() {
		case "task not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		// --- 新增：处理主任务完成前置条件失败的错误 ---
		case "cannot complete main task: there are still incomplete subtasks":
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		// ------------------------------------------
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

// EvaluateTask 允许特定权限用户评价一个已完成的任务 (最终版)
func EvaluateTask(c *gin.Context) {
	// --- 不再在这里做简单的角色校验，将由Service层进行复杂的校验 ---

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

	// 从Context中获取当前用户的ID，并从数据库中查出完整的用户信息
	currentUserIDStr, _ := c.Get("user_id")
	currentUserID, _ := uuid.Parse(currentUserIDStr.(string))
	currentUser, err := repository.FindUserByID(currentUserID) // 复用我们已有的repository方法
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve current user info"})
		return
	}

	// 调用Service层处理业务逻辑，传入完整的用户信息
	err = service.EvaluateTaskService(uint(taskID), currentUser, input.Evaluation)
	if err != nil {
		// 根据Service返回的错误类型，给出不同的HTTP响应
		if err.Error() == "permission denied: you are not authorized to evaluate this task" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "task not found" || err.Error() == "parent task not found" {
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

// --- 新增：驳回与重提任务的处理器 ---

// RejectTaskInput 定义了驳回任务时需要输入的参数
type RejectTaskInput struct {
	Reason string `json:"reason" binding:"required"`
}

// RejectTask 驳回一个待审核的任务
func RejectTask(c *gin.Context) {
	// 1. 权限校验：确保只有管理者可以执行此操作
	userRole, _ := c.Get("user_role")
	if userRole != "manager" && userRole != "system_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied. Only managers can reject tasks."})
		return
	}

	// 2. 解析URL中的任务ID
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// 3. 解析请求体中的驳回理由
	var input RejectTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Reason is required"})
		return
	}

	// 4. 获取当前操作者(即审批人)的ID
	reviewerIDStr, _ := c.Get("user_id")
	reviewerID, _ := uuid.Parse(reviewerIDStr.(string))

	// 5. 调用Service层处理业务逻辑
	err = service.RejectTaskService(uint(taskID), input.Reason, reviewerID)
	if err != nil {
		// 根据Service返回的错误类型，给出不同的HTTP响应
		if err.Error() == "task not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "task is not in pending_review status" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reject task"})
		return
	}

	// 6. 返回成功响应
	c.JSON(http.StatusOK, gin.H{"message": "Task rejected."})
}

// ResubmitTask 重新提交一个被驳回的任务
func ResubmitTask(c *gin.Context) {
	// 1. 解析URL中的任务ID
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// 2. 获取当前操作者(即任务创建人)的ID
	creatorIDStr, _ := c.Get("user_id")
	creatorID, _ := uuid.Parse(creatorIDStr.(string))

	// 3. 调用Service层处理业务逻辑 (Service层内部会进行权限校验)
	err = service.ResubmitTaskService(uint(taskID), creatorID)
	if err != nil {
		// 根据Service返回的错误类型，给出不同的HTTP响应
		if err.Error() == "task not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "permission denied: only the creator can resubmit the task" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "task is not in rejected status" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resubmit task"})
		return
	}

	// 4. 返回成功响应
	c.JSON(http.StatusOK, gin.H{"message": "Task resubmitted successfully."})
}

// CreateSubtaskInput (最终锁定版)
// - 移除了 Priority (将自动继承)
// - 新增了 DueDate (必需)
type CreateSubtaskInput struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	Effort      int        `json:"effort" binding:"required,gte=1"`
	DueDate     *time.Time `json:"due_date" binding:"required"`
}

// -----------------------------------------

// CreateSubtask (最终锁定版)
func CreateSubtask(c *gin.Context) {
	// 1. 解析父任务ID (逻辑不变)
	parentTaskIDStr := c.Param("id")
	parentTaskID, err := strconv.ParseUint(parentTaskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parent task ID"})
		return
	}

	// 2. 使用新的输入结构体来解析请求体
	var input CreateSubtaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. 获取当前用户ID (逻辑不变)
	creatorID, _ := uuid.Parse(c.GetString("user_id"))

	// 4. 组装将要传递给Service层的数据模型
	// 注意：这里不再包含Priority
	subtaskInput := model.Task{
		Title:       input.Title,
		Description: input.Description,
		Effort:      input.Effort,
		DueDate:     input.DueDate,
	}

	// 5. 调用Service层处理核心业务逻辑
	createdSubtask, err := service.CreateSubtaskService(uint(parentTaskID), creatorID, subtaskInput)
	if err != nil {
		// 6. 完整的错误处理，包含我们新增的所有校验
		switch err.Error() {
		case "subtask due date cannot be after the parent task's due date":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case "total effort of subtasks cannot exceed parent task's original effort":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case "permission denied: only the assignee of the main task can create subtasks":
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case "parent task not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "only in-progress tasks can have subtasks":
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subtask"})
		}
		return
	}

	// 7. 返回成功响应
	c.JSON(http.StatusCreated, createdSubtask)
}
