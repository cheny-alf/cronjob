package cronjob

import (
	"strconv"
	"strings"
	"time"
)

/*
IsCronTime
curTime: current time
schedue: crontab expression: seconds/minutes/hours/day/week/month
Determine whether the current time meet crontab expression
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
CheckCronKey
schedule: crontab Expression of split corresponding seconds/minutes/hours/day/week/month
curTime: The current time of seconds/minutes/hours/day/week/month
Determine whether the current time a crontab time format
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
