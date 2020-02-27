package app

import "google.golang.org/api/calendar/v3"

// App interface for Google or Outlook calendar test app
type App interface {
	ListEvents(count int) []*calendar.Event
	CreateEvents(count int, rate int) []*calendar.Event
	UpdateEvents(count int, rate int) []*calendar.Event
	DeleteEvents(count int, rate int) []*calendar.Event
}
