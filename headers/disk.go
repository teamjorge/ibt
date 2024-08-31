package headers

import (
	"fmt"
	"math"
	"time"

	"github.com/teamjorge/ibt/utilities"
)

const (
	DISK_HEADER_BYTES_SIZE int = 32
)

// DiskHeader is the ibt file header indicating start, end, and amount of records
type DiskHeader struct {
	// Unix timestamp indicating the start date and time of the file
	StartDate int64
	// Start time of file relative to the seconds since start of the session
	StartTime float64
	// End time of file relative to the seconds since start of the session
	EndTime float64
	// Number of laps telemetry exists for
	LapCount int
	// Number of telemetry variable records
	RecordCount int
}

// ReadDiskHeader attempts to parse the Disk SubHeader from the given Reader (a loaded .ibt file)
//
// This function assumes that the reader is already seeked to the current offset in the buffer. If the call
// was not preceded by headers.ReadTelemetryHeader(), you can seek to the current offset with:
//
//	_, err := reader.Seek(int64(TELEMETRY_HEADER_BYTES_SIZE), 0)
//
// Validation will be performed to ensure that the values are as expected
func ReadDiskHeader(reader Reader) (*DiskHeader, error) {
	diskHeaderBuf := make([]byte, DISK_HEADER_BYTES_SIZE)

	_, err := reader.ReadAt(diskHeaderBuf, int64(TELEMETRY_HEADER_BYTES_SIZE))
	if err != nil {
		return nil, fmt.Errorf("failed to read disk header buffer: %v", err)
	}

	h := DiskHeader{
		StartDate:   utilities.Byte8ToInt64(diskHeaderBuf[0:8]),
		StartTime:   utilities.Byte8ToFloat(diskHeaderBuf[8:16]),
		EndTime:     utilities.Byte8ToFloat(diskHeaderBuf[16:24]),
		LapCount:    utilities.Byte4ToInt(diskHeaderBuf[24:28]),
		RecordCount: utilities.Byte4ToInt(diskHeaderBuf[28:32]),
	}

	if h.EndTime < 0 || h.StartTime < 0 || h.EndTime > math.Pow(10, 20) || h.StartTime > math.Pow(10, 20) ||
		h.RecordCount == 0 {
		return nil, fmt.Errorf("invalid disk header detected. values received: %+v", h)
	}

	// Determine if StartDate is an invalid time value
	parsedTime := time.Unix(h.StartDate, 0)
	currentYear := time.Now().Year()

	if parsedTime.Year() < currentYear-20 || parsedTime.Year() > currentYear+20 {
		return nil, fmt.Errorf("invalid StartDate detected: %v", parsedTime)
	}

	return &h, nil
}
