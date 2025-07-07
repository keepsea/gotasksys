// internal/repository/config_repository.go
package repository

import (
	"gotasksys/internal/config"
	"gotasksys/internal/model"
	"time"

	"github.com/google/uuid"
)

// --- 法定节假日 (Holidays) 相关 ---

// ListAllHolidays 获取所有已设置的节假日
func ListAllHolidays() ([]model.Holiday, error) {
	var holidays []model.Holiday
	err := config.DB.Order("holiday_date asc").Find(&holidays).Error
	return holidays, err
}

// CreateHoliday 创建一条新的节假日记录
func CreateHoliday(holiday *model.Holiday) error {
	return config.DB.Create(holiday).Error
}

// DeleteHoliday 删除一条节假日记录
func DeleteHoliday(id uuid.UUID) error {
	return config.DB.Where("id = ?", id).Delete(&model.Holiday{}).Error
}

// GetHolidaysInRange 获取一个日期范围内的所有节假日，返回一个map以方便快速查找
func GetHolidaysInRange(start, end time.Time) (map[string]bool, error) {
	var holidays []model.Holiday
	err := config.DB.Where("holiday_date BETWEEN ? AND ?", start, end).Find(&holidays).Error
	if err != nil {
		return nil, err
	}

	holidayMap := make(map[string]bool)
	for _, h := range holidays {
		holidayMap[h.HolidayDate.Format("2006-01-02")] = true
	}
	return holidayMap, nil
}

// --- 系统配置 (System Configs) 相关 ---

// GetSystemConfigValueByKey 根据键获取一个系统配置的值
func GetSystemConfigValueByKey(key string) (string, error) {
	var cfg model.SystemConfig
	err := config.DB.First(&cfg, "config_key = ?", key).Error
	return cfg.ConfigValue, err
}
func UpdateSystemConfig(key, value string) error {
	return config.DB.Model(&model.SystemConfig{}).Where("config_key = ?", key).Update("config_value", value).Error
}
