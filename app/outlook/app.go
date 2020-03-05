package outlook

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/oauth2/microsoft"
	"google.golang.org/api/calendar/v3"
)

const eventPrefix = "--TEST_EVENT--"

var config clientcredentials.Config = clientcredentials.Config{
	ClientID:     "9387a899-c2bb-449a-b532-3e1d7c866806",
	ClientSecret: "CLIENT_SECRET",
	TokenURL:     microsoft.AzureADEndpoint("b43b0d1b-165d-4fbf-a6b7-11583bf5bc63").TokenURL,
	Scopes:       []string{"https://graph.microsoft.com/.default"},
}

// App for testing Outlook Calendar API
type App struct {
	ownerID    string
	calendarID string
	deltaLink  string
	svc        *Service
}

// New Creates and retuns a new Outlook Test App
func New(ctx context.Context, ownerID string, calendarID string) (*App, error) {
	app := &App{
		ownerID:    ownerID,
		calendarID: calendarID,
	}

	startDateTime := time.Now()
	r := rand.Intn(1440 * 30) // 30 days
	endDateTime := time.Now().Add(time.Minute * time.Duration(r))
	app.deltaLink = fmt.Sprintf(
		"/users/%s/calendars/%s/calendarView/delta?%s",
		ownerID,
		calendarID,
		fmt.Sprintf(
			"startDateTime=%s&endDateTime=%s",
			startDateTime.Format(time.RFC3339),
			endDateTime.Format(time.RFC3339),
		),
	)

	svc, err := NewService(ctx, config.Client(ctx))
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to create service")
	}
	app.svc = svc
	return app, nil
}

// ListCalendars fetch calendars for given user
func (app *App) ListCalendars() *[]string {
	calendars, err := app.svc.Calendars(app.ownerID)
	if err != nil {
		log.Printf("Error fetching calendars for user: %v\n", app.ownerID)
		return &[]string{}
	}
	return calendars
}

// ListEvents fetch events for given calendar
func (app *App) ListEvents(count int) []*calendar.Event {
	events, err := app.svc.List(app.deltaLink, eventPrefix, count)
	if err != nil {
		log.Printf("Error fetching events for calendar: %v\n", app.calendarID)
		return []*calendar.Event{}
	}
	return events
}

// CreateEvents create a specified number of events
func (app *App) CreateEvents(count int, rate int) []*calendar.Event {
	var events []*calendar.Event
	for i := 0; i < count; i++ {
		rand.Seed(time.Now().UnixNano())
		r := rand.Intn(1440 * 30) // 30 days
		t := time.Now().Add(time.Minute * time.Duration(r))
		event := &calendar.Event{
			Summary: fmt.Sprintf(eventPrefix+": %v", t.Unix()),
			Start: &calendar.EventDateTime{
				DateTime: t.Format(time.RFC3339),
			},
			End: &calendar.EventDateTime{
				DateTime: t.Add(time.Minute * time.Duration(15)).Format(time.RFC3339),
			},
		}
		events = append(events, event)
	}

	delay := int(math.Round(1.0 / (float64(rate) / 60)))
	for i, event := range events {
		goEvent, err := app.svc.Insert(app.deltaLink, event)
		if err != nil {
			log.Printf("Error creating event: %v\n", event.Summary)
		} else {
			events[i] = goEvent
			fmt.Printf("CREATE: %v (%v)\n", event.Summary, event.Start.DateTime)
		}

		time.Sleep(time.Second * time.Duration(delay))
	}

	return events
}

// UpdateEvents updates a specified number of events
func (app *App) UpdateEvents(count int, rate int) []*calendar.Event {
	events := app.ListEvents(count)
	for i, event := range events {
		t, err := time.Parse(time.RFC3339, event.End.DateTime)
		if err != nil {
			log.Printf("Error parsing End DateTime for event: %v", event.Summary)
			continue
		}
		events[i] = &calendar.Event{
			Id:      event.Id,
			Summary: event.Summary,
			Start:   event.Start,
			End: &calendar.EventDateTime{
				DateTime: t.Add(time.Minute * time.Duration(5)).Format(time.RFC3339),
			},
		}
	}

	delay := int(math.Round(1.0 / (float64(rate) / 60)))
	for i, event := range events {
		goEvent, err := app.svc.Update(app.deltaLink, event)
		if err != nil {
			log.Printf("Error updating event: %v\n", event.Id)
		} else {
			events[i] = goEvent
			fmt.Printf("UPDATE: %v (%v)\n", event.Summary, event.Start.DateTime)
		}

		time.Sleep(time.Second * time.Duration(delay))
	}

	return events
}

// DeleteEvents finds and deletes all test events
func (app *App) DeleteEvents(count int, rate int) []*calendar.Event {
	events := app.ListEvents(count)
	delay := int(math.Round(1.0 / (float64(rate) / 60)))
	for _, event := range events {
		err := app.svc.Delete(app.deltaLink, event.Id)
		if err != nil {
			log.Printf("Error deleting event: %v\n", event.Id)
		} else {
			fmt.Printf("DELETE: %v (%v)\n", event.Summary, event.Start.DateTime)
		}

		time.Sleep(time.Second * time.Duration(delay))
	}

	return events
}
