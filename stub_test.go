package ibt

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/teamjorge/ibt/headers"
)

func TestStubs(t *testing.T) {
	f, err := os.Open(".testing/valid_test_file.ibt")
	if err != nil {
		t.Errorf("failed to open testing file - %v", err)
		return
	}
	defer f.Close()

	header, err := headers.ParseHeaders(f)
	if err != nil {
		t.Errorf("expected headers to parse correctly: %v", err)
	}

	testStub := Stub{
		filepath: ".testing/valid_test_file.ibt",
		header:   header,
	}

	t.Run("stubs Filename()", func(t *testing.T) {
		filename := testStub.Filename()

		if filename != ".testing/valid_test_file.ibt" {
			t.Errorf("expected filename to be .testing/valid_test_file.ibt but got %s", filename)
		}
	})

	t.Run("stubs Headers()", func(t *testing.T) {
		h := testStub.Headers()

		if reflect.DeepEqual(h, header) {
			t.Errorf("expected %v = %v", h, header)
		}
	})

	t.Run("stubs Time() valid time", func(t *testing.T) {
		parsedTime, err := testStub.Time()

		if err != nil {
			t.Errorf("expected err to be nil but got %v", err)
		}

		expectedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-06-24 16:02:03")
		if parsedTime.Equal(expectedTime) {
			t.Errorf("expected time to be %v but got %v", expectedTime, parsedTime)
		}
	})

	t.Run("stubs Time() invalid time", func(t *testing.T) {
		// Value below threshold
		h := headers.NewHeader(&headers.TelemetryHeader{}, &headers.DiskHeader{StartDate: 192138}, map[string]headers.VarHeader{}, &headers.Session{})

		invalidTimeStub := Stub{filepath: ".testing/valid_test_file.ibt", header: h}
		_, err := invalidTimeStub.Time()

		if err == nil {
			t.Error("expected err to contain a value but received nil")
		}

		invalidTimeStub.header = headers.NewHeader(
			&headers.TelemetryHeader{},
			&headers.DiskHeader{StartDate: time.Now().Add(time.Hour * (24 * 1000)).Unix()},
			map[string]headers.VarHeader{},
			&headers.Session{},
		)

		invalidTimeStub = Stub{filepath: ".testing/valid_test_file.ibt", header: h}
		_, err = invalidTimeStub.Time()

		if err == nil {
			t.Error("expected err to contain a value but received nil")
		}
	})

	t.Run("stubs DriverIdx()", func(t *testing.T) {
		h := headers.NewHeader(
			&headers.TelemetryHeader{},
			&headers.DiskHeader{},
			map[string]headers.VarHeader{},
			&headers.Session{DriverInfo: headers.DriverInfo{DriverCarIdx: 15}})

		driverIdxStub := Stub{filepath: ".testing/valid_test_file.ibt", header: h}

		if driverIdxStub.DriverIdx() != 15 {
			t.Errorf("expected driver idx to be 15, but got %d", driverIdxStub.DriverIdx())
		}
	})
}

func TestParseStubs(t *testing.T) {
	t.Run("test parseStub valid file", func(t *testing.T) {
		parsedStub, err := parseStub(".testing/valid_test_file.ibt")
		if err != nil {
			t.Errorf("unexpected error received from parseStub(): %v", err)
		}

		if parsedStub.filepath != ".testing/valid_test_file.ibt" {
			t.Errorf("expected parsedStub to have a filepath of .testing/valid_test_file.ibt. received: %s", parsedStub.filepath)
		}

		if parsedStub.header.TelemetryHeader().BufOffset != 53764 {
			t.Errorf("expected parsedStub to have a BufOffset of %d. received: %d", 53764, parsedStub.header.TelemetryHeader().BufOffset)
		}

		if parsedStub.header.DiskHeader().StartDate != 1719258336 {
			t.Errorf("expected parsedStub to have a StartDate of %d. received: %d", 1719258336, parsedStub.header.DiskHeader().StartDate)
		}

		if parsedStub.header.SessionInfo().DriverInfo.DriverCarVersion != "2024.05.28.02" {
			t.Errorf("expected parsedStub to have a DriverCarVersion of %s. received: %s", "2024.05.28.02", parsedStub.header.SessionInfo().DriverInfo.DriverCarVersion)
		}

		if len(parsedStub.header.VarHeader()) != 276 {
			t.Errorf("expected parsedStub to have a %d VarHeaders. received: %d", 276, len(parsedStub.header.VarHeader()))
		}
	})

	t.Run("test parseStub invalid file", func(t *testing.T) {
		_, err := parseStub(".testing/invalid_test_file.ibt")
		if err == nil {
			t.Error("expected an error from parseStub() when reading an invalid file")
		}
	})

	t.Run("test parseStub non-existent file", func(t *testing.T) {
		_, err := parseStub(".testing/dissapear_here.ibt")
		if err == nil {
			t.Error("expected an error from parseStub() when reading an invalid file")
		}
	})

	t.Run("test ParseStubs() valid file", func(t *testing.T) {
		parsedStubs, err := ParseStubs(".testing/valid_test_file.ibt", ".testing/valid_test_file.ibt")
		if err != nil {
			t.Errorf("unexpected error received from ParseStubs(): %v", err)
		}

		if len(parsedStubs) != 2 {
			t.Errorf("expected %d stubs to be parsed. found %d. parsed stubs: %v", 2, len(parsedStubs), parsedStubs)
		}

		if parsedStubs[1].filepath != ".testing/valid_test_file.ibt" {
			t.Errorf("expected the parsed stub to have a filepath of .testing/valid_test_file.ibt. received: %s", parsedStubs[1].filepath)
		}

		if parsedStubs[0].header.TelemetryHeader().BufOffset != 53764 {
			t.Errorf("expected the parsed stub to have a BufOffset of %d. received: %d", 53764, parsedStubs[0].header.TelemetryHeader().BufOffset)
		}

		if parsedStubs[1].header.DiskHeader().StartDate != 1719258336 {
			t.Errorf("expected the parsed stub to have a StartDate of %d. received: %d", 1719258336, parsedStubs[1].header.DiskHeader().StartDate)
		}

		if parsedStubs[0].header.SessionInfo().DriverInfo.DriverCarVersion != "2024.05.28.02" {
			t.Errorf("expected the parsed stub to have a DriverCarVersion of %s. received: %s", "2024.05.28.02", parsedStubs[0].header.SessionInfo().DriverInfo.DriverCarVersion)
		}

		if len(parsedStubs[1].header.VarHeader()) != 276 {
			t.Errorf("expected the parsed stub to have a %d VarHeaders. received: %d", 276, len(parsedStubs[1].header.VarHeader()))
		}
	})

	t.Run("test ParseStubs() one invalid file", func(t *testing.T) {
		parsedStubs, err := ParseStubs(".testing/valid_test_file.ibt", ".testing/invalid_test_file.ibt")
		if err == nil {
			t.Error("expected an error from ParseStubs() when reading an invalid file")
		}

		if len(parsedStubs) != 1 {
			t.Errorf("expected %d stubs to be parsed. found %d. parsed stubs: %v", 1, len(parsedStubs), parsedStubs)
		}

		if parsedStubs[0].header.TelemetryHeader().BufOffset != 53764 {
			t.Errorf("expected the parsed stub to have a BufOffset of %d. received: %d", 53764, parsedStubs[0].header.TelemetryHeader().BufOffset)
		}
	})
}

func TestGroupTestSessionStubs(t *testing.T) {
	sessionResultPos := &headers.Session{SessionInfo: headers.SessionInfo{Sessions: []headers.Sessions{{ResultsPositions: 1}}}}
	sessionNoResultPos := &headers.Session{SessionInfo: headers.SessionInfo{Sessions: []headers.Sessions{{ResultsPositions: nil}}}}

	stub1 := Stub{
		filepath: "stub_1.ibt",
		header:   headers.NewHeader(nil, nil, nil, sessionNoResultPos),
	}

	stub2 := Stub{
		filepath: "stub_2.ibt",
		header:   headers.NewHeader(nil, nil, nil, sessionResultPos),
	}

	stub3 := Stub{
		filepath: "stub_3.ibt",
		header:   headers.NewHeader(nil, nil, nil, sessionResultPos),
	}

	t.Run("test groupTestSessionStubs() with regular pattern", func(t *testing.T) {
		stubs := StubGroup{stub1, stub2, stub3, stub1, stub2, stub3}

		grouped := groupTestSessionStubs(stubs)

		if len(grouped) != 2 {
			t.Errorf("expected length of stub group to be %d. received %d. group: %v", 2, len(grouped), grouped)
		}

		if grouped[0][1].filepath != "stub_2.ibt" {
			t.Errorf("expected second item of the first group to have filename %s. received %s", "stub_2.ibt", grouped[0][1].filepath)
		}

		if grouped[1][2].filepath != "stub_3.ibt" {
			t.Errorf("expected third item of the second group to have filename %s. received %s", "stub_3.ibt", grouped[1][2].filepath)
		}
	})

	t.Run("test groupTestSessionStubs() with irregular pattern", func(t *testing.T) {
		stubs := StubGroup{stub2, stub3, stub1, stub2, stub3, stub1}

		grouped := groupTestSessionStubs(stubs)

		if len(grouped) != 3 {
			t.Errorf("expected length of stub group to be %d. received %d. group: %v", 2, len(grouped), grouped)
		}

		if grouped[0][0].filepath != "stub_2.ibt" {
			t.Errorf("expected first item of the first group to have filename %s. received %s", "stub_2.ibt", grouped[0][0].filepath)
		}

		if grouped[1][1].filepath != "stub_2.ibt" {
			t.Errorf("expected second item of the second group to have filename %s. received %s", "stub_2.ibt", grouped[1][1].filepath)
		}

		if grouped[2][0].filepath != "stub_1.ibt" {
			t.Errorf("expected first item of the third group to have filename %s. received %s", "stub_1.ibt", grouped[2][0].filepath)
		}
	})
}

func TestGroup(t *testing.T) {
	sessionWithSubsession1 := &headers.Session{WeekendInfo: headers.WeekendInfo{SubSessionID: 1}}
	sessionWithSubsession2 := &headers.Session{WeekendInfo: headers.WeekendInfo{SubSessionID: 2}}
	sessionWithSubsession3 := &headers.Session{WeekendInfo: headers.WeekendInfo{SubSessionID: 3}}

	sessionNoResultPos := &headers.Session{SessionInfo: headers.SessionInfo{Sessions: []headers.Sessions{{ResultsPositions: nil}}}}
	sessionResultPos := &headers.Session{SessionInfo: headers.SessionInfo{Sessions: []headers.Sessions{{ResultsPositions: 1}}}}

	stub1 := Stub{
		filepath: "stub_1.ibt",
		header:   headers.NewHeader(nil, nil, nil, sessionWithSubsession1),
	}

	stub2 := Stub{
		filepath: "stub_2.ibt",
		header:   headers.NewHeader(nil, nil, nil, sessionWithSubsession1),
	}

	stub3 := Stub{
		filepath: "stub_3.ibt",
		header:   headers.NewHeader(nil, nil, nil, sessionWithSubsession1),
	}

	stub4 := Stub{
		filepath: "stub_4.ibt",
		header:   headers.NewHeader(nil, nil, nil, sessionWithSubsession2),
	}

	stub5 := Stub{
		filepath: "stub_5.ibt",
		header:   headers.NewHeader(nil, nil, nil, sessionWithSubsession2),
	}

	stub6 := Stub{
		filepath: "stub_6.ibt",
		header:   headers.NewHeader(nil, nil, nil, sessionWithSubsession3),
	}

	stub7 := Stub{
		filepath: "stub_7.ibt",
		header:   headers.NewHeader(nil, nil, nil, sessionNoResultPos),
	}

	stub8 := Stub{
		filepath: "stub_8.ibt",
		header:   headers.NewHeader(nil, nil, nil, sessionResultPos),
	}

	t.Run("test Group() with regular pattern", func(t *testing.T) {
		stubs := StubGroup{stub1, stub2, stub3, stub4, stub5, stub6, stub7, stub8}

		grouped := stubs.Group()

		if len(grouped) != 4 {
			t.Errorf("expected length of stub group to be %d. received %d. group: %v", 4, len(grouped), grouped)
		}

		if grouped[0][1].filepath != "stub_8.ibt" {
			t.Errorf("expected second item of the first group to have filename %s. received %s", "stub_8.ibt", grouped[0][1].filepath)
		}

		if grouped[1][2].filepath != "stub_3.ibt" {
			t.Errorf("expected third item of the second group to have filename %s. received %s", "stub_3.ibt", grouped[1][2].filepath)
		}

		if grouped[2][1].filepath != "stub_5.ibt" {
			t.Errorf("expected second item of the third group to have filename %s. received %s", "stub_3.ibt", grouped[1][2].filepath)
		}
	})

	t.Run("test groupTestSessionStubs() with irregular pattern", func(t *testing.T) {
		stubs := StubGroup{stub8, stub7, stub1, stub2, stub3}

		grouped := stubs.Group()

		if len(grouped) != 3 {
			t.Errorf("expected length of stub group to be %d. received %d. group: %v", 2, len(grouped), grouped)
		}

		if grouped[0][0].filepath != "stub_8.ibt" {
			t.Errorf("expected first item of the first group to have filename %s. received %s", "stub_8.ibt", grouped[0][0].filepath)
		}

		if grouped[1][0].filepath != "stub_7.ibt" {
			t.Errorf("expected first item of the second group to have filename %s. received %s", "stub_7.ibt", grouped[1][0].filepath)
		}

		if grouped[2][0].filepath != "stub_1.ibt" {
			t.Errorf("expected first item of the third group to have filename %s. received %s", "stub_1.ibt", grouped[2][0].filepath)
		}
	})
}
