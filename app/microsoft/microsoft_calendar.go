package microsoft

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/connectors-test-client/app/microsoft/types"
	"github.com/pkg/errors"
)

// Service a MS Outlook Calendar Service
type Service struct {
	ctx    context.Context
	client *Client
}

// NewService returns a MS Outlook Calendar Service
func NewService(ctx context.Context, client *http.Client) (*Service, error) {
	msClient := NewClient(ctx, client)
	svc := &Service{
		ctx:    ctx,
		client: msClient,
	}
	return svc, nil
}

// Calendars fetches Calendars
func (svc *Service) Calendars(ownerID string) (*[]string, error) {
	url := fmt.Sprintf("/users/%s/calendars", ownerID)
	msCalendars := &types.Calendars{}
	_, err := svc.client.Get(url, msCalendars)
	if err != nil {
		return &[]string{}, errors.Wrap(err, "Unable to perform list")
	}

	calendars := []string{}
	for _, msCal := range msCalendars.Calendars {
		calendars = append(calendars, fmt.Sprintf("%s: %s", msCal.Name, msCal.ID))
	}

	return &calendars, nil
}

// List fetches calendar events
func (svc *Service) List(ownerID string, search string, maxResults int) ([]*types.Event, error) {
	events := &types.Events{}
	url := fmt.Sprintf("%s&$top=%v&$filter=startswith(subject,'%s')", deltaLink(ownerID), maxResults, search)
	// fmt.Println(url)
	_, err := svc.client.Get(url, events)
	if err != nil {
		return []*types.Event{}, errors.Wrap(err, "Unable to perform list")
	}

	return events.Events, nil
}

// Insert inserts calendar event
func (svc *Service) Insert(ownerID string, event *types.Event) (*types.Event, error) {
	// bytes, _ := json.Marshal(msEvent)
	// fmt.Printf(string(bytes))
	resp, err := svc.client.Post(eventsURL(ownerID), event)
	if err != nil {
		return &types.Event{}, errors.Wrap(err, "Unable to perform Insert")
	}
	if resp.StatusCode != http.StatusCreated {
		err := errors.New(fmt.Sprintf("Request failed: %v\n", resp.StatusCode))
		body, _ := ReadHTTPResponse(resp)
		err = errors.Wrap(err, string(body))
		return &types.Event{}, err
	}
	return event, nil
}

// Update updates existing calendar event
func (svc *Service) Update(ownerID string, event *types.Event) (*types.Event, error) {
	_, err := svc.client.Patch(fmt.Sprintf("%s/%s", eventsURL(ownerID), event.ID), event)
	if err != nil {
		return &types.Event{}, errors.Wrap(err, "Unable to perform Update")
	}
	return event, nil
}

// Delete updates existing calendar event
func (svc *Service) Delete(ownerID string, eventID string) error {
	_, err := svc.client.Delete(fmt.Sprintf("%s/%s", eventsURL(ownerID), eventID))
	if err != nil {
		return errors.Wrap(err, "Unable to perform Update")
	}
	return nil
}

func calendarURL(ownerID string) string {
	return fmt.Sprintf("/users/%s/calendar", ownerID)
}

func eventsURL(ownerID string) string {
	return fmt.Sprintf("%s/events", calendarURL(ownerID))
}

func deltaLink(ownerID string) string {
	startDateTime := time.Now().UTC()
	endDateTime := startDateTime.Add(time.Minute * time.Duration(1440*30)) // 30 days
	return fmt.Sprintf(
		"%s/calendarView/?%s",
		calendarURL(ownerID),
		fmt.Sprintf(
			"startDateTime=%s&endDateTime=%s",
			startDateTime.Format(time.RFC3339),
			endDateTime.Format(time.RFC3339),
		),
	)
}

// func convertToOutlookEvent(event *calendar.Event) *types.Event {
// 	const DateTimeTimeZoneFormat = "2006-01-02T15:04:05.9999999"
// 	t, _ := time.Parse(time.RFC3339, event.Start.DateTime)
// 	start := &types.DateTimeTimeZone{
// 		DateTime: t.Format(DateTimeTimeZoneFormat),
// 		TimeZone: t.Location().String(),
// 	}
// 	t, _ = time.Parse(time.RFC3339, event.End.DateTime)
// 	end := &types.DateTimeTimeZone{
// 		DateTime: t.Format(DateTimeTimeZoneFormat),
// 		TimeZone: t.Location().String(),
// 	}
// 	msEvent := &types.Event{
// 		ID:      event.Id,
// 		Subject: event.Summary,
// 		Start:   start,
// 		End:     end,
// 		Body: &types.Body{
// 			ContentType: "HTML",
// 			Content:     event.Description,
// 		},
// 	}
// 	return msEvent
// }

// func convertToGoogleEvent(msEvent *types.Event) *calendar.Event {
// 	start, _ := msEvent.Start.Time()
// 	end, _ := msEvent.End.Time()
// 	event := &calendar.Event{
// 		Id:      msEvent.ID,
// 		Summary: msEvent.Subject,
// 		Start: &calendar.EventDateTime{
// 			DateTime: start.Format(time.RFC3339),
// 		},
// 		End: &calendar.EventDateTime{
// 			DateTime: end.Format(time.RFC3339),
// 		},
// 	}
// 	return event
// }
