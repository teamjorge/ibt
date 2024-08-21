package main

import (
	"context"
	"log"

	"github.com/teamjorge/ibt"
	"github.com/teamjorge/ibt/examples"
)

func main() {
	// Parse the files into stubs
	stubs, err := examples.ParseExampleStubs()
	if err != nil {
		log.Fatal(err)
	}

	// Create our storage client
	storage := newStorage()
	if err := storage.Connect(); err != nil {
		log.Fatal(err)
	}
	// Close it when the application ends
	defer storage.Close()

	// We group our stubs mainly to be able to identify the batches we are loading
	// This might not be necessary on your use case
	groups := stubs.Group()

	for groupNumber, group := range groups {
		// Create a new processor for this group and set the groupNumber.
		// It embeds our storage and we set our loading threshold to 100
		processor := newLoaderProcessor(storage, groupNumber, 100)

		// Process the group
		if err := ibt.Process(context.Background(), group, processor); err != nil {
			log.Fatalf("failed to process telemetry for stubs %v: %v", stubs, err)
		}

		// Print the number of batches loaded after each group
		log.Printf("%d batches loaded after group %d\n", storage.Loaded(), groupNumber)
	}
}
