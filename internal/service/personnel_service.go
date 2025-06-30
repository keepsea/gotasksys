// internal/service/personnel_service.go
package service

import (
	"gotasksys/internal/repository"
	"gotasksys/pkg/utils" // <-- 新增导入我们自己的工具包
	"time"
)

// --- 新增：用于API响应的、更丰富的绩效数据结构 ---
type PerformanceMetricsDto struct {
	CompositeScore   float64 `json:"composite_score"`
	AvgTimeliness    float64 `json:"avg_timeliness"`
	AvgQuality       float64 `json:"avg_quality"`
	AvgCollaboration float64 `json:"avg_collaboration"`
	AvgComplexity    float64 `json:"avg_complexity"`
}

// --- 修改 PersonnelStatus 结构体，使其包含上面这个新结构 ---
type PersonnelStatus struct {
	UserID             string                 `json:"user_id"`
	RealName           string                 `json:"real_name"`
	Role               string                 `json:"role"`
	CurrentLoad        int                    `json:"current_load"`
	StatusLight        string                 `json:"status_light"`
	ActiveTasks        []TaskInfo             `json:"active_tasks"`
	HasOverdueTask     bool                   `json:"has_overdue_task"`
	PerformanceMetrics *PerformanceMetricsDto `json:"performance_metrics,omitempty"` // 使用指针，如果没有数据则为null
}

type TaskInfo struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}

// GetPersonnelStatusService 获取人员看板状态数据 (线性负载分配版)
func GetPersonnelStatusService() ([]PersonnelStatus, error) {
	members, err := repository.FindAllActiveMembers()
	if err != nil {
		return nil, err
	}

	var statuses []PersonnelStatus
	today := time.Now()

	for _, member := range members {
		tasks, _ := repository.FindInProgressTasksByAssigneeID(member.ID)

		var dailyLoad float64 // 使用浮点数以提高精度
		var activeTasks []TaskInfo
		var hasOverdueTask bool

		for _, task := range tasks {
			activeTasks = append(activeTasks, TaskInfo{ID: task.ID, Title: task.Title})

			if task.DueDate == nil {
				// 如果任务没有截止日期，为简化处理，我们将其全部工时算作当天负载
				dailyLoad += float64(task.Effort)
				continue
			}

			if task.DueDate.Before(today) {
				// 如果任务已超期，其所有剩余工时都压在今天
				dailyLoad += float64(task.Effort)
				hasOverdueTask = true
				continue
			}

			// --- 核心：线性分配逻辑 ---
			workingDays := utils.CalculateWorkingDays(today, *task.DueDate)
			if workingDays > 0 {
				dailyShare := float64(task.Effort) / float64(workingDays)
				dailyLoad += dailyShare
			} else { // 如果截止日期就是今天，且今天不是周末
				dailyLoad += float64(task.Effort)
			}
			// -------------------------
		}

		// 获取历史绩效分 (逻辑不变)
		var performanceMetrics *PerformanceMetricsDto
		if member.Role == "executor" { /* ... */
		}

		status := PersonnelStatus{
			UserID:             member.ID.String(),
			RealName:           member.RealName,
			Role:               member.Role,
			CurrentLoad:        int(dailyLoad), // 返回时可以取整
			StatusLight:        calculateStatusLight(dailyLoad),
			ActiveTasks:        activeTasks,
			HasOverdueTask:     hasOverdueTask,
			PerformanceMetrics: performanceMetrics,
		}
		statuses = append(statuses, status)
	}
	return statuses, nil
}

// calculateStatusLight 修改为接收浮点数
func calculateStatusLight(load float64) string {
	if load == 0 {
		return "idle"
	} else if load > 0 && load <= 6.0 {
		return "normal"
	} else if load > 6.0 && load <= 8.0 {
		return "busy"
	} else {
		return "overloaded"
	}
}
