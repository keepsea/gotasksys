// internal/api/middleware/manager_auth.go
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ManagerAuthMiddleware 检查用户角色是否为 "manager" 或 "system_admin"
func ManagerAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")

		// 如果上下文中没有角色信息，或者角色不是manager或system_admin，则拒绝访问
		if !exists || (userRole.(string) != "manager" && userRole.(string) != "system_admin") {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Permission denied. Manager role required."})
			return
		}

		// 权限校验通过，继续处理请求
		c.Next()
	}
}
