package ibt

import (
	"bytes"
	"os"
	"reflect"
	"sort"
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

		if !reflect.DeepEqual(h, header) {
			t.Errorf("expected %v = %v", h, header)
		}
	})

	t.Run("stubs Time() valid time", func(t *testing.T) {
		parsedTime := testStub.Time()

		expectedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-06-24 19:45:36")
		if !parsedTime.Equal(expectedTime) {
			t.Errorf("expected time to be %v but got %v", expectedTime.UTC(), parsedTime.UTC())
		}
	})

	t.Run("stubs DriverIdx()", func(t *testing.T) {
		h := headers.Header{
			TelemetryHeader: &headers.TelemetryHeader{},
			DiskHeader:      &headers.DiskHeader{},
			VarHeader:       map[string]headers.VarHeader{},
			SessionInfo:     &headers.Session{DriverInfo: headers.DriverInfo{DriverCarIdx: 15}},
		}

		driverIdxStub := Stub{filepath: ".testing/valid_test_file.ibt", header: &h}

		if driverIdxStub.DriverIdx() != 15 {
			t.Errorf("expected driver idx to be 15, but got %d", driverIdxStub.DriverIdx())
		}
	})

	t.Run("stubs Open() valid file", func(t *testing.T) {
		stub := Stub{filepath: ".testing/valid_test_file.ibt"}

		f, err := stub.Open()
		if err != nil {
			t.Errorf("did not expect an error when opening file %s. received: %v", ".testing/valid_test_file.ibt", err)
		}

		buf := make([]byte, 2)
		if _, err := f.Read(buf); err != nil {
			t.Errorf("did not expect an error when reading 2 bytes from file %s. received: %v", ".testing/valid_test_file.ibt", err)
		}

		if !bytes.Equal(buf, []byte{0x2, 0x0}) {
			t.Errorf("expected buf to be %v. received %v", []byte{0x2, 0x0}, buf)
		}
	})

	t.Run("stubs Open() invalid file", func(t *testing.T) {
		stub := Stub{filepath: ".testing/disappear_here.ibt"}

		_, err := stub.Open()
		if err == nil {
			t.Errorf("expected an error when opening a non-existent file %s", ".testing/disappear_here.ibt")
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

		if parsedStub.header.TelemetryHeader.BufOffset != 53764 {
			t.Errorf("expected parsedStub to have a BufOffset of %d. received: %d", 53764, parsedStub.header.TelemetryHeader.BufOffset)
		}

		if parsedStub.header.DiskHeader.StartDate != 1719258336 {
			t.Errorf("expected parsedStub to have a StartDate of %d. received: %d", 1719258336, parsedStub.header.DiskHeader.StartDate)
		}

		if parsedStub.header.SessionInfo.DriverInfo.DriverCarVersion != "2024.05.28.02" {
			t.Errorf("expected parsedStub to have a DriverCarVersion of %s. received: %s", "2024.05.28.02", parsedStub.header.SessionInfo.DriverInfo.DriverCarVersion)
		}

		if len(parsedStub.header.VarHeader) != 276 {
			t.Errorf("expected parsedStub to have a %d VarHeaders. received: %d", 276, len(parsedStub.header.VarHeader))
		}
	})

	t.Run("test parseStub invalid file", func(t *testing.T) {
		_, err := parseStub(".testing/invalid_test_file.ibt")
		if err == nil {
			t.Error("expected an error from parseStub() when reading an invalid file")
		}
	})

	t.Run("test parseStub non-existent file", func(t *testing.T) {
		_, err := parseStub(".testing/disappear_here.ibt")
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

		if parsedStubs[0].header.TelemetryHeader.BufOffset != 53764 {
			t.Errorf("expected the parsed stub to have a BufOffset of %d. received: %d", 53764, parsedStubs[0].header.TelemetryHeader.BufOffset)
		}

		if parsedStubs[1].header.DiskHeader.StartDate != 1719258336 {
			t.Errorf("expected the parsed stub to have a StartDate of %d. received: %d", 1719258336, parsedStubs[1].header.DiskHeader.StartDate)
		}

		if parsedStubs[0].header.SessionInfo.DriverInfo.DriverCarVersion != "2024.05.28.02" {
			t.Errorf("expected the parsed stub to have a DriverCarVersion of %s. received: %s", "2024.05.28.02", parsedStubs[0].header.SessionInfo.DriverInfo.DriverCarVersion)
		}

		if len(parsedStubs[1].header.VarHeader) != 276 {
			t.Errorf("expected the parsed stub to have a %d VarHeaders. received: %d", 276, len(parsedStubs[1].header.VarHeader))
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

		if parsedStubs[0].header.TelemetryHeader.BufOffset != 53764 {
			t.Errorf("expected the parsed stub to have a BufOffset of %d. received: %d", 53764, parsedStubs[0].header.TelemetryHeader.BufOffset)
		}
	})
}

func TestGroupTestSessionStubs(t *testing.T) {
	makeHeader := func(subSessionId int, ResultsPositions interface{}, ts int64) *headers.Header {
		return &headers.Header{
			DiskHeader: &headers.DiskHeader{StartDate: ts},
			SessionInfo: &headers.Session{
				WeekendInfo: headers.WeekendInfo{SubSessionID: subSessionId},
				SessionInfo: headers.SessionInfo{
					Sessions: []headers.Sessions{{ResultsPositions: ResultsPositions}},
				},
			},
		}
	}

	now := time.Now()

	stub1 := Stub{
		filepath: "stub_1.ibt",
		header:   makeHeader(3, nil, now.Add(-120*time.Minute).Unix()),
	}

	stub2 := Stub{
		filepath: "stub_2.ibt",
		header:   makeHeader(0, 1, now.Add(-60*time.Minute).Unix()),
	}

	stub3 := Stub{
		filepath: "stub_3.ibt",
		header:   makeHeader(0, 1, now.Unix()),
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
	makeHeader := func(subSessionId int, ResultsPositions interface{}, ts int64) *headers.Header {
		return &headers.Header{
			DiskHeader: &headers.DiskHeader{StartDate: ts},
			SessionInfo: &headers.Session{
				WeekendInfo: headers.WeekendInfo{SubSessionID: subSessionId},
				SessionInfo: headers.SessionInfo{
					Sessions: []headers.Sessions{
						{ResultsPositions: ResultsPositions},
					},
				},
			},
		}
	}

	now := time.Now()

	stub1 := Stub{
		filepath: "stub_1.ibt",
		header:   makeHeader(1, nil, now.Add(-420*time.Minute).Unix()),
	}

	stub2 := Stub{
		filepath: "stub_2.ibt",
		header:   makeHeader(1, nil, now.Add(-360*time.Minute).Unix()),
	}

	stub3 := Stub{
		filepath: "stub_3.ibt",
		header:   makeHeader(1, nil, now.Add(-300*time.Minute).Unix()),
	}

	stub4 := Stub{
		filepath: "stub_4.ibt",
		header:   makeHeader(2, nil, now.Add(-240*time.Minute).Unix()),
	}

	stub5 := Stub{
		filepath: "stub_5.ibt",
		header:   makeHeader(2, nil, now.Add(-180*time.Minute).Unix()),
	}

	stub6 := Stub{
		filepath: "stub_6.ibt",
		header:   makeHeader(3, nil, now.Add(-120*time.Minute).Unix()),
	}

	stub7 := Stub{
		filepath: "stub_7.ibt",
		header:   makeHeader(0, nil, now.Add(-60*time.Minute).Unix()),
	}

	stub8 := Stub{
		filepath: "stub_8.ibt",
		header:   makeHeader(0, 1, now.Unix()),
	}

	t.Run("test Group() with regular pattern", func(t *testing.T) {
		stubs := StubGroup{stub1, stub2, stub3, stub4, stub5, stub6, stub7, stub8}

		grouped := stubs.Group()

		if len(grouped) != 4 {
			t.Errorf("expected length of stub group to be %d. received %d. group: %v", 4, len(grouped), grouped)
		}

		if grouped[0][1].filepath != "stub_2.ibt" {
			t.Errorf("expected second item of the first group to have filename %s. received %s", "stub_2.ibt", grouped[0][1].filepath)
		}

		if grouped[1][1].filepath != "stub_5.ibt" {
			t.Errorf("expected second item of the second group to have filename %s. received %s", "stub_5.ibt", grouped[1][1].filepath)
		}

		if grouped[2][0].filepath != "stub_6.ibt" {
			t.Errorf("expected first item of the third group to have filename %s. received %s", "stub_6.ibt", grouped[2][0].filepath)
		}

		if grouped[3][1].filepath != "stub_8.ibt" {
			t.Errorf("expected first item of the third group to have filename %s. received %s", "stub_8.ibt", grouped[3][1].filepath)
		}
	})

	t.Run("test groupTestSessionStubs() with irregular pattern", func(t *testing.T) {
		stubs := StubGroup{stub7, stub8, stub1, stub4, stub5, stub2, stub3, stub6}

		grouped := stubs.Group()

		if len(grouped) != 4 {
			t.Errorf("expected length of stub group to be %d. received %d. group: %v", 4, len(grouped), grouped)
		}

		if grouped[0][1].filepath != "stub_2.ibt" {
			t.Errorf("expected second item of the first group to have filename %s. received %s", "stub_2.ibt", grouped[0][1].filepath)
		}

		if grouped[1][1].filepath != "stub_5.ibt" {
			t.Errorf("expected second item of the second group to have filename %s. received %s", "stub_5.ibt", grouped[1][1].filepath)
		}

		if grouped[2][0].filepath != "stub_6.ibt" {
			t.Errorf("expected first item of the third group to have filename %s. received %s", "stub_6.ibt", grouped[2][0].filepath)
		}

		if grouped[3][1].filepath != "stub_8.ibt" {
			t.Errorf("expected first item of the third group to have filename %s. received %s", "stub_8.ibt", grouped[3][1].filepath)
		}
	})
}

func TestStubGroupSorting(t *testing.T) {
	makeHeader := func(ts int64) *headers.Header {
		return &headers.Header{DiskHeader: &headers.DiskHeader{StartDate: ts}}
	}

	t.Run("test stub group sort", func(t *testing.T) {
		stubGroup := StubGroup{
			Stub{filepath: "5.ibt", header: makeHeader(time.Now().Unix())},
			Stub{filepath: "3.ibt", header: makeHeader(time.Now().Add(-120 * time.Minute).Unix())},
			Stub{filepath: "1.ibt", header: makeHeader(time.Now().Add(-360 * time.Minute).Unix())},
			Stub{filepath: "4.ibt", header: makeHeader(time.Now().Add(-60 * time.Minute).Unix())},
			Stub{filepath: "2.ibt", header: makeHeader(time.Now().Add(-240 * time.Minute).Unix())},
		}

		sort.Sort(stubGroup)

		expectedOrder := []string{"1.ibt", "2.ibt", "3.ibt", "4.ibt", "5.ibt"}

		for idx, item := range expectedOrder {
			if item != stubGroup[idx].filepath {
				t.Errorf("item at index %d did not match expected value. expected %s. received %s", idx, item, stubGroup[idx].filepath)
			}
		}
	})

	t.Run("test stub group sort", func(t *testing.T) {
		now := time.Now()

		stubGroup := StubGroup{
			Stub{filepath: "5.ibt", header: makeHeader(now.Unix())},
			Stub{filepath: "3.ibt", header: makeHeader(now.Add(-120 * time.Minute).Unix())},
			Stub{filepath: "1.ibt", header: makeHeader(now.Add(-360 * time.Minute).Unix())},
			Stub{filepath: "4.ibt", header: makeHeader(now.Add(-60 * time.Minute).Unix())},
			Stub{filepath: "2.ibt", header: makeHeader(now.Add(-240 * time.Minute).Unix())},
		}

		sort.Sort(stubGroup)

		expectedOrder := []string{"1.ibt", "2.ibt", "3.ibt", "4.ibt", "5.ibt"}

		for idx, item := range expectedOrder {
			if item != stubGroup[idx].filepath {
				t.Errorf("item at index %d did not match expected value. expected %s. received %s", idx, item, stubGroup[idx].filepath)
			}
		}
	})

	t.Run("test stub group grouping sort", func(t *testing.T) {
		now := time.Now()

		stubGroup1 := StubGroup{
			Stub{filepath: "7.ibt", header: makeHeader(now.Add(-60 * time.Minute).Unix())},
			Stub{filepath: "8.ibt", header: makeHeader(now.Unix())},
		}

		stubGroup2 := StubGroup{
			Stub{filepath: "4.ibt", header: makeHeader(now.Add(-240 * time.Minute).Unix())},
			Stub{filepath: "5.ibt", header: makeHeader(now.Add(-180 * time.Minute).Unix())},
			Stub{filepath: "6.ibt", header: makeHeader(now.Add(-120 * time.Minute).Unix())},
		}

		stubGroup3 := StubGroup{
			Stub{filepath: "1.ibt", header: makeHeader(now.Add(-420 * time.Minute).Unix())},
			Stub{filepath: "2.ibt", header: makeHeader(now.Add(-360 * time.Minute).Unix())},
			Stub{filepath: "3.ibt", header: makeHeader(now.Add(-300 * time.Minute).Unix())},
		}

		stubGroupGrouping := StubGroupGrouping{stubGroup3, stubGroup1, stubGroup2}

		sort.Sort(stubGroupGrouping)

		expectedOrder := []int64{
			now.Add(-420 * time.Minute).Unix(),
			now.Add(-240 * time.Minute).Unix(),
			now.Add(-60 * time.Minute).Unix(),
		}

		for idx, item := range expectedOrder {
			if item != stubGroupGrouping[idx][0].Time().Unix() {
				t.Errorf("item at index %d did not match expected value. expected %d. received %d", idx, item, stubGroupGrouping[idx][0].Time().Unix())
			}
		}
	})
}

// Order is not always preserved with the slices
// Check based on length
