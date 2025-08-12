package tools

import (
	"time"
)

// 获取下一个周五日期字符串
func GetNextFriday() string {
	// 获取当前时间
	now := time.Now()

	// 计算下一个周五的日期
	daysUntilFriday := (5 - int(now.Weekday()) + 7) % 7
	if daysUntilFriday == 0 {
		daysUntilFriday = 7 // 如果今天是周五，则获取下一个周五
	}
	nextFriday := now.AddDate(0, 0, daysUntilFriday)

	// 返回格式化后的日期字符串
	return nextFriday.Format("2006-01-02")
}

// 获取下一个家族战日期
func GetNextBattleDate() string {
	// 获取当前时间
	now := time.Now()

	// 计算今天是星期几
	weekday := int(now.Weekday())

	// 如果是周六晚上 8 点之后，返回周日的日期
	if weekday == 6 && now.Hour() >= 20 {
		return now.AddDate(0, 0, 1).Format("2006-01-02")
	}

	// 如果是周日晚上 8 点之后，返回下周六的日期
	if weekday == 0 && now.Hour() >= 20 {
		return now.AddDate(0, 0, 6).Format("2006-01-02")
	}

	// 如果是周六且时间未到晚上 8 点，返回今天的日期
	if weekday == 6 {
		return now.Format("2006-01-02")
	}

	// 如果是周日且时间未到晚上 8 点，返回今天的日期
	if weekday == 0 {
		return now.Format("2006-01-02")
	}

	// 其他情况，返回下一个周六的日期
	daysUntilSaturday := (6 - weekday + 7) % 7
	return now.AddDate(0, 0, daysUntilSaturday).Format("2006-01-02")
}
