// internal/service/personnel_service.go
package service

import (
	"gotasksys/internal/repository"
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

// --- 修改 GetPersonnelStatusService 函数，填充新的数据结构 ---
func GetPersonnelStatusService() ([]PersonnelStatus, error) {
	members, err := repository.FindAllActiveMembers()
	if err != nil {
		return nil, err
	}

	var statuses []PersonnelStatus
	for _, member := range members {
		// 获取进行中的任务和负载 (逻辑不变)
		tasks, _ := repository.FindInProgressTasksByAssigneeID(member.ID)
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

		// --- 核心修改：获取并计算历史绩效分 ---
		var performanceMetrics *PerformanceMetricsDto
		// 我们只为 'executor' 计算绩效分
		if member.Role == "executor" {
			metrics, err := repository.GetPerformanceMetricsForUser(member.ID)
			if err == nil { // 如果查询出错，绩效分部分就为null
				compositeScore := (metrics.AvgTimeliness + metrics.AvgQuality + metrics.AvgCollaboration + metrics.AvgComplexity) / 4.0

				performanceMetrics = &PerformanceMetricsDto{
					CompositeScore:   compositeScore,
					AvgTimeliness:    metrics.AvgTimeliness,
					AvgQuality:       metrics.AvgQuality,
					AvgCollaboration: metrics.AvgCollaboration,
					AvgComplexity:    metrics.AvgComplexity,
				}
			}
		}
		// ------------------------------------

		status := PersonnelStatus{
			UserID:             member.ID.String(),
			RealName:           member.RealName,
			Role:               member.Role,
			CurrentLoad:        currentLoad,
			StatusLight:        calculateStatusLight(currentLoad),
			ActiveTasks:        activeTasks,
			HasOverdueTask:     hasOverdueTask,
			PerformanceMetrics: performanceMetrics, // 赋值
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
