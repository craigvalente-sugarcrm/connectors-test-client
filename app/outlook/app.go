package outlook

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
	"google.golang.org/api/calendar/v3"
)

const eventPrefix = "{TEST EVENT}"

var config = oauth2.Config{
	ClientID:     "89d631ac-6d06-4b07-992b-ca78eed50787",
	ClientSecret: "M0OF4S?[AdHp8S/8.B[dpiMxab=2+NZc",
	Scopes:       []string{"offline_access", "Calendars.ReadWrite"},
	Endpoint:     microsoft.AzureADEndpoint("tenent"),
	RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",
}

// App for testing Outlook Calendar API
type App struct {
	calendarID string
	svc        *Service
}

// New Creates and retuns a new Outlook Test App
func New(ctx context.Context, token *oauth2.Token, calendarID string) (*App, error) {
	app := &App{
		calendarID: calendarID,
	}

	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))
	svc, err := NewService(ctx, client)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to create service")
	}
	app.svc = svc
	return app, nil
}

// GetTokenFromWeb Request a token from the web, then returns the retrieved token.
func GetTokenFromWeb() {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Authorization Code: \n%v\n", authURL)

	// authCode := ""
	// token, err := config.Exchange(context.TODO(), authCode)
	// if err != nil {
	// 	log.Fatalf("Unable to retrieve token from web: %v", err)
	// } else {
	// 	fmt.Printf("Authorization Code: %v\n", token.AccessToken)
	// 	fmt.Printf("Authorization Code: %v\n", token.RefreshToken)
	// 	fmt.Printf("Authorization Code: %v\n", token.Expiry.Format(time.RFC3339))
	// }
}

// ListEvents fetch events for given calendar
func (app *App) ListEvents(count int) []*calendar.Event {
	events, err := app.svc.List(app.calendarID, eventPrefix, count)
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
		goEvent, err := app.svc.Insert(app.calendarID, event)
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
		goEvent, err := app.svc.Update(app.calendarID, event)
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
		err := app.svc.Delete(app.calendarID, event.Id)
		if err != nil {
			log.Printf("Error deleting event: %v\n", event.Id)
		} else {
			fmt.Printf("DELETE: %v (%v)\n", event.Summary, event.Start.DateTime)
		}

		time.Sleep(time.Second * time.Duration(delay))
	}

	return events
}
