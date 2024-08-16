package ibt

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/teamjorge/ibt/headers"
)

type IbtStub struct {
	filepath string
	header   headers.Header
}

func (stub *IbtStub) Filename() string         { return stub.filepath }
func (stub *IbtStub) Headers() *headers.Header { return &stub.header }

func (stub *IbtStub) Time() (time.Time, error) {
	datePattern, err := regexp.Compile(`\d{4}-\d{2}-\d{2}\s\d{2}-\d{2}-\d{2}`)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse date pattern for stub time - %v", err)
	}

	filename := filepath.Base(stub.filepath)

	foundPattern := datePattern.FindString(filename)
	if foundPattern == "" {
		return time.Time{}, fmt.Errorf("failed to find date pattern in filename %s", stub.filepath)
	}

	t, err := time.Parse("2006-01-02 15-04-05", foundPattern)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse date pattern for filename %s - %v", stub.filepath, err)
	}

	return t, nil
}

func (stub *IbtStub) DriverIdx() int {
	return stub.header.SessionInfo().DriverInfo.DriverCarIdx
}

type StubGroup []IbtStub

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

func parseStub(filename string) (IbtStub, error) {
	var stub IbtStub

	f, err := os.Open(filename)
	if err != nil {
		return stub, fmt.Errorf("failed to open file %s for reading: %v", filename, err)
	}
	defer f.Close()

	header, err := headers.ParseHeaders(f)
	if err != nil {
		return stub, fmt.Errorf("failed to parse headers for file %s - %v", filename, err)
	}

	return IbtStub{filename, header}, nil
}

func (stubs StubGroup) Group() []StubGroup {
	sessionStubMap := make(map[int]StubGroup)

	for _, stub := range stubs {
		subSessionId := stub.header.SessionInfo().WeekendInfo.SubSessionID
		sessionStubMap[subSessionId] = append(sessionStubMap[subSessionId], stub)
	}

	groups := make([]StubGroup, 0)
	if testSessionStubGroup, ok := sessionStubMap[0]; ok {
		groups = append(groups, groupTestSessionStubs(testSessionStubGroup)...)
	}

	delete(sessionStubMap, 0)

	for _, stub := range sessionStubMap {
		groups = append(groups, stub)
	}

	return groups
}

func groupTestSessionStubs(stubs StubGroup) []StubGroup {
	groups := make([]StubGroup, 0)

	currentGroup := make([]IbtStub, 0)
	for _, stub := range stubs {
		if stub.header.SessionInfo().SessionInfo.Sessions[0].ResultsPositions != nil {
			currentGroup = append(currentGroup, stub)
		} else {
			if len(currentGroup) > 0 {
				groups = append(groups, currentGroup)
			}
			currentGroup = StubGroup{stub}
		}
	}
	if len(currentGroup) > 0 {
		groups = append(groups, currentGroup)
	}

	return groups
}
