// internal/api/handler/profile_handler.go
package handler

import (
	"gotasksys/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateProfileInput struct {
	RealName string `json:"real_name"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
	Team     string `json:"team"`
}

// UpdateMyProfile 处理用户更新自己的姓名、头像、邮箱、团队
func UpdateMyProfile(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))
	var input UpdateProfileInput
	// 使用Bind而不是ShouldBindJSON，如果json没有这个字段，它会是零值而不是错误
	if err := c.Bind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.UpdateMyProfileService(userID, input.RealName, input.Avatar, input.Email, input.Team); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

type ChangePasswordInput struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// ChangeMyPassword 处理用户修改自己的密码
func ChangeMyPassword(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))
	var input ChangePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := service.ChangeMyPasswordService(userID, input.OldPassword, input.NewPassword); err != nil {
		// 根据错误类型返回不同响应
		if err.Error() == "old password is incorrect" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}
