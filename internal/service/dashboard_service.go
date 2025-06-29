// internal/service/dashboard_service.go
package service

import "gotasksys/internal/repository"

type DashboardSummary struct {
	PendingReviewCount int64 `json:"pending_review_count"`
	InPoolCount        int64 `json:"in_pool_count"`
	InProgressCount    int64 `json:"in_progress_count"`
}

// GetDashboardSummaryService 获取驾驶舱的统计数据
func GetDashboardSummaryService() (DashboardSummary, error) {
	var summary DashboardSummary
	var err error

	summary.PendingReviewCount, err = repository.CountTasksByStatus("pending_review")
	if err != nil {
		return summary, err
	}

	summary.InPoolCount, err = repository.CountTasksByStatus("in_pool")
	if err != nil {
		return summary, err
	}

	summary.InProgressCount, err = repository.CountTasksByStatus("in_progress")
	if err != nil {
		return summary, err
	}

	return summary, nil
}
