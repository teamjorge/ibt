package main

import (
	"fmt"
	"sort"

	"github.com/teamjorge/ibt"
	"github.com/teamjorge/ibt/headers"
	"github.com/teamjorge/ibt/utilities"
	"golang.org/x/exp/maps"
)

// TrackTempProcessors tracks the track temperature for each lap of the ibt file
type trackTempProcessor struct {
	tempMap map[int]float32
}

// NewTrackTempProcessor creates and initialises a new trackTempProcessor
func newTrackTempProcessor() *trackTempProcessor {
	t := new(trackTempProcessor)

	// tempMap will store a temperature value against a lap number
	t.tempMap = make(map[int]float32)

	return t
}

// Display name of the processor
func (t *trackTempProcessor) Name() string { return "Track Temp" }

// Method used for processing every tick of telemetry
func (t *trackTempProcessor) Process(input ibt.Tick, hasNext bool, session *headers.Session) error {
	trackTemp, err := ibt.GetTickValue[float32](input, "TrackTempCrew")
	if err != nil {
		return err
	}

	lap, err := ibt.GetTickValue[int](input, "Lap")
	if err != nil {
		return err
	}

	t.tempMap[lap] = trackTemp

	return nil
}

// Utility function for create a result that can be joined with other processors.
//
// This will convert the results to map[int]interface{}, where the keys will refer to laps.
// Result is not yet required by any interfaces, but is useful when using multiple processors
// that summarise telemetry based by lap.
func (t *trackTempProcessor) Result() map[int]interface{} {
	return utilities.CreateGenericMap(t.tempMap)
}

// Columns required for the processor
func (t *trackTempProcessor) Whitelist() []string { return []string{"Lap", "TrackTempCrew"} }

// Print the summarised Track Temperature
func (t *trackTempProcessor) Print() {
	fmt.Println("Track Temp:")
	laps := maps.Keys(t.tempMap)
	sort.Ints(laps)

	for _, lap := range laps {
		fmt.Printf("%03d - %.3f\n", lap, t.tempMap[lap])
	}
}
