package google

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

const googleAPICreds = "/.../github.com/connectors-test-client/config/credentials.json"
const googleTokenFile = "/.../github.com/connectors-test-client/config/token.json"
const eventPrefix = "{TEST EVENT}"

// NewGoogleClient creates a new client for making request to Google Calendar API
func NewGoogleClient() (*http.Client, error) {
	b, err := ioutil.ReadFile(googleAPICreds)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	googleConfig, err := google.ConfigFromJSON(b, calendar.CalendarEventsScope)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to parse client secret file to config")
	}
	return getClient(googleConfig), nil
}

// NewGoogleService creates a new service for making request to Google Calendar API
func NewGoogleService(ctx context.Context) (*calendar.Service, error) {
	client, err := NewGoogleClient()
	if err != nil {
		return nil, err
	}
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to create service")
	}
	return srv, nil
}

// CreateEvents create a specified number of events
func CreateEvents(svc *calendar.Service, calendarID string, count int, rate int) []*calendar.Event {
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
		goEvent, err := svc.Events.Insert(calendarID, event).Do()
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
func UpdateEvents(svc *calendar.Service, calendarID string, count int, rate int) []*calendar.Event {
	events, err := ListEvents(svc, calendarID, count)
	if err != nil {
		log.Println("Error fetching events to update")
		return []*calendar.Event{}
	}

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
		goEvent, err := svc.Events.Update(calendarID, event.Id, event).Do()
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

// DeleteTestEvents finds and deletes all test events
func DeleteTestEvents(svc *calendar.Service, calendarID string, count int, rate int) []*calendar.Event {
	events, err := ListEvents(svc, calendarID, count)
	if err != nil {
		log.Println("Error fetching events to delete")
		return []*calendar.Event{}
	}

	delay := int(math.Round(1.0 / (float64(rate) / 60)))
	for _, event := range events {
		err := svc.Events.Delete(calendarID, event.Id).Do()
		if err != nil {
			log.Printf("Error deleting event: %v\n", event.Id)
		} else {
			fmt.Printf("DELETE: %v (%v)\n", event.Summary, event.Start.DateTime)
		}

		time.Sleep(time.Second * time.Duration(delay))
	}

	return events
}

// ListEvents fetch events for given calendar
func ListEvents(svc *calendar.Service, calendarID string, count int) ([]*calendar.Event, error) {
	t := time.Now().Format(time.RFC3339)
	events, err := svc.Events.List(calendarID).ShowDeleted(false).
		SingleEvents(true).TimeMin(t).MaxResults(int64(count)).Q(eventPrefix).OrderBy("startTime").Do()
	if err != nil {
		return nil, err
	}
	// strings.HasPrefix(event.Summary, eventPrefix)
	return events.Items, err
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := googleTokenFile
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}
	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
