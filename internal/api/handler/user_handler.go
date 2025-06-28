// internal/api/handler/user_handler.go
package handler

import (
	"gotasksys/internal/config"
	"gotasksys/internal/model"
	"gotasksys/internal/service" // 引入service层
	"net/http"

	"github.com/gin-gonic/gin"
)

// ... RegisterInput 和 LoginInput 结构体定义保持不变 ...

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

func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := model.User{
		Username:     input.Username,
		PasswordHash: input.Password, // 临时将原始密码存在这里
		RealName:     input.RealName,
		Role:         input.Role,
	}

	userID, err := service.RegisterUser(user)
	if err != nil {
		if err.Error() == "username already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user_id": userID})
}

func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	cfg, _ := config.LoadConfig("config.yaml")
	token, err := service.LoginUser(input.Username, input.Password, cfg.JWT.Secret, cfg.JWT.ExpirationHours)

	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// ... GetProfile 函数可以暂时保持不变 ...

// GetProfile 获取当前登录用户的信息
func GetProfile(c *gin.Context) {
	// 从中间件设置的Context中获取用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
		return
	}

	var user model.User
	// 使用从token中解析出的ID来查询数据库
	if result := config.DB.First(&user, "id = ?", userID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 返回用户信息（注意不要返回密码哈希）
	c.JSON(http.StatusOK, gin.H{
		"id":        user.ID,
		"username":  user.Username,
		"real_name": user.RealName,
		"role":      user.Role,
	})
}
