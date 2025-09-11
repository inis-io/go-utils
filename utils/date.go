package utils

import (
	"fmt"
	"time"
	
	"github.com/spf13/cast"
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
func (this *DateClass) FormatSeconds(seconds any) string {

	value := cast.ToInt(seconds)

	if value < 60 {
		return fmt.Sprintf("%d秒", value)
	}

	if value < 3600 {

		minutes := value / 60
		remain  := value % 60

		if remain == 0 {
			return fmt.Sprintf("%d分钟", minutes)
		}
		if remain < 10 {
			return fmt.Sprintf("%d分钟0%d秒", minutes, remain)
		}
		return fmt.Sprintf("%d分钟%d秒", minutes, remain)
	}

	if value < 86400 {

		hours   := value / 3600
		remain  := value % 3600
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

	days   := value / 86400
	remain := value % 86400
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