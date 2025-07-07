// internal/repository/leave_repository.go
package repository

import (
	"gotasksys/internal/config"
	"gotasksys/internal/model"
	"time"

	"github.com/google/uuid"
)

// CreateLeave 在数据库中创建一条新的请假记录
func CreateLeave(leave *model.Leave) error {
	return config.DB.Create(leave).Error
}

// ListLeavesByUserID 获取一个用户的所有请假记录
func ListLeavesByUserID(userID uuid.UUID) ([]model.Leave, error) {
	var leaves []model.Leave
	err := config.DB.Where("user_id = ?", userID).Order("start_date desc").Find(&leaves).Error
	return leaves, err
}

// DeleteLeave 从数据库中删除一条请假记录
func DeleteLeave(leaveID, userID uuid.UUID) error {
	// 在删除时，也校验用户ID，确保用户只能删除自己的请假记录
	result := config.DB.Where("id = ? AND user_id = ?", leaveID, userID).Delete(&model.Leave{})
	// 如果没有记录被删除，说明该请假记录不属于该用户或不存在
	if result.RowsAffected == 0 {
		return config.DB.Error // 可以返回一个自定义的 "not found or permission denied" 错误
	}
	return result.Error
}

// FindInProgressTasksForUser 获取一个用户所有正在进行的任务
func FindInProgressTasksForUser(userID uuid.UUID) ([]model.Task, error) {
	var tasks []model.Task
	err := config.DB.Where("assignee_id = ? AND status = ?", userID, "in_progress").Find(&tasks).Error
	return tasks, err
}

// ListLeaveDatesInRange 获取一个用户在指定日期范围内的所有请假日期
func ListLeaveDatesInRange(userID uuid.UUID, start, end time.Time) (map[string]bool, error) {
	var leaves []model.Leave
	// 查询所有与[start, end]有交集的请假记录
	err := config.DB.Where("user_id = ? AND start_date <= ? AND end_date >= ?", userID, end, start).Find(&leaves).Error
	if err != nil {
		return nil, err
	}

	leaveDateMap := make(map[string]bool)
	for _, leave := range leaves {
		// 将请假期间的每一天都加入到map中
		current := leave.StartDate
		for !current.After(leave.EndDate) {
			leaveDateMap[current.Format("2006-01-02")] = true
			current = current.AddDate(0, 0, 1)
		}
	}
	return leaveDateMap, nil
}
