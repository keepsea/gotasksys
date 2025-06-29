/*// internal/api/middleware/admin_auth.go
package middleware

import (
	"net/http"




	"github.com/gin-gonic/gin"
)

func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists || userRole.(string) != "system_admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Permission denied. System administrator role required."})
			return
		}
		c.Next()
	}
}
*/
// internal/api/middleware/admin_auth.go (带有日志的调试版本)
package middleware

import (
	"github.com/gin-gonic/gin"
	"log" // <-- 新增导入 "log" 包
	"net/http"
)

func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// --- 添加的调试日志 ---
		log.Println("--- [AdminAuthMiddleware] IS RUNNING! ---")
		// ---------------------

		userRole, exists := c.Get("user_role")

		// --- 添加的调试日志 ---
		log.Printf("[AdminAuthMiddleware] Checking role. Role found in context: %v, Role value: '%v'", exists, userRole)
		// ---------------------

		if !exists || userRole.(string) != "system_admin" {
			// --- 添加的调试日志 ---
			log.Printf("[AdminAuthMiddleware] Permission DENIED for role: '%v'", userRole)
			// ---------------------
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Permission denied. System administrator role required."})
			return
		}

		// --- 添加的调试日志 ---
		log.Printf("[AdminAuthMiddleware] Permission GRANTED for role: '%v'", userRole)
		// ---------------------
		c.Next()
	}
}
