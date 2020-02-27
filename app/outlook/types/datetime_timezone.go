package types

import (
	"time"

	"github.com/pkg/errors"
)

// DateTimeTimeZoneFormat is RFC3339Nano without the time zone offset (Z or
// ±hh:mm).
const DateTimeTimeZoneFormat = "2006-01-02T15:04:05.9999999"

// DateTimeTimeZone represents the Microsoft Graph API's dateTimeTimeZone
// resource type.
type DateTimeTimeZone struct {
	DateTime string `json:"dateTime,omitempty"`
	TimeZone string `json:"timeZone,omitempty"`
}

// Time converts a DateTimeTimeZone to a time.Time. DateTime is assumed to be in
// UTC if TimeZone is empty or invalid.
func (d *DateTimeTimeZone) Time() (time.Time, error) {
	tz, err := time.LoadLocation(d.TimeZone)
	if err != nil {
		tz = time.UTC
	}

	t, err := time.ParseInLocation(DateTimeTimeZoneFormat, d.DateTime, tz)
	if err != nil {
		return time.Time{}, errors.Wrapf(err, "parse failed")
	}

	return t, nil
}

// UTC converts a DateTimeTimeZone to a time.Time in UTC.
func (d *DateTimeTimeZone) UTC() (time.Time, error) {
	t, err := d.Time()
	if err != nil {
		return t, err
	}

	return t.In(time.UTC), nil
}
