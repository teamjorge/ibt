package ibt

import (
	"errors"
	"os"
	"testing"

	"github.com/teamjorge/ibt/headers"
)

type testProcessor struct {
	results []map[string]headers.VarHeader
	session *headers.Session
}

func (t *testProcessor) Process(input map[string]headers.VarHeader, hasNext bool, session *headers.Session) error {
	t.results = append(t.results, input)
	t.session = session

	return nil
}

func (t *testProcessor) Whitelist() []string { return []string{"LapCurrentLapTime"} }

type testErrorProcessor struct{}

func (t *testErrorProcessor) Process(input map[string]headers.VarHeader, hasNext bool, session *headers.Session) error {
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

	t.Run("test process normal processor", func(t *testing.T) {
		proc := testProcessor{}

		if err := Process(f, testHeaders, &proc); err != nil {
			t.Errorf("expected Process() to run without err. received error: %v", err)
		}

		valueToCheck := proc.results[0]["LapCurrentLapTime"].Value.(float32)
		if valueToCheck != 37.678566 {
			t.Errorf("expected value to check to be %f. got %f", 37.678566, valueToCheck)
		}

		valueToCheck = proc.results[69]["LapCurrentLapTime"].Value.(float32)
		if valueToCheck != 38.828568 {
			t.Errorf("expected value to check to be %f. got %f", 38.828568, valueToCheck)
		}
	})

	t.Run("test process err processor", func(t *testing.T) {
		proc := testErrorProcessor{}

		if err := Process(f, testHeaders, &proc); err == nil {
			t.Error("expected Process() to return an error")
		}
	})
}
