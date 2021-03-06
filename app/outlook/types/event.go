package types

import (
	"time"
)

// Event represents the Microsoft Graph API's Event resource.
type Event struct {
	Calendar                   *Calendar         `json:"calendar,omitempty"`
	Categories                 []string          `json:"categories,omitempty"`
	ChangeKey                  string            `json:"changeKey,omitempty"`
	CreatedDateTime            *time.Time        `json:"createdDateTime,omitempty"`
	End                        *DateTimeTimeZone `json:"end,omitempty"`
	ICalUID                    string            `json:"iCalUId,omitempty"`
	ID                         string            `json:"id,omitempty"`
	OriginalStart              *time.Time        `json:"originalStart,omitempty"`
	OriginalStartTimeZone      string            `json:"originalStartTimeZone,omitempty"`
	ReminderMinutesBeforeStart int32             `json:"reminderMinutesBeforeStart,omitempty"`
	Start                      *DateTimeTimeZone `json:"start,omitempty"`
	Subject                    string            `json:"subject,omitempty"`
	Organizer                  *Recipient        `json:"organizer,omitempty"`
	Body                       *Body             `json:"body,omitempty"`
	Attendees                  []*Recipient      `json:"attendees,omitempty"`
}

// Events represents a list of events from the Microsoft Graph API's
// Event resource.
type Events struct {
	Events []*Event `json:"value,omitempty"`
}

// Body message body of an event
type Body struct {
	ContentType string `json:"contentType,omitempty"`
	Content     string `json:"content,omitempty"`
}
