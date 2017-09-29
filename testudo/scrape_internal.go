package testudo

import (
	"strings"
	"time"
)

func parseDays(dayString string) []time.Weekday {
	days := make([]time.Weekday, 0, 7)
	if strings.Contains(dayString, "M") {
		days = append(days, time.Monday)
	}

	if strings.Contains(dayString, "Tu") {
		days = append(days, time.Tuesday)
	}
	if strings.Contains(dayString, "W") {
		days = append(days, time.Wednesday)
	}
	if strings.Contains(dayString, "Th") {
		days = append(days, time.Thursday)
	}
	if strings.Contains(dayString, "F") {
		days = append(days, time.Friday)
	}

	return days
}
