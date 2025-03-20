package utils

import (
	"fmt"
	"time"
)

var Date *DateClass

type DateClass struct {
	Time time.Time
}

// Format 格式化日期
func (this *DateClass) Format(layout string) string {
	return this.Time.Format(layout)
}

// FormatSeconds 将秒数格式化为友好的时间字符串
func (this *DateClass) FormatSeconds(seconds int) string {

	if seconds < 60 {
		return fmt.Sprintf("%d秒", seconds)
	}

	if seconds < 3600 {

		minutes := seconds / 60
		remain  := seconds % 60

		if remain == 0 {
			return fmt.Sprintf("%d分钟", minutes)
		}
		if remain < 10 {
			return fmt.Sprintf("%d分钟0%d秒", minutes, remain)
		}
		return fmt.Sprintf("%d分钟%d秒", minutes, remain)
	}

	if seconds < 86400 {

		hours   := seconds / 3600
		remain  := seconds % 3600
		minutes := remain  / 60
		remain   = remain % 60

		result := fmt.Sprintf("%d小时", hours)

		if minutes > 0 {
			result += fmt.Sprintf("%d分钟", minutes)
		}

		if remain > 0 {
			if remain < 10 {
				result += fmt.Sprintf("0%d秒", remain)
			} else {
				result += fmt.Sprintf("%d秒", remain)
			}
		}

		return result
	}

	days   := seconds / 86400
	remain := seconds % 86400
	hours  := remain  / 3600

	remain   = remain % 3600
	minutes := remain / 60

	remain  = remain % 60
	result := fmt.Sprintf("%d天", days)

	if hours > 0 {
		result += fmt.Sprintf("%d小时", hours)
	}

	if minutes > 0 {
		result += fmt.Sprintf("%d分钟", minutes)
	}

	if remain > 0 {
		if remain < 10 {
			result += fmt.Sprintf("0%d秒", remain)
		} else {
			result += fmt.Sprintf("%d秒", remain)
		}
	}

	return result
}