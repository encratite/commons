package commons

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"
)

const (
	dateLayout = "2006-01-02"
	minutesLayout = "2006-01-02 15:04"
	timestampLayout = "2006-01-02 15:04:05"
	minutesPerHour = 60
	hoursPerDay = 24
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

func ParseTimeOfDay(timeOfDayString string) (time.Duration, error) {
	pattern := regexp.MustCompile(`^(?:(\d+)d )?(\d+):(\d{2})$`)
	matches := pattern.FindStringSubmatch(timeOfDayString)
	zero := time.Duration(0)
	if matches == nil {
		return zero, fmt.Errorf("Unable to parse duration: %s", timeOfDayString)
	}
	dayMatch := matches[1]
	var days int
	if dayMatch == "" {
		days = 0
	} else {
		d, err := ParseInt(matches[1])
		if err != nil {
			return zero, err
		}
		days = d
	}
	hours, err := ParseInt(matches[2])
	if err != nil {
		return zero, err
	}
	minutes, err := ParseInt(matches[3])
	if err != nil {
		return zero, err
	}
	output := time.Duration(hoursPerDay * days + hours) * time.Hour + time.Duration(minutes) * time.Minute
	return output, nil
}

func MustParseTimeOfDay(timeOfDayString string) time.Duration {
	timeOfDay, err := ParseTimeOfDay(timeOfDayString)
	if err != nil {
		log.Fatal(err)
	}
	return timeOfDay
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

func GetDurationString(duration time.Duration) string {
	hours := int(duration.Hours())
	days := hours / hoursPerDay
	if days > 0 {
		hours %= days
	}
	minutes := int(duration.Minutes()) % minutesPerHour
	if days > 0 {
		return fmt.Sprintf("%dd %02dh %02dm", days, hours, minutes)
	} else {
		return fmt.Sprintf("%02dh %02dm", hours, minutes)
	}
}