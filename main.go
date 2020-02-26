package main

import (
	"context"
	"fmt"
	"log"

	"github.com/connectors-test-client/app/google"
)

func main() {

	ctx := context.Background()
	// client := google.NewGoogleClient()
	svc, err := google.NewGoogleService(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// _, _ = google.ListEvents(svc, "craig.valente@sugarcrm.com", 10)

	events := google.CreateEvents(svc, "craig.valente@sugarcrm.com", 5, 30)
	fmt.Printf("Created %v events", len(events))

	// events := google.UpdateEvents(svc, "craig.valente@sugarcrm.com", 5, 30)
	// fmt.Printf("Updated %v events", len(events))

	// events := google.DeleteTestEvents(svc, "craig.valente@sugarcrm.com", 5, 30)
	// fmt.Printf("Updated %v events", len(events))

	// svc := outlook.NewOutlookService(ctx)
}
