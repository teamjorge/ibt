package main

import (
	"fmt"

	"github.com/teamjorge/ibt"
	"github.com/teamjorge/ibt/headers"
)

type loaderProcessor struct {
	// Our storage client
	*storage
	// Cache for holding the number of telemetry ticks equal to threshold
	cache []map[string]interface{}
	// Number
	groupNumber int
	threshold   int
}

// Simple Constructor for creating our processor
func newLoaderProcessor(storage *storage, groupNumber int, threshold int) *loaderProcessor {
	return &loaderProcessor{storage, make([]map[string]interface{}, 0), groupNumber, threshold}
}

// Columns we want to parse from telemetry
func (l *loaderProcessor) Whitelist() []string {
	return []string{
		"Lap", "ThrottleRaw", "BrakeRaw", "ClutchRaw", "LapDistPct", "Lat", "Lon",
	}
}

// Our method for processing a single tick of telemetry.
func (l *loaderProcessor) Process(input ibt.Tick, hasNext bool, session *headers.Session) error {
	// Add our group number to the tick of telemetry.
	// This will be useful to seperate ticks by group in our storage.
	input["groupNum"] = l.groupNumber

	// Add it to the cache
	l.cache = append(l.cache, input)

	// If our cache is past the threshold, that means we can now do a bulk load
	// to our storage.
	if len(l.cache) >= l.threshold {
		if err := l.loadBatch(); err != nil {
			return fmt.Errorf("failed to load batch - %v", err)
		}
		// Empty the cache again
		l.cache = make([]map[string]interface{}, 0)
	}

	return nil
}

func (l *loaderProcessor) loadBatch() error {
	// Bulk load our batch to storage.
	return l.Exec(l.cache)
}
