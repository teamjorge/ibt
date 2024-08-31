package headers

import (
	"os"
	"reflect"
	"testing"
)

var (
	expectedLiveVarBuffers []VarBuffer = []VarBuffer{
		{
			TickCount: 2004,
			BufOffset: 1114224,
		},
		{
			TickCount: 2005,
			BufOffset: 1138800,
		},
		{
			TickCount: 2003,
			BufOffset: 1163376,
		},
	}
	expectedVarBuffers []VarBuffer = []VarBuffer{
		{
			TickCount: 54661,
			BufOffset: 53764,
		},
	}
)

func TestReadVarBufferHeaders(t *testing.T) {
	t.Run("test live telemetry file", func(t *testing.T) {
		f, err := os.Open("../.testing/live_test_file.ibt")
		if err != nil {
			t.Errorf("failed to open testing file - %v", err)
			return
		}
		defer f.Close()

		output, err := ReadVarBufferHeaders(f, 3)
		if err != nil {
			t.Errorf("failed to parse telemetry header for testing file - %v", err)
			return
		}

		if !reflect.DeepEqual(output[0], expectedLiveVarBuffers[0]) {
			t.Errorf("expected varBuffer header does not match actual. \nexpected: %+v\n \nactual: %+v\n", expectedLiveVarBuffers[0], output[0])
		}

		if !reflect.DeepEqual(output[1], expectedLiveVarBuffers[1]) {
			t.Errorf("expected varBuffer header does not match actual. \nexpected: %+v\n \nactual: %+v\n", expectedLiveVarBuffers[1], output[1])
		}

		if !reflect.DeepEqual(output[2], expectedLiveVarBuffers[2]) {
			t.Errorf("expected varBuffer header does not match actual. \nexpected: %+v\n \nactual: %+v\n", expectedLiveVarBuffers[2], output[2])
		}
	})

	t.Run("test valid file", func(t *testing.T) {
		f, err := os.Open("../.testing/valid_test_file.ibt")
		if err != nil {
			t.Errorf("failed to open testing file - %v", err)
			return
		}
		defer f.Close()

		output, err := ReadVarBufferHeaders(f, expectedTelemetryHeader.NumBuf)
		if err != nil {
			t.Errorf("failed to parse telemetry header for testing file - %v", err)
			return
		}

		if !reflect.DeepEqual(output[0], expectedVarBuffers[0]) {
			t.Errorf("expected varBuffer header does not match actual. \nexpected: %+v\n \nactual: %+v\n", expectedVarBuffers[0], output[0])
		}
	})

	t.Run("invalid header file", func(t *testing.T) {
		f, err := os.Open("../.testing/invalid_test_file.ibt")
		if err != nil {
			t.Errorf("failed to open testing file - %v", err)
			return
		}
		defer f.Close()

		_, err = ReadVarBufferHeaders(f, expectedTelemetryHeader.NumBuf)
		if err == nil {
			t.Error("expected telemHeader parsing of invalid file to return an error")
		}
	})

	t.Run("empty file", func(t *testing.T) {
		f, err := os.Open("../.testing/empty_test_file.ibt")
		if err != nil {
			t.Errorf("failed to open testing file - %v", err)
			return
		}
		defer f.Close()

		_, err = ReadVarBufferHeaders(f, expectedTelemetryHeader.NumBuf)
		if err == nil {
			t.Error("expected telemHeader parsing of empty file to return an error")
		}
	})
}
