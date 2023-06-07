package common

import (
	"fmt"
	"strconv"
	"time"
)

func ToDateString(date time.Time) string {
	year, month, day := date.Date()

	monthStr := strconv.Itoa(int(month))
	dayStr := strconv.Itoa(day)
	if int(month) < 10 {
		monthStr = fmt.Sprintf("0%s", monthStr)
	}
	if int(day) < 10 {
		dayStr = fmt.Sprintf("0%s", dayStr)
	}

	return fmt.Sprintf("%v-%s-%s", year, monthStr, dayStr)
}

func ToTimeString(date time.Time) string {
	hour, minute := date.Hour(), date.Minute()

	hourStr := strconv.Itoa(hour)
	minuteStr := strconv.Itoa(minute)
	clockHour := "PM"
	if hour < 10 {
		hourStr = fmt.Sprintf("0%s", hourStr)
		clockHour = "AM"
	}
	if minute < 10 {
		minuteStr = fmt.Sprintf("0%s", minuteStr)
	}

	return fmt.Sprintf("%s:%s %s", hourStr, minuteStr, clockHour)
}
