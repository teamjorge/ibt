package headers

import (
	"os"
	"reflect"
	"testing"
)

var expectedTelemetryHeader = TelemetryHeader{
	Version:           2,
	Status:            1,
	TickRate:          60,
	SessionInfoUpdate: 0,
	SessionInfoOffset: 39888,
	SessionInfoLength: 13876,
	NumVars:           276,
	VarHeaderOffset:   144,
	NumBuf:            1,
	BufLen:            1072,
	BufOffset:         53764,
}

func TestTelemetryHeaders(t *testing.T) {
	t.Run("valid header file", func(t *testing.T) {
		f, err := os.Open("../.testing/valid_test_file.ibt")
		if err != nil {
			t.Errorf("failed to open testing file - %v", err)
			return
		}
		defer f.Close()

		output, err := ReadTelemetryHeader(f)
		if err != nil {
			t.Errorf("failed to parse telemetry header for testing file - %v", err)
			return
		}

		if !reflect.DeepEqual(*output, expectedTelemetryHeader) {
			t.Errorf("expected telemHeader does not match actual. \nexpected: %+v\n \nactual: %+v\n", expectedTelemetryHeader, *output)
		}
	})

	t.Run("invalid header file", func(t *testing.T) {
		f, err := os.Open("../.testing/invalid_test_file.ibt")
		if err != nil {
			t.Errorf("failed to open testing file - %v", err)
			return
		}
		defer f.Close()

		_, err = ReadTelemetryHeader(f)
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

		_, err = ReadTelemetryHeader(f)
		if err == nil {
			t.Error("expected telemHeader parsing of empty file to return an error")
		}
	})
}
