package outlook

import (
	"context"
	"net/http"
	"time"

	"github.com/connectors-test-client/app/outlook/types"
	"github.com/pkg/errors"
	"google.golang.org/api/calendar/v3"
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

// List fetches calendar events
func (svc *Service) List(calenderID string, query string, maxResults int) ([]*calendar.Event, error) {
	msEvents := &types.Events{}
	_, err := svc.client.Get("agg.DeltaLink", msEvents)
	if err != nil {
		return []*calendar.Event{}, errors.Wrap(err, "Unable to perform list")
	}

	events := []*calendar.Event{}
	for _, msEvent := range msEvents.Events {
		event := convertToGoogleEvent(msEvent)
		events = append(events, event)
	}

	return events, nil
}

// Insert inserts calendar event
func (svc *Service) Insert(calendarID string, event *calendar.Event) (*calendar.Event, error) {
	msEvent := convertToOutlookEvent(event)
	_, err := svc.client.Post("agg.DeltaLink", msEvent)
	if err != nil {
		return &calendar.Event{}, errors.Wrap(err, "Unable to perform Insert")
	}
	return event, nil
}

// Update updates existing calendar event
func (svc *Service) Update(calendarID string, event *calendar.Event) (*calendar.Event, error) {
	msEvent := convertToOutlookEvent(event)
	_, err := svc.client.Put("agg.DeltaLink", msEvent)
	if err != nil {
		return &calendar.Event{}, errors.Wrap(err, "Unable to perform Update")
	}
	return event, nil
}

// Delete updates existing calendar event
func (svc *Service) Delete(calendarID string, eventID string) error {
	_, err := svc.client.Delete("agg.DeltaLink")
	if err != nil {
		return errors.Wrap(err, "Unable to perform Update")
	}
	return nil
}

func convertToOutlookEvent(event *calendar.Event) *types.Event {
	const DateTimeTimeZoneFormat = "2006-01-02T15:04:05.9999999"
	t, _ := time.Parse(time.RFC3339, event.Start.DateTime)
	z, _ := t.Zone()
	start := &types.DateTimeTimeZone{
		DateTime: t.Format(DateTimeTimeZoneFormat),
		TimeZone: z,
	}
	t, _ = time.Parse(time.RFC3339, event.End.DateTime)
	z, _ = t.Zone()
	end := &types.DateTimeTimeZone{
		DateTime: t.Format(DateTimeTimeZoneFormat),
		TimeZone: z,
	}
	msEvent := &types.Event{
		ID:      event.Id,
		Subject: event.Summary,
		Start:   start,
		End:     end,
	}
	return msEvent
}

func convertToGoogleEvent(msEvent *types.Event) *calendar.Event {
	start, _ := msEvent.Start.Time()
	end, _ := msEvent.End.Time()
	event := &calendar.Event{
		Id:      msEvent.ID,
		Summary: msEvent.Subject,
		Start: &calendar.EventDateTime{
			DateTime: start.Format(time.RFC3339),
		},
		End: &calendar.EventDateTime{
			DateTime: end.Format(time.RFC3339),
		},
	}
	return event
}
