package commons

import (
	"log"
	"time"
)

type StopWatch struct {
	start time.Time
}

func NewStopWatch() StopWatch {
	start := time.Now()
	stopWatch := StopWatch{
		start: start,
	}
	return stopWatch
}

func (s *StopWatch) Stop(message string) {
	end := time.Now()
	duration := end.Sub(s.start)
	log.Printf("%s in %.1f s", message, duration.Seconds())
}