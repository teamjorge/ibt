package ibt

import (
	"context"
	"errors"
	"os"
	"sort"
	"testing"

	"github.com/teamjorge/ibt/headers"
)

type testProcessor struct {
	results   []Tick
	session   *headers.Session
	whitelist []string
}

func (t *testProcessor) Process(input Tick, hasNext bool, session *headers.Session) error {
	t.results = append(t.results, input)
	t.session = session

	return nil
}

func (t *testProcessor) Whitelist() []string { return t.whitelist }

type testErrorProcessor struct{}

func (t *testErrorProcessor) Process(input Tick, hasNext bool, session *headers.Session) error {
	return errors.New("unit test error")
}

func (t *testErrorProcessor) Whitelist() []string { return []string{"LapCurrentLapTime"} }

func TestProcess(t *testing.T) {
	f, err := os.Open(".testing/valid_test_file.ibt")
	if err != nil {
		t.Errorf("failed to open testing file - %v", err)
		return
	}
	defer f.Close()

	testHeaders, err := headers.ParseHeaders(f)
	if err != nil {
		t.Errorf("failed to parse header for testing file - %v", err)
		return
	}

	stubs := StubGroup{
		{filepath: ".testing/valid_test_file.ibt", header: testHeaders},
	}

	t.Run("test Process() normal processor", func(t *testing.T) {
		proc := testProcessor{whitelist: []string{"LapCurrentLapTime"}}

		if err := Process(context.Background(), stubs, &proc); err != nil {
			t.Errorf("expected Process() to run without err. received error: %v", err)
		}

		valueToCheck := proc.results[0]["LapCurrentLapTime"].(float32)
		if valueToCheck != 37.678566 {
			t.Errorf("expected value to check to be %f. got %f", 37.678566, valueToCheck)
		}

		valueToCheck = proc.results[69]["LapCurrentLapTime"].(float32)
		if valueToCheck != 38.828568 {
			t.Errorf("expected value to check to be %f. got %f", 38.828568, valueToCheck)
		}
	})

	t.Run("test Process() err processor", func(t *testing.T) {
		proc := testErrorProcessor{}

		if err := Process(context.Background(), stubs, &proc); err == nil {
			t.Error("expected Process() to return an error")
		}
	})

	t.Run("test process() invalid file", func(t *testing.T) {
		proc := testProcessor{whitelist: []string{"LapCurrentLapTime"}}

		invalidStub := Stub{
			filepath: "disappear_here",
		}

		if err := process(context.Background(), invalidStub, &proc); err == nil {
			t.Errorf("expected Process() to exit with a file error")
		}
	})

	t.Run("test process() invalid file", func(t *testing.T) {
		proc := testProcessor{whitelist: []string{"LapCurrentLapTime"}}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		if err := process(ctx, stubs[0], &proc); err == nil {
			t.Errorf("expected process() to exit with a context done error")
		}
	})
}

func TestWhitelistParsing(t *testing.T) {
	f, err := os.Open(".testing/valid_test_file.ibt")
	if err != nil {
		t.Errorf("failed to open testing file - %v", err)
		return
	}
	defer f.Close()

	testHeaders, err := headers.ParseHeaders(f)
	if err != nil {
		t.Errorf("failed to parse header for testing file - %v", err)
		return
	}

	varHeader := testHeaders.VarHeader()

	t.Run("test parseAndValidateWhitelist empty", func(t *testing.T) {
		proc := testProcessor{whitelist: []string{}}

		cols := parseAndValidateWhitelist(varHeader, &proc)

		if len(cols) != 276 {
			t.Errorf("expected %d columns to be in whitelist when returning an empty Whitelist() value. found %d", 276, len(cols))
		}
	})

	t.Run("test parseAndValidateWhitelist *", func(t *testing.T) {
		proc := testProcessor{whitelist: []string{"something", "*"}}

		cols := parseAndValidateWhitelist(varHeader, &proc)

		if len(cols) != 276 {
			t.Errorf("expected %d columns to be in whitelist when returning an empty Whitelist() value. found %d", 276, len(cols))
		}
	})

	t.Run("test parseAndValidateWhitelist valid and invalid columns", func(t *testing.T) {
		proc := testProcessor{whitelist: []string{"something", "Speed", "is", "Gear", "wrong"}}

		cols := parseAndValidateWhitelist(varHeader, &proc)

		if len(cols) != 2 {
			t.Errorf("expected %d columns to be in whitelist when returning an empty Whitelist() value. found %d", 2, len(cols))
		}

		sort.Strings(cols)

		if cols[0] != "Gear" || cols[1] != "Speed" {
			t.Errorf("expected columns to be %v. received %v", []string{"Gear", "Speed"}, cols)
		}
	})

	t.Run("test buildWhitelist with 1 * and 2 normal", func(t *testing.T) {
		proc1 := testProcessor{whitelist: []string{"something", "Speed", "is", "Gear", "wrong"}}
		proc2 := testProcessor{whitelist: []string{"*"}}

		cols := buildWhitelist(varHeader, []Processor{&proc1, &proc2}...)

		if len(cols) != 276 {
			t.Errorf("expected %d columns to be in whitelist when returning an empty Whitelist() value. found %d", 276, len(cols))
		}
	})

	t.Run("test buildWhitelist with 1 empty and 2 normal", func(t *testing.T) {
		proc1 := testProcessor{whitelist: []string{"something", "Speed", "is", "Gear", "wrong"}}
		proc2 := testProcessor{whitelist: nil}

		cols := buildWhitelist(varHeader, []Processor{&proc1, &proc2}...)

		if len(cols) != 276 {
			t.Errorf("expected %d columns to be in whitelist when returning an empty Whitelist() value. found %d", 276, len(cols))
		}
	})

	t.Run("test buildWhitelist with duplicate and invalid columns", func(t *testing.T) {
		proc1 := testProcessor{whitelist: []string{"something", "Speed", "is", "Gear", "wrong"}}
		proc2 := testProcessor{whitelist: []string{"BrakeRaw", "Speed", "ThrottleRaw", "Gear", "wrong"}}

		cols := buildWhitelist(varHeader, []Processor{&proc1, &proc2}...)

		if len(cols) != 4 {
			t.Errorf("expected %d columns to be in whitelist when returning an empty Whitelist() value. found %d", 4, len(cols))
		}

		sort.Strings(cols)

		if cols[0] != "BrakeRaw" || cols[1] != "Gear" || cols[2] != "Speed" || cols[3] != "ThrottleRaw" {
			t.Errorf("expected columns to be %v. received %v", []string{"BrakeRaw", "Gear", "Speed", "ThrottleRaw"}, cols)
		}
	})
}
