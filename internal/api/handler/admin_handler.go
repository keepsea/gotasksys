// internal/api/handler/admin_handler.go
package handler

import (
	"gotasksys/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// --- 任务类型管理 ---

func ListTaskTypes(c *gin.Context) {
	types, err := service.ListTaskTypesService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list task types"})
		return
	}
	c.JSON(http.StatusOK, types)
}

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

// --- 用户管理 ---
func ListUsers(c *gin.Context) {
	users, err := service.ListUsersService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

type UpdateUserRoleInput struct {
	Role string `json:"role" binding:"required"`
}

func UpdateUserRole(c *gin.Context) {
	userID, _ := uuid.Parse(c.Param("id"))
	var input UpdateUserRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 此处也应增加对role值的校验
	allowedRoles := map[string]bool{"system_admin": true, "manager": true, "executor": true, "creator": true}
	if !allowedRoles[input.Role] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role specified"})
		return
	}

	if err := service.UpdateUserRoleService(userID, input.Role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user role"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User role updated successfully"})
}

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

func DeleteUser(c *gin.Context) {
	userID, _ := uuid.Parse(c.Param("id"))
	if err := service.DeleteUserService(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

type UpdateTaskTypeInput struct {
	Name      string `json:"name" binding:"required"`
	IsEnabled bool   `json:"is_enabled"`
}

// UpdateTaskType 修改任务类型
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

// DeleteTaskType 删除任务类型
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
