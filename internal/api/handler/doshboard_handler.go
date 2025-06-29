// internal/api/handler/dashboard_handler.go
package handler

import (
	"gotasksys/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetDashboardSummary 获取驾驶舱统计信息
func GetDashboardSummary(c *gin.Context) {
	userRole, _ := c.Get("user_role")
	if userRole != "manager" && userRole != "system_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	summary, err := service.GetDashboardSummaryService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get dashboard summary"})
		return
	}

	c.JSON(http.StatusOK, summary)
}
