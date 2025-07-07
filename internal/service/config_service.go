// internal/service/config_service.go (最终修正版)
package service

import (
	"errors"
	"gotasksys/internal/model"
	"gotasksys/internal/repository"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// --- 系统配置相关 ---

func UpdateSystemConfigService(key, value string) error {
	if key == "global_daily_work_hours" {
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			return errors.New("invalid value for daily work hours, must be a number")
		}
	}
	// 调用 config_repository.go 中的正确函数
	return repository.UpdateSystemConfig(key, value)
}

// --- 法定节假日相关 ---

func ListHolidaysService() ([]model.Holiday, error) {
	return repository.ListAllHolidays()
}
func CreateHolidayService(date time.Time, description string) (model.Holiday, error) {
	holiday := model.Holiday{
		HolidayDate: date,
		Description: description,
	}
	err := repository.CreateHoliday(&holiday)
	return holiday, err
}
func DeleteHolidayService(id uuid.UUID) error {
	return repository.DeleteHoliday(id)
}

// --- 用户个人工时配置相关 ---

// UpdateUserCapacityService 由管理员更新一个用户的个人可用工时
func UpdateUserCapacityService(userID uuid.UUID, hours *float64) error {
	// 【核心修正】调用 user_repository.go 中已有的、更灵活的 UpdateUserProfile 函数
	updates := map[string]interface{}{
		"daily_capacity_hours": hours,
	}
	return repository.UpdateUserProfile(userID, updates)
}
