// pkg/utils/date.go (最终智能版)
package utils

import (
	"gotasksys/internal/repository"
	"time"

	"github.com/google/uuid"
)

// CalculateAvailableWorkingDays 计算一个用户在指定日期范围内的“实际可用工作日”数量
// 它会综合排除：周末、全局法定节假日、个人请假
func CalculateAvailableWorkingDays(userID uuid.UUID, start, end time.Time) (int, error) {
	// 将时间戳统一到当天的零点，避免小时和分钟的干扰
	startDate := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	endDate := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, end.Location())

	if startDate.After(endDate) {
		return 0, nil
	}

	// 1. 从数据库一次性获取此时间段内所有的全局节假日和个人请假
	holidays, err := repository.GetHolidaysInRange(startDate, endDate)
	if err != nil {
		return 0, err
	}
	leaves, err := repository.ListLeaveDatesInRange(userID, startDate, endDate)
	if err != nil {
		return 0, err
	}

	// 2. 循环遍历每一天，进行“三层过滤”
	availableDays := 0
	currentDate := startDate
	for !currentDate.After(endDate) {
		dateStr := currentDate.Format("2006-01-02")
		weekday := currentDate.Weekday()

		// 过滤条件1：不是周六或周日
		isWeekday := weekday != time.Saturday && weekday != time.Sunday
		// 过滤条件2：不在全局节假日列表中
		_, isHoliday := holidays[dateStr]
		// 过滤条件3：不在个人请假列表中
		_, isOnLeave := leaves[dateStr]

		if isWeekday && !isHoliday && !isOnLeave {
			availableDays++
		}

		currentDate = currentDate.AddDate(0, 0, 1) // 日期加一天
	}

	return availableDays, nil
}
