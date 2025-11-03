package commons

import (
	"time"

	"gopkg.in/yaml.v3"
)

type SerializableDate struct {
	time.Time
}

type SerializableDuration struct {
	time.Duration
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