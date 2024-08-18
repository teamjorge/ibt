package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/teamjorge/ibt"
	"github.com/teamjorge/ibt/headers"
	"github.com/teamjorge/ibt/utilities"
	"golang.org/x/exp/maps"
)

func main() {
	flag.Parse()

	// Use the default testing file if no ibt files were provided
	// Provided files can use wildcards, for example ./telemetry/*.ibt
	var filePattern string
	if flag.Arg(0) == "" {
		filePattern = ".testing/valid_test_file.ibt"
	} else {
		filePattern = flag.Arg(0)
	}

	// Find all files for parsing in case it had included a wildcard
	files, err := filepath.Glob(filePattern)
	if err != nil {
		log.Fatalf("could not glob the given input files: %v", err)
	}

	// Parse the files into stubs
	stubs, err := ibt.ParseStubs(files...)
	if err != nil {
		log.Fatalf("failed to parse stubs for %v. error - %v", files, err)
	}

	// Group the stubs into iRacing sessions
	stubGroups := stubs.Group()

	for _, stubGroup := range stubGroups {
		for _, stub := range stubGroup {
			stubFile, err := os.Open(stub.Filename())
			if err != nil {
				log.Fatalf("failed to open stub file %s for reading: %v", stub.Filename(), err)
			}

			// Create the instance(s) of your processor(s)
			processor := NewTrackTempProcessor()

			// Process the available telemetry for the ibt file. This is currently only utilising the Track Temp processor,
			// but can include as many as you want.
			if err := ibt.Process(stubFile, *stub.Headers(), processor); err != nil {
				log.Fatalf("failed to process telemetry for stub %s: %v", stub.Filename(), err)
			}

			// Print the summarised track temperature
			processor.Print()
		}
	}

}

type TrackTempProcessor struct {
	tempMap map[int]float32
}

func NewTrackTempProcessor() *TrackTempProcessor {
	t := new(TrackTempProcessor)

	t.tempMap = make(map[int]float32)

	return t
}

// Display name of the processor
func (t *TrackTempProcessor) Name() string { return "Track Temp" }

// Method used for processing every tick of telemetry
func (t *TrackTempProcessor) Process(input map[string]headers.VarHeader, hasNext bool, session *headers.Session) error {
	TrackTempProcessor := input["TrackTempCrew"].Value.(float32)
	lap := input["Lap"].Value.(int)

	t.tempMap[lap] = TrackTempProcessor

	return nil
}

// Utility function for create a result that can be joined with other processors.
//
// This will convert the results to map[int]interface{}, where the keys will refer to laps.
// Result is not yet required by any interfaces, but is useful when using multiple processors
// that summarise telemetry based by lap.
func (t *TrackTempProcessor) Result() map[int]interface{} {
	return utilities.CreateGenericMap(t.tempMap)
}

// Columns required for the processor
func (t *TrackTempProcessor) Whitelist() []string { return []string{"Lap", "TrackTempCrew"} }

// Print the summarised Track Temperature
func (t *TrackTempProcessor) Print() {
	fmt.Println("Track Temp:")
	laps := maps.Keys(t.tempMap)
	sort.Ints(laps)

	for _, lap := range laps {
		fmt.Printf("%03d - %.3f\n", lap, t.tempMap[lap])
	}
}
