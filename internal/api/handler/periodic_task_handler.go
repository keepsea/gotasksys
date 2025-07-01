// internal/api/handler/periodic_task_handler.go
package handler

import (
	"gotasksys/internal/model"
	"gotasksys/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PeriodicTaskInput struct {
	Title             string     `json:"title" binding:"required"`
	Description       string     `json:"description"`
	CronExpression    string     `json:"cron_expression" binding:"required"`
	DefaultAssigneeID *uuid.UUID `json:"default_assignee_id"`
	DefaultEffort     int        `json:"default_effort" binding:"required"`
	DefaultPriority   string     `json:"default_priority" binding:"required"`
	DefaultTaskTypeID *uuid.UUID `json:"default_task_type_id"`
}

func ListPeriodicTasks(c *gin.Context) {
	tasks, err := service.ListPeriodicTasksService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list periodic tasks"})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func CreatePeriodicTask(c *gin.Context) {
	var input PeriodicTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	creatorID, _ := uuid.Parse(c.GetString("user_id"))

	pt := model.PeriodicTask{
		Title: input.Title, Description: input.Description, CronExpression: input.CronExpression,
		DefaultAssigneeID: input.DefaultAssigneeID, DefaultEffort: input.DefaultEffort,
		DefaultPriority: input.DefaultPriority, DefaultTaskTypeID: input.DefaultTaskTypeID,
	}

	createdTask, err := service.CreatePeriodicTaskService(pt, creatorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, createdTask)
}

func UpdatePeriodicTask(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	var input PeriodicTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pt := model.PeriodicTask{
		Title: input.Title, Description: input.Description, CronExpression: input.CronExpression,
		DefaultAssigneeID: input.DefaultAssigneeID, DefaultEffort: input.DefaultEffort,
		DefaultPriority: input.DefaultPriority, DefaultTaskTypeID: input.DefaultTaskTypeID,
	}

	updatedTask, err := service.UpdatePeriodicTaskService(id, pt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedTask)
}

func DeletePeriodicTask(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	if err := service.DeletePeriodicTaskService(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete periodic task"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Periodic task deleted successfully"})
}

func TogglePeriodicTask(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	var input struct {
		IsActive *bool `json:"is_active" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := service.TogglePeriodicTaskService(id, *input.IsActive); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Periodic task status updated successfully"})
}
