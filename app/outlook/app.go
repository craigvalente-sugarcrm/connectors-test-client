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
	ownerID string
	svc     *Service
}

// New Creates and retuns a new Outlook Test App
func New(ctx context.Context, ownerID string) (*App, error) {
	app := &App{
		ownerID: ownerID,
	}

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
	events, err := app.svc.List(app.ownerID, eventPrefix, count)
	if err != nil {
		log.Printf("Error fetching events for owner: %v\n", app.ownerID)
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
		t := time.Now().UTC().Add(time.Minute * time.Duration(r))
		event := &calendar.Event{
			Summary: fmt.Sprintf(eventPrefix+": %v", t.Unix()),
			Start: &calendar.EventDateTime{
				DateTime: t.Format(time.RFC3339),
			},
			End: &calendar.EventDateTime{
				DateTime: t.Add(time.Minute * time.Duration(15)).Format(time.RFC3339),
			},
			Description: fmt.Sprintf("Project Andaman eest event for user: %s", app.ownerID),
		}
		events = append(events, event)
	}

	delay := int(math.Round(1.0 / (float64(rate) / 60)))
	for i, event := range events {
		goEvent, err := app.svc.Insert(app.ownerID, event)
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
	app.UpdateTheseEvents(events, rate)
	return events
}

// UpdateTheseEvents updates the provided events
func (app *App) UpdateTheseEvents(events []*calendar.Event, rate int) {
	delay := int(math.Round(1.0 / (float64(rate) / 60)))
	for i, event := range events {
		goEvent, err := app.svc.Update(app.ownerID, event)
		if err != nil {
			log.Printf("Error updating event: %v\n", event.Id)
		} else {
			events[i] = goEvent
			fmt.Printf("UPDATE: %v (%v)\n", event.Summary, event.Start.DateTime)
		}

		time.Sleep(time.Second * time.Duration(delay))
	}
}

// DeleteEvents finds and deletes all test events
func (app *App) DeleteEvents(count int, rate int) []*calendar.Event {
	events := app.ListEvents(count)
	app.DeleteTheseEvents(events, rate)
	return events
}

// DeleteTheseEvents deletes the provided events
func (app *App) DeleteTheseEvents(events []*calendar.Event, rate int) {
	delay := int(math.Round(1.0 / (float64(rate) / 60)))
	for _, event := range events {
		err := app.svc.Delete(app.ownerID, event.Id)
		if err != nil {
			log.Printf("Error deleting event: %v\n", event.Id)
		} else {
			fmt.Printf("DELETE: %v (%v)\n", event.Summary, event.Start.DateTime)
		}

		time.Sleep(time.Second * time.Duration(delay))
	}
}
