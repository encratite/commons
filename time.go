package commons

import (
	"fmt"
	"log"
	"strings"
	"time"
)

const (
	dateLayout = "2006-01-02"
	minutesLayout = "2006-01-02 15:04"
	timestampLayout = "2006-01-02 15:04:05"
	minutesPerHour = 60
)

func GetDate(timestamp time.Time) time.Time {
	return time.Date(
		timestamp.Year(),
		timestamp.Month(),
		timestamp.Day(),
		0,
		0,
		0,
		0,
		timestamp.Location(),
	)
}

func GetHourTimestamp(timestamp time.Time) time.Time {
	return time.Date(
		timestamp.Year(),
		timestamp.Month(),
		timestamp.Day(),
		timestamp.Hour(),
		0,
		0,
		0,
		timestamp.Location(),
	)
}

func GetDateString(date time.Time) string {
	return date.Format(dateLayout)
}

func GetTimeString(timestamp time.Time) string {
	return timestamp.Format(timestampLayout)
}

func ParseTime(timeString string) (time.Time, error) {
	layouts := []string{
		time.RFC3339,
		timestampLayout,
		minutesLayout,
		dateLayout,
	}
	for _, layout := range layouts {
		output, err := time.Parse(layout, timeString)
		if err != nil {
			continue
		}
		return output, nil
	}
	err := fmt.Errorf("Unable to parse time string: %s", timeString)
	return time.Time{}, err
}

func MustParseTime(timeString string) time.Time {
	timestamp, err := ParseTime(timeString)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return timestamp
}

func ParseWeekday(weekdayString string) (time.Weekday, error) {
	switch strings.ToLower(weekdayString) {
	case "monday":
		return time.Monday, nil
	case "tuesday":
		return time.Tuesday, nil
	case "wednesday":
		return time.Wednesday, nil
	case "thursday":
		return time.Thursday, nil
	case "friday":
		return time.Friday, nil
	case "saturday":
		return time.Saturday, nil
	case "sunday":
		return time.Sunday, nil
	}
	err := fmt.Errorf("Invalid weekday string: %s", weekdayString)
	return time.Monday, err
}

func MustParseWeekday(weekdayString string) time.Weekday {
	weekday, err := ParseWeekday(weekdayString)
	if err != nil {
		log.Fatal(err)
	}
	return weekday
}

func GetTimeOfDay(timestamp time.Time) time.Duration {
	duration := time.Duration(timestamp.Hour()) * time.Hour + time.Duration(timestamp.Minute()) * time.Minute
	return duration
}

func GetTimeOfDayString(timeOfDay time.Duration) string {
	hours := int(timeOfDay.Hours())
	minutes := int(timeOfDay.Minutes()) % minutesPerHour
	output := fmt.Sprintf("%02d:%02d", hours, minutes)
	return output
}