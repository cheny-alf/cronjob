package cronjob

import (
	"strconv"
	"strings"
	"time"
)

/*
	curTime: 当前时间
	schedue: crontab表达式: 秒 分钟 小时 日 月 星期
	判断当前时间是否满足crontab表达式
*/

func IsCronTime(curTime time.Time, schedue string) bool {
	curSec := curTime.Second()
	curMin := curTime.Minute()
	curHour := curTime.Hour()
	curDay := curTime.Day()
	curMonth := int(curTime.Month())
	curWeekday := int(curTime.Weekday())

	schedule := strings.Split(schedue, " ")
	if !CheckCronKey(schedule[0], curSec) || !CheckCronKey(schedule[1], curMin) ||
		!CheckCronKey(schedule[2], curHour) || !CheckCronKey(schedule[3], curDay) ||
		!CheckCronKey(schedule[4], curMonth) || !CheckCronKey(schedule[5], curWeekday) {
		return false
	}
	return true
}

/*
schedule: crontab 表达式拆分的对应 秒/分钟/小时/日/月/星期
curTime: 当前时间的 秒/分钟/小时/日/月/星期
判断当前时间是否命中crontab时间格式
*/
func CheckCronKey(schedule string, curTime int) bool {
	if strings.Contains(schedule, "/") {
		scheduleSlice := strings.Split(schedule, "/")
		scheInt, err := strconv.Atoi(scheduleSlice[1])
		if err != nil {
			return false
		}
		if curTime%scheInt != 0 {
			return false
		}
		schedule = scheduleSlice[0]
	}
	if strings.Contains(schedule, "-") {
		scheduleSlice := strings.Split(schedule, "-")
		min, err1 := strconv.Atoi(scheduleSlice[0])
		max, err2 := strconv.Atoi(scheduleSlice[1])
		if err1 != nil || err2 != nil {
			return false
		}
		if curTime < min || curTime > max {
			return false
		}
		return true
	}
	if schedule != "*" && schedule != strconv.Itoa(curTime) {
		return false
	}
	return true
}
