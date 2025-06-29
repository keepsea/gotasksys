// internal/repository/dashboard_repository.go
package repository

import (
	"gotasksys/internal/config"
	"gotasksys/internal/model"
)

// CountTasksByStatus 根据状态统计任务数量
func CountTasksByStatus(status string) (int64, error) {
	var count int64
	result := config.DB.Model(&model.Task{}).Where("status = ?", status).Count(&count)
	return count, result.Error
}
