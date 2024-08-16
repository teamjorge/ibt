package metric

import (
	"fmt"
	"math"
)

type LapTime float32

func (l LapTime) ToString() string {
	lapMS := float64(l * 1000)

	minutes := math.Floor(lapMS / 60000.00)

	secondsRemaining := lapMS - (minutes * 60000)

	seconds := math.Floor(secondsRemaining / 1000.00)

	milliseconds := secondsRemaining - (seconds * 1000)

	return fmt.Sprintf("%02d:%02d.%03d", int(minutes), int(seconds), int(milliseconds))
}
