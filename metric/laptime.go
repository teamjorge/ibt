package metric

import (
	"fmt"
	"math"
)

// LapTime is a lap time measured in seconds
type LapTime float32

// ToString converts the given lap time to an 0m:0s.0ms form
//
// For example: 78,583 seconds will be representated as 01:18.583
func (l LapTime) ToString() string {
	lapMS := float64(l * 1000)

	minutes := math.Floor(lapMS / 60000.00)

	secondsRemaining := lapMS - (minutes * 60000)

	seconds := math.Floor(secondsRemaining / 1000.00)

	milliseconds := secondsRemaining - (seconds * 1000)

	return fmt.Sprintf("%02d:%02d.%03d", int(minutes), int(seconds), int(milliseconds))
}
