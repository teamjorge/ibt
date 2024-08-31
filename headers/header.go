package headers

import "fmt"

// Header contains all sub-headers present in the ibt file.
type Header struct {
	TelemetryHeader *TelemetryHeader
	DiskHeader      *DiskHeader
	VarHeader       map[string]VarHeader
	SessionInfo     *Session
	VarBuffers      []VarBuffer
}

// ParseHeader parses each of the required sub-headers of the ibt file in sequence.
func ParseHeaders(r Reader) (*Header, error) {
	telemHeader, err := ReadTelemetryHeader(r)
	if err != nil {
		return nil, fmt.Errorf("failed to parse telemetry header: %v", err)
	}

	diskHeader, err := ReadDiskHeader(r)
	if err != nil {
		return nil, fmt.Errorf("failed to parse disk header: %v", err)
	}

	varHeader, err := ReadVarHeader(r, telemHeader.NumVars, telemHeader.VarHeaderOffset)
	if err != nil {
		return nil, fmt.Errorf("failed to parse variable header: %v", err)
	}

	varBuffers, err := ReadVarBufferHeaders(r, telemHeader.NumBuf)
	if err != nil {
		return nil, fmt.Errorf("failed to parse var buffer header: %v", err)
	}

	sessionInfo, err := ReadSessionInfo(r, telemHeader.SessionInfoOffset, telemHeader.SessionInfoLength)
	if err != nil {
		return nil, fmt.Errorf("failed to parse session info: %v", err)
	}

	return &Header{
		TelemetryHeader: telemHeader,
		DiskHeader:      diskHeader,
		VarHeader:       varHeader,
		SessionInfo:     sessionInfo,
		VarBuffers:      varBuffers,
	}, nil
}

func (h *Header) UpdateVarBuffer(r Reader) error {
	varBuffers, err := ReadVarBufferHeaders(r, h.TelemetryHeader.NumBuf)
	if err != nil {
		return fmt.Errorf("failed to parse var buffer header: %v", err)
	}

	h.VarBuffers = varBuffers

	return nil
}
