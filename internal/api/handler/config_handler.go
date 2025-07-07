// internal/api/handler/config_handler.go
package handler

import (
	"gotasksys/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// --- 系统配置相关 Handler ---

type UpdateConfigInput struct {
	Value string `json:"value" binding:"required"`
}

// UpdateSystemConfig 更新一个系统配置项，例如'global_daily_work_hours'
func UpdateSystemConfig(c *gin.Context) {
	// 从URL路径中获取要更新的配置的键(key)
	key := c.Param("key")
	var input UpdateConfigInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := service.UpdateSystemConfigService(key, input.Value); err != nil {
		// service层会进行校验，如果值不合法，会返回错误
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "System config updated successfully"})
}

// --- 法定节假日相关 Handler ---

type CreateHolidayInput struct {
	Date        string `json:"date" binding:"required"` // 格式 YYYY-MM-DD
	Description string `json:"description"`
}

// ListHolidays 获取所有节假日
func ListHolidays(c *gin.Context) {
	holidays, err := service.ListHolidaysService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list holidays"})
		return
	}
	c.JSON(http.StatusOK, holidays)
}

// CreateHoliday 添加一个新的节假日
func CreateHoliday(c *gin.Context) {
	var input CreateHolidayInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	date, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD."})
		return
	}
	holiday, err := service.CreateHolidayService(date, input.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create holiday"})
		return
	}
	c.JSON(http.StatusCreated, holiday)
}

// DeleteHoliday 删除一个节假日
func DeleteHoliday(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid holiday ID"})
		return
	}
	if err := service.DeleteHolidayService(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete holiday"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Holiday deleted successfully"})
}
