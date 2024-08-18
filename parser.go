package ibt

import (
	"github.com/teamjorge/ibt/headers"
	"golang.org/x/exp/maps"
)

// Parser is used to iterate and process telemetry variables for a given ibt file and it's headers.
type Parser struct {
	reader       headers.Reader
	varHeader    map[string]headers.VarHeader
	whitelist    []string
	length       int
	bufferOffset int
	current      int
}

// NewParser creates a new parser from a given ibt file, it's headers, and a variable whitelist.
//
// reader - Opened ibt file.
//
// header - Parsed headers of ibt file.
//
// whitelist - Variables to process. For example, "gear", "speed", "rpm" etc. If no values or a
// single value of "*" is received, all variables will be processed.
func NewParser(reader headers.Reader, header headers.Header, whitelist ...string) *Parser {
	p := new(Parser)

	p.reader = reader
	p.whitelist = computeVars(header.VarHeader(), whitelist...)

	p.length = header.TelemetryHeader().BufLen
	p.bufferOffset = header.TelemetryHeader().BufOffset
	p.varHeader = header.VarHeader()
	p.current = 1

	return p
}

// Next parses and returns the next tick of telemetry variables and whether it can be called again.
//
// A return of false will indicate that the buffer has reached the end. If the buffer has reached the end and Next() is called again,
// a nil and false will be returned.
//
// Should expected variable values be missing, please ensure that they are added to the Parser whitelist.
func (p *Parser) Next() (map[string]headers.VarHeader, bool) {
	start := p.bufferOffset + (p.current * p.length)
	currentBuf := p.read(start)
	if currentBuf == nil {
		return nil, false
	}

	// Read in the next buffer to determine if more telemetry ticks are available.
	nextStart := p.bufferOffset + ((p.current + 1) * p.length)
	nextBuf := p.read(nextStart)

	newVars := make(map[string]headers.VarHeader)

	for _, variable := range p.whitelist {
		item := p.varHeader[variable]
		val := readVarValue(currentBuf, item)
		item.Value = val
		newVars[variable] = item
	}

	p.current++

	return newVars, nextBuf != nil
}

// read the next buffer from offset to the current length set by the parser
func (p *Parser) read(start int) []byte {
	buf := make([]byte, p.length)
	_, err := p.reader.ReadAt(buf, int64(start))
	if err != nil {
		defer p.reader.Close()
		return nil
	}

	return buf
}

// compareVars will retrieve vars when * is used and ensure a unique list
//
// Variables that are not found in the VarHeader will automatically be excluded.
func computeVars(vars map[string]headers.VarHeader, whitelist ...string) []string {
	if len(whitelist) == 0 {
		return headers.AvailableVars(vars)
	}

	for _, col := range whitelist {
		if col == "*" {
			return headers.AvailableVars(vars)
		}
	}

	// de-duplicate the columns
	varMap := make(map[string]struct{})

	for _, col := range whitelist {
		// ensure it's a valid column
		if _, ok := vars[col]; ok {
			varMap[col] = struct{}{}
		}
	}

	return maps.Keys(varMap)
}
