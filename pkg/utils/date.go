// pkg/utils/date.go

package utils

import "time"

// CalculateWorkingDays 计算两个日期之间的工作日天数 (包含首尾)
func CalculateWorkingDays(start, end time.Time) int {
	// 将时间戳统一到当天的零点，避免小时和分钟的干扰
	startDate := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	endDate := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, end.Location())

	if startDate.After(endDate) {
		return 0
	}

	workingDays := 0
	currentDate := startDate

	for !currentDate.After(endDate) {
		weekday := currentDate.Weekday()
		if weekday != time.Saturday && weekday != time.Sunday {
			workingDays++
		}
		currentDate = currentDate.AddDate(0, 0, 1) // 日期加一天
	}

	return workingDays
}
