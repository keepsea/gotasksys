// internal/service/personnel_service.go
package service

import (
	"gotasksys/internal/repository"
	"time"
)

type PersonnelStatus struct {
	UserID         string     `json:"user_id"`
	RealName       string     `json:"real_name"`
	Role           string     `json:"role"`
	CurrentLoad    int        `json:"current_load"`     // 当前总工时
	StatusLight    string     `json:"status_light"`     // 状态灯: idle, normal, busy, overloaded
	ActiveTasks    []TaskInfo `json:"active_tasks"`     // 进行中的任务列表
	HasOverdueTask bool       `json:"has_overdue_task"` // 是否有超期任务
}

type TaskInfo struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}

func GetPersonnelStatusService() ([]PersonnelStatus, error) {
	members, err := repository.FindAllActiveMembers()
	if err != nil {
		return nil, err
	}

	var statuses []PersonnelStatus
	for _, member := range members {
		tasks, err := repository.FindInProgressTasksByAssigneeID(member.ID)
		if err != nil {
			// 如果单个用户查询失败，可以记录日志但继续处理其他用户
			continue
		}

		var currentLoad int
		var activeTasks []TaskInfo
		var hasOverdueTask bool

		for _, task := range tasks {
			currentLoad += task.Effort
			activeTasks = append(activeTasks, TaskInfo{ID: task.ID, Title: task.Title})
			if task.DueDate != nil && task.DueDate.Before(time.Now()) {
				hasOverdueTask = true
			}
		}

		status := PersonnelStatus{
			UserID:         member.ID.String(),
			RealName:       member.RealName,
			Role:           member.Role,
			CurrentLoad:    currentLoad,
			ActiveTasks:    activeTasks,
			HasOverdueTask: hasOverdueTask,
			StatusLight:    calculateStatusLight(currentLoad),
		}
		statuses = append(statuses, status)
	}
	return statuses, nil
}

func calculateStatusLight(load int) string {
	// 阈值定义 (与我们之前设计的一致)
	// 假设8小时工作日，6小时(75%)为繁忙分界线
	if load == 0 {
		return "idle" // 🟢
	} else if load > 0 && load <= 6 {
		return "normal" // 🟡
	} else if load > 6 && load <= 8 {
		return "busy" // 🟠
	} else {
		return "overloaded" // 🔴
	}
}
