package ibt

import (
	"github.com/teamjorge/ibt/headers"
)

// Parser is used to iterate and process telemetry variables for a given ibt file and it's headers.
type Parser struct {
	// File or Live Telemetry reader
	reader headers.Reader
	// List of columns to parse
	whitelist []string
	header    *headers.Header

	current int
}

// NewParser creates a new parser from a given ibt file, it's headers, and a variable whitelist.
//
// reader - Opened ibt file.
//
// header - Parsed headers of ibt file.
//
// whitelist - Variables to process. For example, "gear", "speed", "rpm" etc. If no values or a
// single value of "*" is received, all variables will be processed.
func NewParser(reader headers.Reader, header *headers.Header, whitelist ...string) *Parser {
	p := new(Parser)

	p.reader = reader
	p.whitelist = whitelist
	p.header = header

	p.current = 1

	return p
}

// Next parses and returns the next tick of telemetry variables and whether it can be called again.
//
// A return of false will indicate that the buffer has reached the end. If the buffer has reached the end and Next() is called again,
// a nil and false will be returned.
//
// Should expected variable values be missing, please ensure that they are added to the Parser whitelist.
func (p *Parser) Next() (Tick, bool) {
	start := p.header.TelemetryHeader.BufOffset + (p.current * p.header.TelemetryHeader.BufLen)

	currentBuf := p.read(start)
	if currentBuf == nil {
		return nil, false
	}

	// Read in the next buffer to determine if more telemetry ticks are available.
	nextStart := p.header.TelemetryHeader.BufOffset + ((p.current + 1) * p.header.TelemetryHeader.BufLen)
	nextBuf := p.read(nextStart)

	newVars := p.readVarsFromBuffer(currentBuf)

	p.current++

	return newVars, nextBuf != nil
}

// ParseAt the given buffer offset and return a processed tick.
//
// ParseAt is useful if a specific offset is known. An example would be the
// telemetry variable buffers that are provided during live telemetry parsing.
func (p *Parser) ParseAt(offset int) Tick {
	currentBuf := p.read(offset)
	if currentBuf == nil {
		return nil
	}

	newVars := p.readVarsFromBuffer(currentBuf)

	return newVars
}

// read the next buffer from offset to the current length set by the parser.
func (p *Parser) read(start int) []byte {
	buf := make([]byte, p.header.TelemetryHeader.BufLen)
	_, err := p.reader.ReadAt(buf, int64(start))
	if err != nil {
		defer p.reader.Close()
		return nil
	}

	return buf
}

// readVarsFromBuffer reads each of the specified (whitelist) fields from the given buffer into a new Tick.
func (p *Parser) readVarsFromBuffer(buf []byte) Tick {
	newVars := make(Tick)

	for _, variable := range p.whitelist {
		item := p.header.VarHeader[variable]
		val := readVarValue(buf, item)
		newVars[variable] = val
	}

	return newVars
}

// Seek the parser to a specific tick within the ibt file.
func (p *Parser) Seek(iter int) { p.current = iter }

// UpdateWhitelist replaces the current whitelist with the given fields
func (p *Parser) UpdateWhitelist(whitelist ...string) {
	p.whitelist = whitelist
}
