// internal/api/handler/admin_handler.go
package handler

import (
	"gotasksys/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// --- 任务类型管理 ---

func ListTaskTypes(c *gin.Context) {
	types, err := service.ListTaskTypesService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list task types"})
		return
	}
	c.JSON(http.StatusOK, types)
}

type CreateTaskTypeInput struct {
	Name string `json:"name" binding:"required"`
}

func CreateTaskType(c *gin.Context) {
	var input CreateTaskTypeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdType, err := service.CreateTaskTypeService(input.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task type"})
		return
	}
	c.JSON(http.StatusCreated, createdType)
}
