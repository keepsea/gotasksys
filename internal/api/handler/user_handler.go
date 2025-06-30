// internal/api/handler/user_handler.go

package handler

import (
	"gotasksys/internal/config"
	"gotasksys/internal/model"
	"gotasksys/internal/repository"
	"gotasksys/internal/service" // 引入service层
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// --- 用于请求参数绑定的结构体 ---

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	RealName string `json:"real_name" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// --- API处理器函数 ---

// Register 由管理员创建一个新用户
func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 角色有效性验证 (纠偏计划任务1)
	allowedRoles := map[string]bool{
		"system_admin": true,
		"manager":      true,
		"executor":     true,
		"creator":      true,
	}
	if !allowedRoles[input.Role] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role specified: " + input.Role})
		return
	}

	// 将输入参数组装成模型对象
	// 注意：我们将原始密码传递给Service层，由Service层负责调用哈希工具
	user := model.User{
		Username:     input.Username,
		PasswordHash: input.Password, // 临时存储原始密码
		RealName:     input.RealName,
		Role:         input.Role,
	}

	// 将业务逻辑委托给Service层
	userID, err := service.RegisterUser(user)
	if err != nil {
		// 根据Service层返回的错误类型，给出不同的HTTP响应
		if err.Error() == "username already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 其他业务错误，如密码太短
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user_id": userID})
}

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
