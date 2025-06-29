// internal/service/admin_service.go
package service

import (
	"gotasksys/internal/model"
	"gotasksys/internal/repository"
)

// 我们暂时把service和repository的函数一一对应
func ListTaskTypesService() ([]model.TaskType, error) {
	return repository.ListTaskTypes()
}

func CreateTaskTypeService(name string) (model.TaskType, error) {
	taskType := model.TaskType{
		Name:      name,
		IsEnabled: true,
	}
	err := repository.CreateTaskType(&taskType)
	return taskType, err
}
