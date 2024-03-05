package utils

import (
	"fmt"
	"time"
)

const DefaultDateLayout = "2006-01-02 15:04:05"
const DefaultDayLayout = "2006-01-02"
const DefaultMonthLayout = "2006-01"

// FormatDate DefaultDateLayout = "2006-01-02 15:04:05"
func FormatDate(day string, hour int) string {
	if hour < 10 {
		return fmt.Sprintf("%s %0d:00:00", day, hour)
	} else {
		return fmt.Sprintf("%s %d:00:00", day, hour)
	}
}

func ConvertStringToTime(dateStr string, layout string) time.Time {
	date, err := time.Parse(layout, dateStr)
	if err != nil {
		panic(err)
	}
	return date
}

func ParseDayToDateWithNow(day string) time.Time {
	return ConvertStringToTime(fmt.Sprintf("%s-%s", time.Now().Format(DefaultMonthLayout), day), DefaultDayLayout)
}

func ParseDayToDateWithBase(day string, baseTime time.Time) time.Time {
	return ConvertStringToTime(fmt.Sprintf("%s-%s", baseTime.Format(DefaultMonthLayout), day), DefaultDayLayout)
}

func GetCurrentTimestamp(timeLocation string) int64 {
	loc, _ := time.LoadLocation(timeLocation)
	timestamp := time.Now().In(loc).UnixNano() / 1e6
	return timestamp
}

func GetFirstDateOfMonth(date time.Time) time.Time {
	return date.AddDate(0, 0, -date.Day()+1)
}

func GetLastDateOfMonth(date time.Time) time.Time {
	return GetFirstDateOfMonth(date).AddDate(0, 1, -1)
}

func GetDaysOfMonth(date time.Time) int {
	return GetLastDateOfMonth(date).Day()
}

func GetFirstDateOfYear(date time.Time) time.Time {
	return time.Date(date.Year(), time.January, 1, 0, 0, 0, 0, date.Location())
}

func GetLastDateOfYear(date time.Time) time.Time {
	return GetFirstDateOfYear(date).AddDate(1, 0, -1)
}

func GetDaysOfYear(date time.Time) int {
	return GetLastDateOfYear(date).YearDay()
}

func ParseDateFromString(timeStr string) time.Time {
	return ConvertStringToTime(timeStr, DefaultDateLayout)
}
