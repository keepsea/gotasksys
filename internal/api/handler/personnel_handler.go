// internal/api/handler/personnel_handler.go
package handler

import (
	"gotasksys/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPersonnelStatus(c *gin.Context) {
	userRole, _ := c.Get("user_role")
	if userRole != "manager" && userRole != "system_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	statuses, err := service.GetPersonnelStatusService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get personnel status"})
		return
	}

	c.JSON(http.StatusOK, statuses)
}
