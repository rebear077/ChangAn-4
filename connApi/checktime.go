package receive

import (
	"strconv"
	"time"
)

// 时间戳验证
func checkTimeStamp(formatTimeStr string) bool {
	currentTime := time.Now() //当前时间
	m, _ := time.ParseDuration("-10m")
	checktime := currentTime.Add(m)
	checktime_str := checktime.String()
	timeTemplate1 := "2006-01-02 15:04:05"
	checktime1, _ := time.ParseInLocation(timeTemplate1, checktime_str[0:19], time.Local)
	timestamp1, _ := time.ParseInLocation(timeTemplate1, formatTimeStr, time.Local)
	return timestamp1.After(checktime1)
}

func convertimeStamp(timestamp string) string {
	timeUnix, _ := strconv.ParseInt(timestamp, 10, 64)
	formatTime := time.Unix(timeUnix, 0).Format("2006-01-02 15:04:05")
	return formatTime
}
