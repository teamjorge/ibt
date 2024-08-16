package headers

import "fmt"

// Header contains all sub-headers present in the ibt file.
type Header struct {
	telemHeader *TelemetryHeader
	diskHeader  *DiskHeader
	varHeader   map[string]VarHeader
	sessionInfo *Session
}

// TelemetryHeader returns the loaded telemetry sub-header.
func (h *Header) TelemetryHeader() *TelemetryHeader { return h.telemHeader }

// DiskHeader returns the loaded disk sub-header.
func (h *Header) DiskHeader() *DiskHeader { return h.diskHeader }

// VarHeader returns the loaded variable sub-header.
func (h *Header) VarHeader() map[string]VarHeader { return h.varHeader }

// SessionInfo returns the parsed session information.
func (h *Header) SessionInfo() *Session { return h.sessionInfo }

// ParseHeader parses each of the required sub-headers of the ibt file in sequence.
func ParseHeaders(r Reader) (Header, error) {
	var header Header

	telemHeader, err := ReadTelemetryHeader(r)
	if err != nil {
		return header, fmt.Errorf("failed to parse telemetry header: %v", err)
	}

	diskHeader, err := ReadDiskHeader(r)
	if err != nil {
		return header, fmt.Errorf("failed to parse disk header: %v", err)
	}

	sessionInfo, err := ReadSessionInfo(r, telemHeader.SessionInfoOffset, telemHeader.SessionInfoLength)
	if err != nil {
		return header, fmt.Errorf("failed to parse session info: %v", err)
	}

	varHeader, err := ReadVarHeader(r, telemHeader.NumVars, telemHeader.VarHeaderOffset)
	if err != nil {
		return header, fmt.Errorf("failed to parse variable header: %v", err)
	}

	return Header{
		telemHeader: telemHeader,
		diskHeader:  diskHeader,
		sessionInfo: sessionInfo,
		varHeader:   varHeader,
	}, nil
}
