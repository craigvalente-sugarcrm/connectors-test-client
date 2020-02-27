package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/connectors-test-client/app/google"
	"golang.org/x/oauth2"
)

func main() {

	// Helper func to request authCode and exchange it for token
	// google.GetTokenFromWeb()

	app, err := initGoogleApp()
	if err != nil {
		log.Fatal(err)
	}

	events := app.ListEvents(10)
	fmt.Printf("Fetched %v events", len(events))

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
