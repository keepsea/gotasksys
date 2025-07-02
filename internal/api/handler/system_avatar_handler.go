// internal/api/handler/system_avatar_handler.go
package handler

import (
	"gotasksys/internal/model"
	"gotasksys/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// --- 头像库管理输入结构体 ---

type CreateAvatarInput struct {
	URL         string `json:"url" binding:"required,url"`
	Description string `json:"description"`
}

type UpdateAvatarInput struct {
	URL         string `json:"url" binding:"required,url"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active" binding:"required"`
}

// --- 处理器函数 ---

// ListAvailableAvatars (供普通用户选择)
func ListAvailableAvatars(c *gin.Context) {
	avatars, err := service.ListSystemAvatarsService(false) // false表示只获取启用的
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list available avatars"})
		return
	}
	c.JSON(http.StatusOK, avatars)
}

// ListAllSystemAvatars (供管理员管理)
func ListAllSystemAvatars(c *gin.Context) {
	// 权限已由AdminAuthMiddleware处理
	avatars, err := service.ListSystemAvatarsService(true) // true表示获取全部
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list all system avatars"})
		return
	}
	c.JSON(http.StatusOK, avatars)
}

// CreateSystemAvatar (管理员新增)
func CreateSystemAvatar(c *gin.Context) {
	var input CreateAvatarInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	avatar := model.SystemAvatar{
		URL:         input.URL,
		Description: input.Description,
	}

	createdAvatar, err := service.CreateSystemAvatarService(avatar)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create system avatar"})
		return
	}
	c.JSON(http.StatusCreated, createdAvatar)
}

// UpdateSystemAvatar (管理员修改)
func UpdateSystemAvatar(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid avatar ID"})
		return
	}

	var input UpdateAvatarInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 准备要更新的数据
	avatarUpdateData := model.SystemAvatar{
		URL:         input.URL,
		Description: input.Description,
		IsActive:    *input.IsActive,
	}

	// 【核心修正】调用Service时，同时传入ID和更新数据
	updatedAvatar, err := service.UpdateSystemAvatarService(id, avatarUpdateData)
	if err != nil {
		if err.Error() == "system avatar not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update system avatar"})
		return
	}

	c.JSON(http.StatusOK, updatedAvatar)
}

// DeleteSystemAvatar (管理员删除)
func DeleteSystemAvatar(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid avatar ID"})
		return
	}

	if err := service.DeleteSystemAvatarService(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete system avatar"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "System avatar deleted successfully"})
}
