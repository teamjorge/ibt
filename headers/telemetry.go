package headers

import (
	"fmt"

	"github.com/teamjorge/ibt/utilities"
)

const (
	TELEMETRY_HEADER_BYTES_SIZE int = 112
)

type TelemetryHeader struct {
	// Version of ibt file
	Version int
	// Status indicates whether session is live (0) or completed (1)
	Status int
	// Tickrate indicates the frequency of telemetry data written to the file.
	// A value of 60 indicates 60 times per second
	TickRate int
	// Indicates the number of times SessionInfo was updated for the current file.
	// The value will be 0 for completed sessions and >1 for active sessions
	SessionInfoUpdate int
	// Buffer offset for SessionInfo data
	SessionInfoOffset int
	// Length of the SessionInfo buffer
	SessionInfoLength int
	// Number of available telemetry vars that will be written at the Tickrate frequency
	NumVars int
	// Buffer offset for VarHeader
	VarHeaderOffset int
	// Specifies the number of telemetry data buffers available.
	// This will be 1 for ibt files and 3 for memory-mapped live telemetry
	NumBuf int
	// Length of the buffer for parsing VarHeader telemetry values
	BufLen int
	// Buffer offset for the VarHeader telemetry values
	BufOffset int
}

// ReadTelemetryHeader attempts to parse the TelemetryHeader from the given Reader (a loaded .ibt file)
//
// Validation will be performed to ensure that the values are as expected
func ReadTelemetryHeader(reader Reader) (*TelemetryHeader, error) {
	headerBuf := make([]byte, TELEMETRY_HEADER_BYTES_SIZE)

	_, err := reader.ReadAt(headerBuf, 0)
	if err != nil {
		return nil, err
	}

	h := TelemetryHeader{
		Version:           utilities.Byte4ToInt(headerBuf[0:4]),
		Status:            utilities.Byte4ToInt(headerBuf[4:8]),
		TickRate:          utilities.Byte4ToInt(headerBuf[8:12]),
		SessionInfoUpdate: utilities.Byte4ToInt(headerBuf[12:16]),
		SessionInfoLength: utilities.Byte4ToInt(headerBuf[16:20]),
		SessionInfoOffset: utilities.Byte4ToInt(headerBuf[20:24]),
		NumVars:           utilities.Byte4ToInt(headerBuf[24:28]),
		VarHeaderOffset:   utilities.Byte4ToInt(headerBuf[28:32]),
		NumBuf:            utilities.Byte4ToInt(headerBuf[32:36]),
		BufLen:            utilities.Byte4ToInt(headerBuf[36:40]),
		BufOffset:         utilities.Byte4ToInt(headerBuf[52:56]),
	}

	if h.Version != 2 || h.Status < 0 || h.Status > 1 || h.TickRate < 60 || h.TickRate > 360 {
		return nil, fmt.Errorf("invalid telemetry header detected. values received: %v", h)
	}

	return &h, nil
}
