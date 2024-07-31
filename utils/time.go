package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ParseTimeString 解析时间字符串，返回小时和分钟
func ParseTimeString(timeStr string) (int, int, error) {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid time format")
	}
	hour, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid hour format")
	}
	minute, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid minute format")
	}
	if hour < 0 || hour > 23 {
		return 0, 0, fmt.Errorf("hour must be between 0 and 23")
	}
	if minute < 0 || minute > 59 {
		return 0, 0, fmt.Errorf("minute must be between 0 and 59")
	}
	return hour, minute, nil
}

// ValidateTime 验证时间是否在合法范围内
func ValidateTime(hour, minute int) error {
	if hour < 0 || hour > 23 {
		return fmt.Errorf("hour must be between 0 and 23")
	}
	if minute < 0 || minute > 59 {
		return fmt.Errorf("minute must be between 0 and 59")
	}
	return nil
}

// CompareTimes 比较开始时间和结束时间
func CompareTimes(startHour, startMinute, endHour, endMinute int) error {
	startTime := time.Date(0, 1, 1, startHour, startMinute, 0, 0, time.Local)
	endTime := time.Date(0, 1, 1, endHour, endMinute, 0, 0, time.Local)
	if startTime.After(endTime) {
		return fmt.Errorf("start time cannot be after end time")
	}
	return nil
}

// CheckTime 检查时间
func CheckTime(timeStr string) (int, int, error) {
	var (
		hour   int
		minute int
		err    error
	)
	// 数据格式
	if hour, minute, err = ParseTimeString(timeStr); err != nil {
		return 0, 0, err
	}
	return hour, minute, nil
}

// ConvertToCronSpec 转换为表达式
func ConvertToCronSpec(hour, minute int) string {
	return fmt.Sprintf("%d %d * * ?", minute, hour)
}

// IsCurrentTimeInRange 判断系统当前时间是否在时间范围内
func IsCurrentTimeInRange(start, end string) (bool, error) {
	// 解析 startTime 和 endTime
	startTime, err := time.Parse("15:04", start)
	if err != nil {
		return false, fmt.Errorf("failed to parse startTime: %v", err)
	}
	endTime, err := time.Parse("15:04", end)
	if err != nil {
		return false, fmt.Errorf("failed to parse endTime: %v", err)
	}

	// 获取当前时间的小时和分钟
	now := time.Now()
	currentTime, err := time.Parse("15:04", now.Format("15:04"))
	if err != nil {
		return false, fmt.Errorf("failed to parse currentTime: %v", err)
	}

	// 判断当前时间是否在范围内
	if currentTime.After(startTime) && currentTime.Before(endTime) {
		return true, nil
	}
	return false, nil
}
