// internal/service/leave_service.go (最终锁定版)
package service

import (
	"errors"
	"fmt"
	"gotasksys/internal/model"
	"gotasksys/internal/repository"
	"gotasksys/pkg/utils" // 导入我们强大的工具包
	"strings"
	"time"

	"github.com/google/uuid"
)

// ListLeavesService 获取一个用户的所有请假记录
func ListLeavesService(userID uuid.UUID) ([]model.Leave, error) {
	return repository.ListLeavesByUserID(userID)
}

// CreateLeaveService 创建一条新的请假记录，包含最终版的冲突检测
func CreateLeaveService(input model.Leave) (model.Leave, error) {
	if input.EndDate.Before(input.StartDate) {
		return model.Leave{}, errors.New("end date cannot be before start date")
	}

	tasks, err := repository.FindInProgressTasksForUser(input.UserID)
	if err != nil {
		return model.Leave{}, err
	}

	var conflictingTasks []string
	// 遍历每一项进行中的任务
	for _, task := range tasks {
		if task.DueDate == nil {
			continue
		}

		// 【核心升级】使用新工具进行更精确的冲突检测
		// 计算在请假期间，该任务是否还有需要投入的有效工作日
		availableDaysInLeave, err := utils.CalculateAvailableWorkingDays(input.UserID, input.StartDate, *task.DueDate)
		if err != nil {
			continue
		}
		// 任务的“总”可用天数
		totalAvailableDays, _ := utils.CalculateAvailableWorkingDays(input.UserID, time.Now(), *task.DueDate)

		// 如果任务的总可用天数，小于等于它在请假期间的可用天数，说明任务必须在请假期间完成，产生冲突
		if totalAvailableDays > 0 && totalAvailableDays <= availableDaysInLeave {
			conflictingTasks = append(conflictingTasks, fmt.Sprintf("'%s'(截止于%s)", task.Title, task.DueDate.Format("2006-01-02")))
		}
	}

	if len(conflictingTasks) > 0 {
		return model.Leave{}, fmt.Errorf("leave request conflicts with active tasks: %s", strings.Join(conflictingTasks, ", "))
	}

	input.Status = "approved"
	err = repository.CreateLeave(&input)
	return input, err
}

// DeleteLeaveService 删除一个请假记录
func DeleteLeaveService(leaveID, userID uuid.UUID) error {
	return repository.DeleteLeave(leaveID, userID)
}
