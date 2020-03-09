package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/connectors-test-client/app/google"
	"github.com/connectors-test-client/app/outlook"
	"github.com/spf13/viper"
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

	viper.BindEnv("numMinutes", "DURATION")
	viper.BindEnv("eventsPerMinute", "RATE")

	numMinutes := viper.GetInt("numMinutes")
	eventsPerMinute := viper.GetInt("eventsPerMinute")

	log.Printf("Running for %v minutes at a rate of %v events / minute.", numMinutes, eventsPerMinute)

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
		AccessToken:  "ya29.a0Adw1xeUUFklLr-91F02hhnWIkdQtTvFhizN7rCs9SJx4laVLZMY5EtzpvqQVRQeky4FPNhEdGcTs50k7HJ5y9ZQ78_p8U4H-Dabduu30msP8SEXNvrTWAJdOY5QVUz4aSW8Ul7_tM950ZvWAcb0gLdFoKj1lOgFQhYQ",
		RefreshToken: "1//01yraftveQw5oCgYIARAAGAESNwF-L9IrUjOwiPKsfat2sExOBJFT8k4pku3r-tckkZ9eomcmNSdxbEZce9CKBIN4e5EA5EXvg3g",
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
	batch := int(math.Ceil(float64(s.numEvents) / 4))

	events = app.CreateEvents(batch, s.rate)
	fmt.Printf("Created %v events for user: %s\n", len(events), s.userID)

	syncCh := make(chan int, 2)
	events = app.ListEvents(batch)

	go func(eventsToUpdate []*calendar.Event) {
		app.UpdateTheseEvents(eventsToUpdate, s.rate)
		fmt.Printf("Updated %v events for user: %s\n", len(eventsToUpdate), s.userID)
		syncCh <- 1
	}(events)
	go func() {
		events = app.CreateEvents(batch, s.rate)
		fmt.Printf("Created %v events for user: %s\n", len(events), s.userID)
		syncCh <- 1
	}()
	<-syncCh
	<-syncCh

	events = app.ListEvents(batch * 2)
	go func(eventsToDelete []*calendar.Event) {
		app.DeleteTheseEvents(eventsToDelete, s.rate)
		fmt.Printf("Deleted %v events for user: %s\n", len(eventsToDelete), s.userID)
		syncCh <- 1
	}(events)
	go func() {
		events = app.CreateEvents(batch, s.rate)
		fmt.Printf("Created %v events for user: %s\n", len(events), s.userID)
		syncCh <- 1
	}()
	<-syncCh
	<-syncCh

	events = app.DeleteEvents(batch, s.rate)
	fmt.Printf("Deleted %v events for user: %s\n", len(events), s.userID)

	ch <- s.userID
}
