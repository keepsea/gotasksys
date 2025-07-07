// internal/api/handler/leave_handler.go
package handler

import (
	"gotasksys/internal/model"
	"gotasksys/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LeaveInput struct {
	StartDate string `json:"start_date" binding:"required"` // 使用YYYY-MM-DD格式
	EndDate   string `json:"end_date" binding:"required"`
	Reason    string `json:"reason"`
}

// CreateLeaveHandler 处理新的请假申请
func CreateLeaveHandler(c *gin.Context) {
	var input LeaveInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 解析日期字符串
	startDate, err1 := time.Parse("2006-01-02", input.StartDate)
	endDate, err2 := time.Parse("2006-01-02", input.EndDate)
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Please use YYYY-MM-DD."})
		return
	}

	userID, _ := uuid.Parse(c.GetString("user_id"))

	leave := model.Leave{
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
		Reason:    input.Reason,
	}

	createdLeave, err := service.CreateLeaveService(leave)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()}) // 冲突错误
		return
	}

	c.JSON(http.StatusCreated, createdLeave)
}

// ListMyLeavesHandler 获取我自己的请假记录
func ListMyLeavesHandler(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))
	leaves, err := service.ListLeavesService(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve leave records."})
		return
	}
	c.JSON(http.StatusOK, leaves)
}

// DeleteLeaveHandler 删除一条请假记录
func DeleteLeaveHandler(c *gin.Context) {
	leaveID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid leave ID."})
		return
	}
	userID, _ := uuid.Parse(c.GetString("user_id"))

	err = service.DeleteLeaveService(leaveID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Leave record not found or permission denied."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Leave record deleted successfully."})
}
