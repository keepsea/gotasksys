// internal/api/handler/user_handler.go

package handler

import (
	"gotasksys/internal/config"
	"gotasksys/internal/repository"
	"gotasksys/internal/service" // 引入service层
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// --- 用于请求参数绑定的结构体 ---

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// --- API处理器函数 ---

// Login 处理用户登录请求
func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load system config"})
		return
	}

	// --- 关键修改在这里 ---
	// 现在我们从service层接收 user, token, err 这三个返回值
	user, token, err := service.LoginUser(input.Username, input.Password, cfg.JWT.Secret, cfg.JWT.ExpirationHours)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()}) // "invalid credentials" 或 "failed to generate token"
		return
	}
	// --------------------

	// 现在，'user'变量在这里是已声明且有值的，可以安全使用
	c.JSON(http.StatusOK, gin.H{
		"token":   token,
		"user_id": user.ID,
		"role":    user.Role,
	})
}

// GetProfile 获取当前登录用户的信息
func GetProfile(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format in context"})
		return
	}

	// 对于简单的查询，Handler可以直接调用Repository层
	user, err := repository.FindUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 返回已脱敏的用户信息
	c.JSON(http.StatusOK, gin.H{
		"id":        user.ID,
		"username":  user.Username,
		"real_name": user.RealName,
		"role":      user.Role,
	})
}
