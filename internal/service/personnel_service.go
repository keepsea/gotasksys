// internal/service/personnel_service.go (最终锁定完整版)
package service

import (
	"log"
	"strconv"
	"time"

	"gotasksys/internal/model"
	"gotasksys/internal/repository"
	"gotasksys/pkg/utils"

	"github.com/google/uuid"
)

// --- 数据结构定义 (DTOs) ---

// PerformanceMetricsDto 用于API响应的绩效数据结构
type PerformanceMetricsDto struct {
	CompositeScore   float64 `json:"composite_score"`
	AvgTimeliness    float64 `json:"avg_timeliness"`
	AvgQuality       float64 `json:"avg_quality"`
	AvgCollaboration float64 `json:"avg_collaboration"`
	AvgComplexity    float64 `json:"avg_complexity"`
}

// PersonnelStatus 用于API响应的单个人员的完整状态
type PersonnelStatus struct {
	User               model.User             `json:"user"`
	InProgressTasks    []TaskInfo             `json:"in_progress_tasks"`
	CurrentLoadHours   float64                `json:"current_load_hours"`
	DailyCapacityHours float64                `json:"daily_capacity_hours"`
	LoadPercentage     float64                `json:"load_percentage"`
	StatusLight        string                 `json:"status_light"`
	HasOverdueTask     bool                   `json:"has_overdue_task"`
	PerformanceMetrics *PerformanceMetricsDto `json:"performance_metrics,omitempty"`
}

// TaskInfo 简化的任务信息结构
type TaskInfo struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}

// --- 核心服务函数 ---

// GetPersonnelStatusService 获取所有人员的看板状态数据
func GetPersonnelStatusService() ([]PersonnelStatus, error) {
	// 1. 获取所有需要展示在看板上的成员 (manager, executor)
	members, err := repository.FindAllActiveMembers()
	if err != nil {
		return nil, err
	}

	// 2. 一次性获取全局的每日工时配置，避免在循环中重复查询数据库
	globalDailyHoursStr, err := repository.GetSystemConfigValueByKey("global_daily_work_hours")
	if err != nil {
		globalDailyHoursStr = "8.0" // 如果获取失败，提供一个安全的默认值
	}
	globalDailyHours, _ := strconv.ParseFloat(globalDailyHoursStr, 64)

	var statuses []PersonnelStatus
	today := time.Now()

	for _, member := range members {
		// 3. 获取该成员所有正在进行的任务
		tasks, err := repository.FindInProgressTasksForUser(member.ID)
		if err != nil {
			log.Printf("Failed to get tasks for user %s: %v", member.Username, err)
			continue // 查询单个用户任务失败，跳过该用户，继续处理下一个
		}

		var dailyLoad float64
		var activeTasks []TaskInfo
		var hasOverdueTask bool

		// 4. 【核心算法】循环计算每个任务对今天产生的负载
		for _, task := range tasks {
			activeTasks = append(activeTasks, TaskInfo{ID: task.ID, Title: task.Title})

			// 如果任务没有截止日期，其全部工时都算作“技术债务”，压在今天
			if task.DueDate == nil {
				dailyLoad += float64(task.Effort)
				continue
			}

			// 如果任务已超期，其全部剩余工时也都算作今天的负载
			if task.DueDate.Before(today) {
				dailyLoad += float64(task.Effort)
				hasOverdueTask = true
				continue
			}

			// 对于未超期的任务，进行线性负载分配
			// a. 调用我们强大的新工具，计算从今天到任务截止日期的“实际可用工作日”
			availableDays, err := utils.CalculateAvailableWorkingDays(member.ID, today, *task.DueDate)
			if err != nil {
				log.Printf("Failed to calculate available days for task %d: %v", task.ID, err)
				dailyLoad += float64(task.Effort) // 计算出错则全算
				continue
			}

			if availableDays > 0 {
				// b. 计算该任务每天需要分摊的工时
				dailyEffort := float64(task.Effort) / float64(availableDays)
				dailyLoad += dailyEffort
			} else {
				// 如果可用工作日为0（比如截止日期是今天，但今天是节假日或请假日），则全部工时压在今天
				dailyLoad += float64(task.Effort)
			}
		}

		// 5. 【核心逻辑】确定该成员的“每日可用总工时”
		dailyCapacity := globalDailyHours
		if member.DailyCapacityHours != nil && *member.DailyCapacityHours > 0 {
			dailyCapacity = *member.DailyCapacityHours // 个人配置优先
		}

		// 6. 计算负载百分比和状态灯
		loadPercentage := 0.0
		if dailyCapacity > 0 {
			loadPercentage = (dailyLoad / dailyCapacity) * 100
		}
		statusLight := calculateStatusLight(loadPercentage)

		// 7. 获取历史绩效评分
		performanceMetrics, _ := GetUserPerformanceMetrics(member.ID)

		// 8. 组装最终返回的完整数据
		status := PersonnelStatus{
			User:               member,
			InProgressTasks:    activeTasks,
			CurrentLoadHours:   dailyLoad,
			DailyCapacityHours: dailyCapacity,
			LoadPercentage:     loadPercentage,
			StatusLight:        statusLight,
			HasOverdueTask:     hasOverdueTask,
			PerformanceMetrics: performanceMetrics,
		}
		statuses = append(statuses, status)
	}

	return statuses, nil
}

// --- 辅助函数 ---

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

// GetUserPerformanceMetrics 是一个辅助函数，用于获取并计算单个用户的绩效分
func GetUserPerformanceMetrics(userID uuid.UUID) (*PerformanceMetricsDto, error) {
	// 我们只为 'executor' 计算绩效分
	user, err := repository.FindUserByID(userID)
	if err != nil || user.Role != "executor" {
		return nil, err // 如果不是executor或找不到用户，则不计算
	}

	metrics, err := repository.GetPerformanceMetricsForUser(userID)
	if err != nil {
		return nil, err // 如果查询出错，不计算
	}

	// 只有当有数据时才计算（避免除以0）
	// 此处可以加入更复杂的逻辑，比如完成任务数少于N个则不计算
	compositeScore := (metrics.AvgTimeliness + metrics.AvgQuality + metrics.AvgCollaboration + metrics.AvgComplexity) / 4.0

	return &PerformanceMetricsDto{
		CompositeScore:   compositeScore,
		AvgTimeliness:    metrics.AvgTimeliness,
		AvgQuality:       metrics.AvgQuality,
		AvgCollaboration: metrics.AvgCollaboration,
		AvgComplexity:    metrics.AvgComplexity,
	}, nil
}
