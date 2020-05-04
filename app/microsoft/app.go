package microsoft

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/connectors-test-client/app/microsoft/types"
	"github.com/pkg/errors"
	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/oauth2/microsoft"
)

const eventPrefix = "--TEST_EVENT--"

var config clientcredentials.Config = clientcredentials.Config{
	ClientID:     "9387a899-c2bb-449a-b532-3e1d7c866806",
	ClientSecret: "aZp@W1CipG=-pU/MpTOBT3Y6GHyiBn20",
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
func (app *App) ListEvents(count int) []*types.Event {
	events, err := app.svc.List(app.ownerID, eventPrefix, count)
	if err != nil {
		log.Printf("Error fetching events for owner: %v\n", app.ownerID)
		return []*types.Event{}
	}
	return events
}

// CreateEvents create a specified number of events
func (app *App) CreateEvents(count int) []*types.Event {
	var events []*types.Event
	for i := 0; i < count; i++ {
		rand.Seed(time.Now().UnixNano())
		r := rand.Intn(1440 * 30) // 30 days
		t := time.Now().UTC().Add(time.Minute * time.Duration(r))
		event := &types.Event{
			Subject: fmt.Sprintf(eventPrefix+": %v", t.Unix()),
			Start: &types.DateTimeTimeZone{
				DateTime: t.Format(types.DateTimeTimeZoneFormat),
				TimeZone: t.Location().String(),
			},
			End: &types.DateTimeTimeZone{
				DateTime: t.Add(time.Minute * time.Duration(30)).Format(types.DateTimeTimeZoneFormat),
				TimeZone: t.Location().String(),
			},
			Body: &types.Body{
				ContentType: "HTML",
				Content:     fmt.Sprintf("Project Andaman eest event for user: %s", app.ownerID),
			},
		}
		events = append(events, event)
	}

	// for i, event := range events {
	// 	goEvent, err := app.svc.Insert(app.ownerID, event)
	// 	if err != nil {
	// 		log.Printf("Error creating event: %v\n", event.Subject)
	// 	} else {
	// 		events[i] = goEvent
	// 		fmt.Printf("CREATE: %v (%v)\n", event.Subject, event.Start.DateTime)
	// 	}
	// }

	return events
}

// UpdateEvents updates a specified number of events
func (app *App) UpdateEvents(count int, rate int) []*types.Event {
	events := app.ListEvents(count)
	for i, event := range events {
		t, err := time.Parse(time.RFC3339, event.End.DateTime)
		if err != nil {
			log.Printf("Error parsing End DateTime for event: %v", event.Subject)
			continue
		}
		events[i] = &types.Event{
			ID:      event.ID,
			Subject: event.Subject,
			Start:   event.Start,
			End: &types.DateTimeTimeZone{
				DateTime: t.Add(time.Minute * time.Duration(30)).Format(types.DateTimeTimeZoneFormat),
				TimeZone: t.Location().String(),
			},
		}
	}
	app.UpdateTheseEvents(events, rate)
	return events
}

// UpdateTheseEvents updates the provided events
func (app *App) UpdateTheseEvents(events []*types.Event, rate int) {
	for i, event := range events {
		goEvent, err := app.svc.Update(app.ownerID, event)
		if err != nil {
			log.Printf("Error updating event: %v\n", event.ID)
		} else {
			events[i] = goEvent
			fmt.Printf("UPDATE: %v (%v)\n", event.Subject, event.Start.DateTime)
		}
	}
}

// DeleteEvents finds and deletes all test events
func (app *App) DeleteEvents(count int, rate int) []*types.Event {
	events := app.ListEvents(count)
	app.DeleteTheseEvents(events, rate)
	return events
}

// DeleteTheseEvents deletes the provided events
func (app *App) DeleteTheseEvents(events []*types.Event, rate int) {
	for _, event := range events {
		err := app.svc.Delete(app.ownerID, event.ID)
		if err != nil {
			log.Printf("Error deleting event: %v\n", event.ID)
		} else {
			fmt.Println("EventID", event.ID)
			fmt.Printf("DELETE: %v (%v)\n", event.Subject, event.Start.DateTime)
		}
	}
}
