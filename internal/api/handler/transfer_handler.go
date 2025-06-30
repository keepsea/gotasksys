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
