package tools

import "time"

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
