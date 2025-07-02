// internal/api/handler/transfer_handler.go
package handler

import (
	"gotasksys/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type InitiateTransferInput struct {
	NewAssigneeID string `json:"new_assignee_id" binding:"required,uuid"`
	EffortSpent   int    `json:"effort_spent" binding:"gte=0"`
}

func InitiateTransfer(c *gin.Context) {
	taskID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	initiatorID, _ := uuid.Parse(c.GetString("user_id"))

	var input InitiateTransferInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newAssigneeID, _ := uuid.Parse(input.NewAssigneeID)

	transfer, err := service.InitiateTransferService(uint(taskID), initiatorID, newAssigneeID, input.EffortSpent)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transfer initiated successfully.", "transfer_id": transfer.ID})
}

func AcceptTransfer(c *gin.Context) {
	transferID, _ := uuid.Parse(c.Param("transfer_id"))
	responderID, _ := uuid.Parse(c.GetString("user_id"))

	if err := service.RespondToTransferService(transferID, responderID, "accept"); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Transfer accepted."})
}

func RejectTransfer(c *gin.Context) {
	transferID, _ := uuid.Parse(c.Param("transfer_id"))
	responderID, _ := uuid.Parse(c.GetString("user_id"))

	if err := service.RespondToTransferService(transferID, responderID, "reject"); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Transfer rejected."})
}

// CancelTransfer 允许发起人取消一个待处理的转交请求
func CancelTransfer(c *gin.Context) {
	transferID, err := uuid.Parse(c.Param("transfer_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transfer ID"})
		return
	}

	// 从中间件获取当前操作用户的ID
	initiatorID, _ := uuid.Parse(c.GetString("user_id"))

	err = service.CancelTransferService(transferID, initiatorID)
	if err != nil {
		// 根据错误类型返回不同响应
		if err.Error() == "permission denied: you are not the initiator of this transfer" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transfer cancelled successfully."})
}
