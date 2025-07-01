// internal/service/periodic_task_service.go
package service

import (
	"errors"
	"gotasksys/internal/model"
	"gotasksys/internal/repository"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
)

func ListPeriodicTasksService() ([]model.PeriodicTask, error) {
	return repository.ListPeriodicTasks()
}

func CreatePeriodicTaskService(input model.PeriodicTask, creatorID uuid.UUID) (model.PeriodicTask, error) {
	// 校验Cron表达式格式是否正确
	if _, err := cron.ParseStandard(input.CronExpression); err != nil {
		return model.PeriodicTask{}, errors.New("invalid cron expression format")
	}

	input.CreatedByID = creatorID
	if err := repository.CreatePeriodicTask(&input); err != nil {
		return model.PeriodicTask{}, err
	}
	if input.IsActive {
		AddJob(input)
	}
	return input, nil
}

func UpdatePeriodicTaskService(id uuid.UUID, input model.PeriodicTask) (model.PeriodicTask, error) {
	if _, err := cron.ParseStandard(input.CronExpression); err != nil {
		return model.PeriodicTask{}, errors.New("invalid cron expression format")
	}
	pt, err := repository.FindPeriodicTaskByID(id)
	if err != nil {
		return model.PeriodicTask{}, errors.New("periodic task not found")
	}
	pt.Title = input.Title
	pt.Description = input.Description
	pt.CronExpression = input.CronExpression
	pt.DefaultAssigneeID = input.DefaultAssigneeID
	pt.DefaultEffort = input.DefaultEffort
	pt.DefaultPriority = input.DefaultPriority
	pt.DefaultTaskTypeID = input.DefaultTaskTypeID

	if err := repository.UpdatePeriodicTask(&pt); err != nil {
		return model.PeriodicTask{}, err
	}
	RemoveJob(id.String())
	if pt.IsActive {
		AddJob(pt)
	}
	return pt, nil
}

func DeletePeriodicTaskService(id uuid.UUID) error {
	RemoveJob(id.String())
	return repository.DeletePeriodicTask(id)
}

func TogglePeriodicTaskService(id uuid.UUID, isActive bool) error {
	pt, err := repository.FindPeriodicTaskByID(id)
	if err != nil {
		return errors.New("periodic task not found")
	}
	pt.IsActive = isActive

	if err := repository.UpdatePeriodicTask(&pt); err != nil {
		return err
	}
	RemoveJob(id.String())
	if isActive {
		AddJob(pt)
	}
	return nil
}
