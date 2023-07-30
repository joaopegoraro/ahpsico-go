package utils

import "time"

func GetStartOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func GetEndOfDay(t time.Time) time.Time {
	start := GetStartOfDay(t)
	return start.Add(time.Hour * 24)
}
