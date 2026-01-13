package commons

import (
	"regexp"
	"time"

	"gopkg.in/yaml.v3"
)

type SerializableDate struct {
	time.Time
}

type SerializableDuration struct {
	time.Duration
}

type SerializableRegexp struct {
	*regexp.Regexp
}

func (d *SerializableDate) UnmarshalYAML(value *yaml.Node) error {
	time, err := ParseTime(value.Value)
	if err != nil {
		return err
	}
	d.Time = GetDate(time)
	return nil
}

func (d *SerializableDuration) UnmarshalYAML(value *yaml.Node) error {
	timeOfDay, err := ParseTimeOfDay(value.Value)
	if err != nil {
		return err
	}
	d.Duration = timeOfDay
	return nil
}

func (r *SerializableRegexp) UnmarshalYAML(value *yaml.Node) error {
	pattern, err := regexp.Compile(value.Value)
	if err != nil {
		return err
	}
	r.Regexp = pattern
	return nil
}