// internal/api/handler/admin_handler.go
package handler

import (
	"gotasksys/internal/model"
	"gotasksys/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// --- 用户管理 ---
// 1、查看用户列表
func ListUsers(c *gin.Context) {
	users, err := service.ListUsersService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// 2、创建用户
type CreateUserInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	RealName string `json:"real_name" binding:"required"`
	Role     string `json:"role" binding:"required"`
	Email    string `json:"email" binding:"omitempty,email"`
	Team     string `json:"team"`
}

func CreateUser(c *gin.Context) {
	var input CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 角色有效性验证
	allowedRoles := map[string]bool{"system_admin": true, "manager": true, "executor": true, "creator": true}
	if !allowedRoles[input.Role] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role specified: " + input.Role})
		return
	}
	// 密码复杂度验证
	if len(input.Password) < 12 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 12 characters long"})
		return
	}

	user := model.User{
		Username:     input.Username,
		PasswordHash: input.Password, // Service层会进行哈希
		RealName:     input.RealName,
		Role:         input.Role,
		Email:        input.Email,
		Team:         input.Team,
	}

	userID, err := service.RegisterUser(user) // 复用RegisterUser服务
	if err != nil {
		if err.Error() == "username already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user_id": userID})
}

// 3、更新用户角色
type UpdateUserRoleInput struct {
	Role string `json:"role" binding:"required"`
}

// UpdateUserInput (最终版)
type UpdateUserInput struct {
	RealName           string   `json:"real_name"`
	Role               string   `json:"role"`
	Team               string   `json:"team"`
	Email              string   `json:"email" binding:"omitempty,email"` // <-- 新增
	DailyCapacityHours *float64 `json:"daily_capacity_hours"`
}

// UpdateUser (新) - 统一的用户更新接口
func UpdateUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var input UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateData := model.User{
		RealName:           input.RealName,
		Role:               input.Role,
		Team:               input.Team,
		Email:              input.Email,
		DailyCapacityHours: input.DailyCapacityHours,
	}

	if err := service.UpdateUserByAdminService(userID, updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// 4、重置用户密码
type ResetPasswordInput struct {
	NewPassword string `json:"new_password" binding:"required"`
}

func ResetPassword(c *gin.Context) {
	userID, _ := uuid.Parse(c.Param("id"))
	var input ResetPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := service.ResetPasswordService(userID, input.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User password reset successfully"})
}

// 5、删除用户
func DeleteUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = service.DeleteUserService(userID)
	if err != nil {
		// --- 新增：根据错误类型返回不同响应 ---
		if err.Error() == "cannot delete user: user has unfinished tasks. Please transfer or complete them first" {
			// 对于业务逻辑冲突，返回 409 Conflict
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		// 对于其他错误，返回通用服务器错误
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// --- 任务类型管理 ---
// 1、查看任务类型
func ListTaskTypes(c *gin.Context) {
	types, err := service.ListTaskTypesService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list task types"})
		return
	}
	c.JSON(http.StatusOK, types)
}

// 2、创建任务类型
type CreateTaskTypeInput struct {
	Name string `json:"name" binding:"required"`
}

func CreateTaskType(c *gin.Context) {
	var input CreateTaskTypeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdType, err := service.CreateTaskTypeService(input.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task type"})
		return
	}
	c.JSON(http.StatusCreated, createdType)
}

// 3、更新任务类型
type UpdateTaskTypeInput struct {
	Name      string `json:"name" binding:"required"`
	IsEnabled bool   `json:"is_enabled"`
}

// 4、修改任务类型
func UpdateTaskType(c *gin.Context) {
	typeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task type ID"})
		return
	}
	var input UpdateTaskTypeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := service.UpdateTaskTypeService(typeID, input.Name, input.IsEnabled); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task type"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task type updated successfully"})
}

// 5、删除任务类型
func DeleteTaskType(c *gin.Context) {
	typeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task type ID"})
		return
	}
	if err := service.DeleteTaskTypeService(typeID); err != nil {
		// 如果是“正在使用中”的业务错误，返回409 Conflict更合适
		if err.Error() == "cannot delete task type: it is currently in use by one or more tasks" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task type"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task type deleted successfully"})
}
