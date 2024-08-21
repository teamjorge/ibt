package ibt

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/teamjorge/ibt/headers"
)

// Stub represents the headers and filename parsed from an ibt file.
//
// Stubs are used for initial parsing of ibt files and their metadata. This can be useful
// to make decisions regarding which files should have their telemetry parsed.
type Stub struct {
	filepath string
	header   headers.Header
}

// Open the underlying ibt file for reading
func (stub *Stub) Open() (headers.Reader, error) {
	reader, err := os.Open(stub.Filename())
	if err != nil {
		return nil, fmt.Errorf("failed to open stub file %s for reading: %v", stub.Filename(), err)
	}

	return reader, nil
}

// Filename where the stub originated from
func (stub *Stub) Filename() string { return stub.filepath }

// Headers that were parsed when the stub was created
func (stub *Stub) Headers() *headers.Header { return &stub.header }

// Time when the stub was created
func (stub *Stub) Time() time.Time {
	parsedTime := time.Unix(stub.header.DiskHeader().StartDate, 0)

	return parsedTime
}

// DriverIdx is the index of the current driver.
//
// This value is useful when parsing telemetry or session info where the index of the driver
// is required.
func (stub *Stub) DriverIdx() int {
	return stub.header.SessionInfo().DriverInfo.DriverCarIdx
}

// StubGroup is a grouping of stubs.
//
// This group is not necessarily part of the same session, but can be grouped with Group().
type StubGroup []Stub

// ParseStubs will create a stub for each of the given files by parsing their headers.
func ParseStubs(files ...string) (StubGroup, error) {
	stubs := make(StubGroup, 0)

	for _, file := range files {
		stub, err := parseStub(file)
		if err != nil {
			return stubs, err
		}

		stubs = append(stubs, stub)
	}

	return stubs, nil
}

// parseStub will create a stub from the given file by parsing it's headers.
func parseStub(filename string) (Stub, error) {
	var stub Stub

	f, err := os.Open(filename)
	if err != nil {
		return stub, fmt.Errorf("failed to open file %s for reading: %v", filename, err)
	}
	defer f.Close()

	header, err := headers.ParseHeaders(f)
	if err != nil {
		return stub, fmt.Errorf("failed to parse headers for file %s - %v", filename, err)
	}

	return Stub{filename, header}, nil
}

// Group stubs together by their iRacing session.
//
// The process for grouping is slightly different for official and test sessions.
// Official sessions can utilise the SubSessionID, whereas Test sessions
// use the ResultsPosition field to determine the start/end of a session.
func (stubs StubGroup) Group() []StubGroup {
	sessionStubMap := make(map[int]StubGroup)

	// Group stubs of the same SubSessionID together
	for _, stub := range stubs {
		subSessionId := stub.header.SessionInfo().WeekendInfo.SubSessionID
		sessionStubMap[subSessionId] = append(sessionStubMap[subSessionId], stub)
	}

	groups := make(StubGroupGrouping, 0)
	// Groups that have a SubSessionID of 0 are considered as Test sessions and
	// are grouped using a separate method
	if testSessionStubGroup, ok := sessionStubMap[0]; ok {
		groups = append(groups, groupTestSessionStubs(testSessionStubGroup)...)
	}

	delete(sessionStubMap, 0)

	for _, stubGroup := range sessionStubMap {
		sort.Sort(stubGroup)
		groups = append(groups, stubGroup)
	}

	sort.Sort(groups)

	return groups
}

// groupTestSessionStubs ensures that ibt files from iRacing Test sessions are grouped correctly.
//
// The logic for grouping Test session files is slightly different due to the lack of subSessionIds
// and rely on the ResultsPosition variable to determine start and end of a group.
func groupTestSessionStubs(stubs StubGroup) []StubGroup {
	groups := make([]StubGroup, 0)

	currentGroup := make(StubGroup, 0)
	for _, stub := range stubs {
		// ResultsPosition nil indicate the first ibt file of a new session
		if stub.header.SessionInfo().SessionInfo.Sessions[0].ResultsPositions != nil {
			currentGroup = append(currentGroup, stub)
		} else {
			// Determine if it should end the existing group and create a new one
			if len(currentGroup) > 0 {
				sort.Sort(currentGroup)
				groups = append(groups, currentGroup)
			}
			currentGroup = StubGroup{stub}
		}
	}
	// Ensure group contains stubs
	if len(currentGroup) > 0 {
		sort.Sort(currentGroup)
		groups = append(groups, currentGroup)
	}

	return groups
}

func (a StubGroup) Len() int           { return len(a) }
func (a StubGroup) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a StubGroup) Less(i, j int) bool { return a[i].Time().Before(a[j].Time()) }

// StubGroupGrouping is just a group of StubGroups.
type StubGroupGrouping []StubGroup

func (a StubGroupGrouping) Len() int      { return len(a) }
func (a StubGroupGrouping) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a StubGroupGrouping) Less(i, j int) bool {
	return len(a[i]) > 0 && len(a[j]) > 0 && a[i][0].Time().Before(a[j][0].Time())
}
