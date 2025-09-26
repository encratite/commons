package commons

import (
	"fmt"
	"log"
	"time"
)

const (
	dateLayout = "2006-01-02"
	minutesLayout = "2006-01-02 15:04"
	timestampLayout = "2006-01-02 15:04:05"
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