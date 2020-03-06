package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/connectors-test-client/app/google"
	"github.com/connectors-test-client/app/outlook"
	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
)

var microsoftUsers = []string{
	"c2d82a71-98d8-4645-9e88-e5531ed1ca02", // Greg
	"627f153d-7dfc-44f5-9f66-4393de2bfe54", // Haarika
	"2f45be53-47dd-46df-9e47-8c442aaf48ed", // Tim
}

type settings struct {
	userID    string
	numEvents int
	rate      int
}

func main() {

	args := os.Args[1:]
	if len(args) != 2 {
		log.Fatal("Missing required input params")
		return
	}

	numMinutes, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatal(err, "1st param must be an int")
	}
	eventsPerMinute, err := strconv.Atoi(args[1])
	if err != nil {
		log.Fatal(err, "2nd param must be an int")
	}

	ch := make(chan string, len(microsoftUsers))
	for _, userID := range microsoftUsers {
		app, err := initOutlookApp(userID)
		if err != nil {
			log.Fatal(err)
		}
		s := &settings{
			userID:    userID,
			numEvents: numMinutes * eventsPerMinute,
			rate:      eventsPerMinute,
		}
		go runOutlookTest(s, app, ch)

		// calendars := app.ListCalendars()
		// for _, cal := range *calendars {
		// 	fmt.Println(userID)
		// 	fmt.Println(cal)
		// }
		// ch <- userID
	}

	for range microsoftUsers {
		fmt.Printf("Completed test for user: %s\n", <-ch)
	}

	// var app app.App
	// app, err := initGoogleApp()

	// Helper func to request authCode and exchange it for token
	// google.GetTokenFromWeb()

	// events := app.ListEvents(10)
	// fmt.Printf("Fetched %v events", len(events))

	// events := app.CreateEvents(5, 30)
	// fmt.Printf("Created %v events", len(events))

	// events := app.UpdateEvents(5, 30)
	// fmt.Printf("Updated %v events", len(events))

	// events := app.DeleteEvents(5, 30)
	// fmt.Printf("Deleted %v events", len(events))

	// svc := outlook.NewOutlookService(ctx)
}

func initGoogleApp() (*google.App, error) {
	ctx := context.Background()
	calendarID := "CALENDAR_ID"
	expiry, _ := time.Parse(time.RFC3339, "2020-02-27T10:09:02-05:00")
	token := &oauth2.Token{
		AccessToken:  "",
		RefreshToken: "",
		Expiry:       expiry,
	}
	return google.New(ctx, token, calendarID)
}

func initOutlookApp(ownerID string) (*outlook.App, error) {
	ctx := context.Background()
	return outlook.New(ctx, ownerID)
}

func runOutlookTest(s *settings, app *outlook.App, ch chan string) {
	// events := app.ListEvents(10)
	// fmt.Printf("\tFetched %v events for user %s\n", len(events), s.userID)
	// for _, event := range events {
	// 	fmt.Println(event.Summary)
	// }

	var events []*calendar.Event
	batch := int(math.Floor(float64(s.numEvents) / 4))

	events = app.CreateEvents(batch, s.rate)
	fmt.Printf("Created %v events for user: %s\n", len(events), s.userID)

	events = app.UpdateEvents(batch, s.rate)
	fmt.Printf("Created %v events for user: %s\n", len(events), s.userID)

	events = app.DeleteEvents(batch, s.rate)
	fmt.Printf("Deleted %v events", len(events))

	events = app.CreateEvents(batch, s.rate)
	fmt.Printf("Created %v events\n", len(events))

	events = app.DeleteEvents(batch, s.rate)
	fmt.Printf("Deleted %v events", len(events))

	ch <- s.userID
}
