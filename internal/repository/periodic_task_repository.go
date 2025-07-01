// internal/repository/periodic_task_repository.go
package repository

import (
	"gotasksys/internal/config"
	"gotasksys/internal/model"

	"github.com/google/uuid"
)

func FindPeriodicTaskByID(id uuid.UUID) (model.PeriodicTask, error) {
	var pt model.PeriodicTask
	err := config.DB.First(&pt, "id = ?", id).Error
	return pt, err
}

func ListPeriodicTasks() ([]model.PeriodicTask, error) {
	var pts []model.PeriodicTask
	err := config.DB.Order("created_at desc").Find(&pts).Error
	return pts, err
}

func CreatePeriodicTask(pt *model.PeriodicTask) error {
	return config.DB.Create(pt).Error
}

func UpdatePeriodicTask(pt *model.PeriodicTask) error {
	return config.DB.Save(pt).Error
}

func DeletePeriodicTask(id uuid.UUID) error {
	return config.DB.Where("id = ?", id).Delete(&model.PeriodicTask{}).Error
}
